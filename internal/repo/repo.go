package repo

import (
	"context"
	"log/slog"

	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/internal/repo/human"
	"github.com/P1xart/effective_mobile_service/pkg/postgresql"
)

type Human interface {
	Create(ctx context.Context, body *entity.Human) (*entity.Human, error)
	GetAll(ctx context.Context, filters *entity.HumanFilters) ([]entity.Human, error)
	UpdateByID(ctx context.Context, id string, updates *entity.Human) (*entity.Human, error)
	DeleteByID(ctx context.Context, id string) error
}

type Repositories struct {
	Human
}

func NewRepositories(log *slog.Logger, pg *postgresql.Postgres) *Repositories {
	return &Repositories{
		Human: human.NewRepo(log, pg),
	}
}
