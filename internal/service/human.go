package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/P1xart/effective_mobile_service/internal/config"
	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/internal/repo"
	"github.com/P1xart/effective_mobile_service/internal/repo/repoerrors"
	"github.com/P1xart/effective_mobile_service/pkg/logger"
)

type HumanService struct {
	log        *slog.Logger
	httpClient *http.Client

	apiUrls   *config.ApiUrls
	humanRepo repo.Human
}

func NewHumanService(log *slog.Logger, humanRepo repo.Human, apiUrls *config.ApiUrls) *HumanService {
	log = log.With(slog.String("component", "human service"))

	client := http.DefaultClient

	return &HumanService{
		log:        log,
		httpClient: client,
		apiUrls:    apiUrls,

		humanRepo: humanRepo,
	}
}

func (s *HumanService) Create(ctx context.Context, body *HumanInput) (*entity.Human, error) {
	err := s.fillUserData(ctx, body)
	if err != nil {
		s.log.Error("failed to fill user data from api", logger.Error(err))
		return nil, err
	}

	human, err := s.humanRepo.Create(ctx, &entity.Human{
		Name:        body.Name,
		Surname:     body.Surname,
		Potronymic:  body.Potronymic,
		Age:         body.Age,
		Gender:      body.Gender,
		Nationality: body.Nationality,
	})
	if err != nil {
		s.log.Error("failed to create human", logger.Error(err))
        return nil, err
	}

	return human, nil
}

func (s *HumanService) GetAll(ctx context.Context, filters *entity.HumanFilters) ([]entity.Human, error) {
	return s.humanRepo.GetAll(ctx, filters)
}

func (s *HumanService) DeleteByID(ctx context.Context, id string) error {
	if err := s.humanRepo.DeleteByID(ctx, id); err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return ErrHumanNotFound
		}

		return err
	}

	return nil
}

func (s *HumanService) UpdateByID(ctx context.Context, id string, updates *HumanInput) (*entity.Human, error) {
	updatedHuman, err := s.humanRepo.UpdateByID(ctx, id, &entity.Human{
		Name:        updates.Name,
		Surname:     updates.Surname,
		Potronymic:  updates.Potronymic,
		Age:         updates.Age,
		Gender:      updates.Gender,
		Nationality: updates.Nationality,
	})
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, ErrHumanNotFound
		}

		return nil, err
	}

	return updatedHuman, nil
}

func (s *HumanService) fillUserData(ctx context.Context, body *HumanInput) error {
    groupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	g, groupCtx := errgroup.WithContext(groupCtx)

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

				body.Age = resp.Age
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
					body.Nationality = ""
					return nil
				}
				body.Nationality = resp.Country[0].CountryId
				return nil
			},
		},
	}

	for _, task := range tasks {
		task := task
		g.Go(func() error {
			req, err := http.NewRequestWithContext(groupCtx, "GET", task.url, nil)
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

    return nil
}