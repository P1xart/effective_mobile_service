package human

import (
	"context"
	"log/slog"

	"github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/internal/repo/repoerrors"
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

func (r *Repo) Create(ctx context.Context, body *entity.Human) (*entity.Human, error) {
	q, args, err := r.Builder.Insert("humans").Columns("name", "surname", "potronymic", "age", "gender", "nationality").
		Values(body.Name, body.Surname, body.Potronymic, body.Age, body.Gender, body.Nationality).Suffix("RETURNING id, age, gender, nationality").ToSql()
	if err != nil {
		r.log.Error("failed to make query", logger.Error(err))
		return nil, err
	}

	r.log.Debug("create human query", slog.String("query", q))

	if err := r.Pool.QueryRow(ctx, q, args...).Scan(
		&body.ID,
		&body.Age,
		&body.Gender,
		&body.Nationality,
	); err != nil {
		r.log.Error("failed to create human", logger.Error(err))
	}

	return body, nil
}

func (r *Repo) GetAll(ctx context.Context, filters *entity.HumanFilters) ([]entity.Human, error) {
	qb := r.Builder.
		Select("*").
		From("humans")

	if filters.AgeFrom != 0 {
		qb = qb.Where(squirrel.GtOrEq{"age": filters.AgeFrom})
	}

	if filters.AgeTo != 0 {
		qb = qb.Where(squirrel.LtOrEq{"age": filters.AgeTo})
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

func (r *Repo) UpdateByID(ctx context.Context, id string, updates *entity.Human) (*entity.Human, error) {
	builder := r.Builder.Update("humans")

	if updates.Name != "" {
		builder = builder.Set("name", updates.Name)
	}

	if updates.Surname != "" {
		builder = builder.Set("surname", updates.Surname)
	}

	if updates.Potronymic != "" {
		builder = builder.Set("potronymic", updates.Potronymic)
	}

	if updates.Age != 0 {
		builder = builder.Set("age", updates.Age)
	}

	if updates.Gender != "" {
		builder = builder.Set("gender", updates.Gender)
	}

	if updates.Nationality != "" {
		builder = builder.Set("nationality", updates.Nationality)
	}

	q, args, err := builder.Where(squirrel.Eq{"id": id}).Suffix("RETURNING id, age, gender, nationality").ToSql()
	if err != nil {
		r.log.Error("failed to build SQL query", slog.Any("id", id), logger.Error(err))
		return nil, err
	}

	r.log.Debug("update human query", slog.String("query", q))

	if err := r.Pool.QueryRow(ctx, q, args...).Scan(
		&updates.ID,
		&updates.Age,
		&updates.Gender,
		&updates.Nationality,
	); err != nil {
		r.log.Error("failed to create human", logger.Error(err))
	}

	r.log.Debug("update rows successfuly")

	return updates, nil
}

func (r *Repo) DeleteByID(ctx context.Context, id string) error {
	q, args, err := r.Builder.Delete("humans").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		r.log.Error("failed to make query", logger.Error(err))
		return err
	}

	r.log.Debug("delete human by id query", slog.String("query", q))

	result, err := r.Pool.Exec(ctx, q, args...)
	if err != nil {
		r.log.Error("failed to delete human by id", slog.Any("id", id), logger.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		r.log.Warn("no human found with the given ID to delete", slog.Any("id", id))
		return repoerrors.ErrNotFound
	}

	return nil
}
