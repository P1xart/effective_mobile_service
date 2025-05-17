package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/P1xart/effective_mobile_service/internal/config"
	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/internal/repo"
)

type CreateHuman struct {
	Name       string
	Surname    string
	Potronymic string
	Age        uint8
	Gender     string
	Nationaly  string
}

type Human interface {
	Create(ctx context.Context, body *CreateHuman, apiUrls config.ApiUrls) error
	GetAll(ctx context.Context, filters *entity.HumanFilters) ([]entity.Human, error)
}

type Dependencies struct {
	Log   *slog.Logger
	Repos *repo.Repositories

	SignKey  string
	TokenTTL time.Duration
}

type Services struct {
	Human Human
}

func NewServices(deps *Dependencies, apiUrls *config.ApiUrls) *Services {
	services := &Services{
		Human: NewHumanService(deps.Log, deps.Repos.Human, apiUrls),
	}

	return services
}
