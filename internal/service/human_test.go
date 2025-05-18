//go:build integration

package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/P1xart/effective_mobile_service/internal/repo"
	"github.com/P1xart/effective_mobile_service/pkg/utest"

	"github.com/stretchr/testify/require"
)

func TestHumanService_Create(t *testing.T) {
	require.NoError(t, prepareErr)

	log.Debug("test configuration", slog.Any("cfg", cfg))

	defer utest.TeardownTable(log, pg, "humans")

	repositories := repo.NewRepositories(log, pg)

	humanService := NewHumanService(log, repositories.Human, &cfg.API)

	ctx := context.Background()

	expectedHuman := &HumanInput{
		Name: "Mihail",
		Surname: "Dmitrievich",
	}

	err := humanService.Create(ctx, expectedHuman)
	require.NoError(t, err)
}

func TestHumanService_GetAll(t *testing.T) {
	require.NoError(t, prepareErr)

	log.Debug("test configuration", slog.Any("cfg", cfg))

	defer utest.TeardownTable(log, pg, "humans")

	repositories := repo.NewRepositories(log, pg)

	humanService := NewHumanService(log, repositories.Human, &cfg.API)

	ctx := context.Background()

	expectedHuman := &HumanInput{
		Name: "Mihail",
		Surname: "Dmitrievich",
	}

	err := humanService.Create(ctx, expectedHuman)
	require.NoError(t, err)
}

func TestHumanService_UpdateByID(t *testing.T) {
	require.NoError(t, prepareErr)

	log.Debug("test configuration", slog.Any("cfg", cfg))

	defer utest.TeardownTable(log, pg, "humans")

	repositories := repo.NewRepositories(log, pg)

	humanService := NewHumanService(log, repositories.Human, &cfg.API)

	ctx := context.Background()

	expectedHuman := &HumanInput{
		Name: "Mihail",
		Surname: "Dmitrievich",
	}

	err := humanService.Create(ctx, expectedHuman)
	require.NoError(t, err)

	expectedUpdateHuman := &HumanInput{
		Name: "Artem",
		Surname: "Alexandrovich",
		Age: 10,
		Nationality: "RU",
		Gender: "female",
	}

	err = humanService.UpdateByID(ctx, "1", expectedUpdateHuman)
	require.NoError(t, err)
}

func TestHumanService_Delete(t *testing.T) {
	require.NoError(t, prepareErr)

	log.Debug("test configuration", slog.Any("cfg", cfg))

	defer utest.TeardownTable(log, pg, "humans")

	repositories := repo.NewRepositories(log, pg)

	humanService := NewHumanService(log, repositories.Human, &cfg.API)

	ctx := context.Background()

	expectedHuman := &HumanInput{
		Name: "Mihail",
		Surname: "Dmitrievich",
	}

	err := humanService.Create(ctx, expectedHuman)
	require.NoError(t, err)
	err = humanService.DeleteByID(ctx, "1")
	require.NoError(t, err)
}