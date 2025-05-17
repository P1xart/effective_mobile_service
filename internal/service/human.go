package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/P1xart/effective_mobile_service/internal/config"
	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/internal/repo"
	"github.com/P1xart/effective_mobile_service/pkg/logger"
)

type HumanService struct {
	log        *slog.Logger
	httpClient http.Client

	apiUrls   *config.ApiUrls
	humanRepo repo.Human
}

func NewHumanService(log *slog.Logger, humanRepo repo.Human, apiUrls *config.ApiUrls) *HumanService {
	log = log.With(slog.String("component", "human service"))

	client := http.Client{}

	return &HumanService{
		log:        log,
		httpClient: client,
        apiUrls: apiUrls,

		humanRepo: humanRepo,
	}
}

func (s *HumanService) CreateHuman(ctx context.Context, body *CreateHuman, apiUrls config.ApiUrls) error {
	ErrGroupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	g, ErrGroupCtx := errgroup.WithContext(ErrGroupCtx)

	tasks := []struct {
		name    string
		url     string
		handler func([]byte) error
	}{
		{
			name: "age",
			url:  fmt.Sprintf("%s/?name=%s", s.apiUrls.Age, body.Name),
			handler: func(data []byte) error {
				var resp entity.AgeResp
				if err := json.Unmarshal(data, &resp); err != nil {
					s.log.Error("failed to unmarshal age", logger.Error(err))
					return err
				}
				body.Age = uint8(resp.Age)
				return nil
			},
		},
		{
			name: "gender",
			url:  fmt.Sprintf("%s/?name=%s", s.apiUrls.Gender, body.Name),
			handler: func(data []byte) error {
				var resp entity.GenderResp
				if err := json.Unmarshal(data, &resp); err != nil {
					s.log.Error("failed to unmarshal gender", logger.Error(err))
					return err
				}
				body.Gender = resp.Gender
				return nil
			},
		},
		{
			name: "nationality",
			url:  fmt.Sprintf("%s/?name=%s", s.apiUrls.Nation, body.Name),
			handler: func(data []byte) error {
				var resp entity.NationalyResp
				if err := json.Unmarshal(data, &resp); err != nil {
					s.log.Error("failed to unmarshal nationality", logger.Error(err))
					return err
				}
				if len(resp.Country) == 0 {
					body.Nationaly = ""
					return nil
				}
				body.Nationaly = resp.Country[0].CountryId
				return nil
			},
		},
	}

	for _, task := range tasks {
		task := task
		g.Go(func() error {
			req, err := http.NewRequestWithContext(ErrGroupCtx, "GET", task.url, nil)
			if err != nil {
				s.log.Error("failed to make request", slog.String("query", task.name), logger.Error(err))
				return err
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				s.log.Error("failed to do request", slog.String("query", task.name), logger.Error(err))
				return err
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				s.log.Error("failed to read response body", slog.String("query", task.name), logger.Error(err))
				return err
			}

			if err := task.handler(data); err != nil {
				s.log.Error("failed to use handler", slog.String("query", task.name), logger.Error(err))
				return err
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		s.log.Error("failed to get user data from api", logger.Error(err))
		return err
	}

	err := s.humanRepo.CreateHuman(ctx, &entity.Human{
		Name:       body.Name,
		Surname:    body.Surname,
		Potronymic: body.Potronymic,
		Age:        body.Age,
		Gender:     body.Gender,
		Nationaly:  body.Nationaly,
	})
	if err != nil {
		s.log.Error("failed to create human", logger.Error(err))
	}

	return err
}
