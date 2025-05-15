package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"encoding/json"
	"sync"

	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/internal/repo"
	"github.com/P1xart/effective_mobile_service/pkg/logger"
)

type HumanService struct {
	log        *slog.Logger
	httpClient http.Client

	humanRepo repo.Human
}

func NewHumanService(log *slog.Logger, humanRepo repo.Human) *HumanService {
	log = log.With(slog.String("component", "human service"))

	client := http.Client{}

	return &HumanService{
		log:        log,
		httpClient: client,

		humanRepo: humanRepo,
	}
}

func (s *HumanService) CreateHuman(ctx context.Context, body *CreateHuman) error {
	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	defer close(errChan)

	for i := range(3) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var apiUrl string
			switch i {
			case 0: apiUrl = fmt.Sprintf("https://api.agify.io/?name=%s", body.Name)
			case 1: apiUrl = fmt.Sprintf("https://api.genderize.io/?name=%s", body.Name)
			case 2: apiUrl = fmt.Sprintf("https://api.nationalize.io/?name=%s", body.Name)
			}
	
			resp, err := http.Get(apiUrl)
			if err != nil {
				s.log.Error("failed get age from api", slog.String("api", apiUrl), logger.Error(err))
				errChan <- err
				return
			}
			defer resp.Body.Close()
			
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				s.log.Error("failed read response of body", slog.String("api", apiUrl), logger.Error(err))
				errChan <- err
				return
			}
	
			switch i {
			case 0: 
				var AgeResp entity.AgeResp
				err = json.Unmarshal(respBody, &AgeResp)
				body.Age = uint8(AgeResp.Age)
			case 1:
				var GenderResp entity.GenderResp
				err = json.Unmarshal(respBody, &GenderResp)
				body.Gender = GenderResp.Gender
			case 2:
				var NationalyResp entity.NationalyResp
				err = json.Unmarshal(respBody, &NationalyResp)
				body.Nationaly = NationalyResp.Country[0].CountryId
			}
			if err != nil {
				s.log.Error("failed to unmarshalling response", logger.Error(err))
				errChan <- err
				return
			}

			errChan <- nil
		}()
	}
	
	wg.Wait()
	if err := <- errChan; err != nil {
		s.log.Error("failed to make query on api", logger.Error(err))
		return err
	}

	err := s.humanRepo.CreateHuman(ctx, &entity.Human{
		Name: body.Name,
		Surname: body.Surname,
		Potronymic: body.Potronymic,
		Age: body.Age,
		Gender: body.Gender,
		Nationaly: body.Nationaly,
	}); if err != nil {
		s.log.Error("failed to create human", logger.Error(err))
	}

	return err
}
