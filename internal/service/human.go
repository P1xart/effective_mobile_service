package service

import (
	"context"
	"fmt"
	"io"
	"time"
	"log/slog"
	"net/http"
	"encoding/json"
	"golang.org/x/sync/errgroup"

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
            url:  fmt.Sprintf("https://api.agify.io/?name=%s", body.Name),
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
            url:  fmt.Sprintf("https://api.genderize.io/?name=%s", body.Name),
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
            url:  fmt.Sprintf("https://api.nationalize.io/?name=%s", body.Name),
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
