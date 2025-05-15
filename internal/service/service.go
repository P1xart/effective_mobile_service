package service

import (
	"context"
	"log/slog"
	"time"

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
	CreateHuman(ctx context.Context, body *CreateHuman) error
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

func NewServices(deps *Dependencies) *Services {
	services := &Services{
		Human: NewHumanService(deps.Log, deps.Repos.Human),
	}

	return services
}
