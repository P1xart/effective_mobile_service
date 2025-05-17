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
	q, args, err := r.Builder.Insert("humans").Columns("name", "surname", "potronymic", "age", "gender", "nationality").
	Values(body.Name, body.Surname, body.Potronymic, body.Age, body.Gender, body.Nationaly).ToSql()
	if err != nil {
		r.log.Error("failed to make query", logger.Error(err))
		return err
	}
	
	r.log.Info("create human query", slog.String("query", q))

	_, err = r.Pool.Exec(ctx, q, args...)
	if err != nil {
		r.log.Error("failed to insert new human", logger.Error(err))
		return err
	}

	return nil
}
