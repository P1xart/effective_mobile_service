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

	repositories := repo.NewRepositories(log, pg)

	humanService := NewHumanService(log, repositories.Human, &cfg.API)

	ctx := context.Background()

	expectedHuman1 := &HumanInput{
		Name: "Mihail",
		Surname: "Dmitrievich",
		Age: 19,
		Gender: "male",
		Nationality: "RU",
	}
	expectedHuman2 := &HumanInput{
		Name: "Anna",
		Surname: "Ivanova",
		Age: 23,
		Gender: "female",
		Nationality: "US",
	}
	expectedHuman3 := &HumanInput{
		Name: "Jan",
		Surname: "Casanova",
		Age: 56,
		Gender: "male",
		Nationality: "RU",
	}

	_, err := humanService.Create(ctx, expectedHuman1)
	require.NoError(t, err)
	_, err = humanService.Create(ctx, expectedHuman2)
	require.NoError(t, err)
	_, err = humanService.Create(ctx, expectedHuman3)
	require.NoError(t, err)

	femaleHumans, err := humanService.GetAll(ctx, &entity.HumanFilters{Limit: 10, Offset: 0, Gender: []string{"female"}})
	require.NoError(t, err)
	require.Equal(t, expectedHuman2.Gender, femaleHumans[0].Gender)

	ageHumans, err := humanService.GetAll(ctx, &entity.HumanFilters{Limit: 10, Offset: 0, AgeFrom: 10, AgeTo: 30})
	require.NoError(t, err)
	require.Len(t, ageHumans, 2)

	nationHumans, err := humanService.GetAll(ctx, &entity.HumanFilters{Limit: 10, Offset: 0, Nationaly: []string{"US"}})
	require.NoError(t, err)
	require.Len(t, nationHumans, 1)
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

	_, err = humanService.UpdateByID(ctx, "-1", expectedUpdateHuman)
	require.ErrorIs(t, err, ErrHumanNotFound)
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

	err = humanService.DeleteByID(ctx, "-1")
	require.ErrorIs(t, err, ErrHumanNotFound)
}