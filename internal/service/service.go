package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/P1xart/effective_mobile_service/internal/config"
	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/internal/repo"
)

type Human interface {
	Create(ctx context.Context, body *HumanInput, apiUrls config.ApiUrls) error
	GetAll(ctx context.Context, filters *entity.HumanFilters) ([]entity.Human, error)
	UpdateByID(ctx context.Context, id string, updates *HumanInput) error
	DeleteByID(ctx context.Context, id string) error
}

type HumanInput struct {
	Name        string
	Surname     string
	Potronymic  string
	Age         string
	Gender      string
	Nationality string
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
