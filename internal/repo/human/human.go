package human

import (
	"context"
	"log/slog"

	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/pkg/logger"
	"github.com/P1xart/effective_mobile_service/pkg/postgresql"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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

func (r *Repo) Create(ctx context.Context, body *entity.Human) error {
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

func (r *Repo) GetAll(ctx context.Context, filters *entity.HumanFilters) ([]entity.Human, error) {
	qb := r.Builder.
		Select("*").
		From("humans")

	if filters.AgeFrom != 0 {
		qb = qb.Where("age >= ?", filters.AgeFrom)
	}

	if filters.AgeTo != 0 {
		qb = qb.Where("age <= ?", filters.AgeTo)
	}

	if len(filters.Gender) != 0 {
		qb = qb.Where(squirrel.Eq{"gender": filters.Gender})
	}

	if len(filters.Nationaly) != 0 {
		qb = qb.Where(squirrel.Eq{"nationality": filters.Nationaly})
	}

	q, args, err := qb.Limit(filters.Limit).
		Offset(filters.Offset).
		ToSql()

	if err != nil {
		r.log.Error("failed to build query", logger.Error(err))
		return nil, err
	}

	r.log.Debug("get all humans query", slog.String("query", q))

	rows, err := r.Pool.Query(ctx, q, args...)
	if err != nil {
		r.log.Error("failed to get humans from database", logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	humans, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Human])
	if err != nil {
		r.log.Error("failed to collect rows", logger.Error(err))
		return nil, err
	}

	return humans, nil
}