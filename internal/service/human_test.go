//go:build integration

package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/P1xart/effective_mobile_service/internal/entity"
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

	expectedHuman1 := &HumanInput{
		Name: "Mihail",
		Surname: "Dmitrievich",
	}

	human, err := humanService.Create(ctx, expectedHuman1)
	require.NoError(t, err)
	require.Equal(t, human.Name, expectedHuman1.Name)
	require.Equal(t, human.Surname, expectedHuman1.Surname)
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

	_, err := humanService.Create(ctx, expectedHuman)
	require.NoError(t, err)
	humans, err := humanService.GetAll(ctx, &entity.HumanFilters{Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.Equal(t, expectedHuman.Age, humans[0].Age)
	require.Equal(t, expectedHuman.Name, humans[0].Name)
	require.Equal(t, expectedHuman.Nationality, humans[0].Nationality)
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

	createdHuman, err := humanService.Create(ctx, expectedHuman)
	require.NoError(t, err)

	expectedUpdateHuman := &HumanInput{
		Name: "Artem",
		Surname: "Alexandrovich",
		Age: 10,
		Nationality: "RU",
		Gender: "female",
	}

	updatedHuman, err := humanService.UpdateByID(ctx, createdHuman.ID, expectedUpdateHuman)
	require.NoError(t, err)
	require.Equal(t, expectedUpdateHuman.Name, updatedHuman.Name)
	require.Equal(t, expectedUpdateHuman.Age, updatedHuman.Age)
	require.Equal(t, expectedUpdateHuman.Gender, updatedHuman.Gender)
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

	human, err := humanService.Create(ctx, expectedHuman)
	require.NoError(t, err)
	err = humanService.DeleteByID(ctx, human.ID)
	require.NoError(t, err)
	getHuman, err := humanService.GetAll(ctx, &entity.HumanFilters{Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.Len(t, getHuman, 0)
}