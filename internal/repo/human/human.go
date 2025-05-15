package human

import (
	"context"
	"log/slog"

	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/pkg/logger"
	"github.com/P1xart/effective_mobile_service/pkg/postgresql"
)

type Repo struct {
	log *slog.Logger
	*postgresql.Postgres
}

func NewRepo(log *slog.Logger, pg *postgresql.Postgres) *Repo {
	return &Repo{
		log:      log,
		Postgres: pg,
	}
}

func (r *Repo) CreateHuman(ctx context.Context, body *entity.Human) error {
	q := "INSERT INTO humans(name, surname, potronymic, age, gender, nationality) VALUES($1, $2, $3, $4, $5, $6)"
	r.log.Info("init create human query", slog.String("query", q))

	_, err := r.Pool.Exec(ctx, q, body.Name, body.Surname, body.Potronymic, body.Age, body.Gender, body.Nationaly)
	if err != nil {
		r.log.Error("failed to insert new human", logger.Error(err))
		return err
	}

	return nil
}
