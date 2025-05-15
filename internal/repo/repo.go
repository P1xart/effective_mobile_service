package repo

import (
	"context"
	"log/slog"

	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/internal/repo/human"
	"github.com/P1xart/effective_mobile_service/pkg/postgresql"
)

type Human interface {
	CreateHuman(ctx context.Context, body *entity.Human) error
}

type Repositories struct {
	Human
}

func NewRepositories(log *slog.Logger, pg *postgresql.Postgres) *Repositories {
	return &Repositories{
		Human: human.NewRepo(log, pg),
	}
}
