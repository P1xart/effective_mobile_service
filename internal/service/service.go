package service

import (
	"context"
	"log/slog"

	"github.com/P1xart/effective_mobile_service/internal/config"
	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/internal/repo"
)

type Human interface {
	Create(ctx context.Context, body *HumanInput) error
	GetAll(ctx context.Context, filters *entity.HumanFilters) ([]entity.Human, error)
	UpdateByID(ctx context.Context, id string, updates *HumanInput) error
	DeleteByID(ctx context.Context, id string) error
}

type HumanInput struct {
	Name        string
	Surname     string
	Potronymic  string
	Age         int
	Gender      string
	Nationality string
}

type Dependencies struct {
	Log   *slog.Logger
	Repos *repo.Repositories

	Cfg *config.Config
}

type Services struct {
	Human Human
}

func NewServices(deps *Dependencies) *Services {
	services := &Services{
		Human: NewHumanService(deps.Log, deps.Repos.Human, &deps.Cfg.API),
	}

	return services
}
