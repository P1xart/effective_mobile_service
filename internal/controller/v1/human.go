package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/P1xart/effective_mobile_service/internal/controller/v1/request"
	"github.com/P1xart/effective_mobile_service/internal/controller/v1/response"
	_ "github.com/P1xart/effective_mobile_service/internal/entity"
	"github.com/P1xart/effective_mobile_service/internal/service"
	"github.com/P1xart/effective_mobile_service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type humanRoutes struct {
	log   *slog.Logger
	valid *validator.Validate

	humanService service.Human
}

func newHumanRoutes(log *slog.Logger, g *gin.RouterGroup, humanService service.Human) {
	log = log.With(slog.String("component", "human routes"))

	v := validator.New()

	r := &humanRoutes{
		log:          log,
		valid:        v,
		humanService: humanService,
	}

	g.POST("/", r.createNewHuman)
	g.GET("/", r.getHumans)
	g.DELETE("/:id", r.deleteHumanByID)
	g.PATCH("/:id", r.updateCarByID)
}

// @Summary Создание нового человека
// @Description Создание нового человека
// @Tags люди
// @Accept json
// @Param input body request.CreateHuman true "Тело запроса"
// @Success 201 {object} entity.Human
// @Router /v1/human/ [post]
func (r *humanRoutes) createNewHuman(c *gin.Context) {
	var human request.CreateHuman

	err := c.ShouldBindJSON(&human)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = r.valid.Struct(&human)
	if err != nil {
		r.log.Info("error validating data", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	createdHuman, err := r.humanService.Create(c, &service.HumanInput{
		Name:       human.Name,
		Surname:    human.Surname,
		Potronymic: human.Potronymic,
	})
	if err != nil {
		r.log.Error("failed to create human", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	r.log.Info("created new human")
	c.JSON(http.StatusCreated, gin.H{
		"human": createdHuman,
	})
}

// @Summary Получить всех людей
// @Description Получить всех людей
// @Tags люди
// @Accept json
// @Produce json
// @Param age_from query string false "Возраст от"
// @Param age_to query string false "Возраст до"
// @Param gender query string false "Пол"
// @Param nationaly query string false "Национальность"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.GetAllHumans
// @Router /v1/human [get]
func (r *humanRoutes) getHumans(c *gin.Context) {
	filters := buildFilters(c)
	if filters == nil {
		return
	}

	humans, err := r.humanService.GetAll(c, filters)
	if err != nil {
		r.log.Error("failed to get all humans", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response.GetAllHumans{
		Humans: humans,
	})
}

// @Summary Удалить человека по его идентификатору
// @Description Удалить человека по его идентификатору
// @Tags люди
// @Param id path string true "Идентификатор человека"
// @Success 204 "No Content"
// @Router /v1/human/{id} [delete]
func (r *humanRoutes) deleteHumanByID(c *gin.Context) {
	userID := c.Param("id")

	err := r.humanService.DeleteByID(c, userID)
	if err != nil {
		if errors.Is(err, service.ErrHumanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": service.ErrHumanNotFound.Error(),
			})
			return
		}

		r.log.Error("failed to delete human", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	r.log.Info("deleted user", slog.String("userID", userID))
	c.Status(http.StatusNoContent)
}

// @Summary Обновить человека по его идентификатору
// @Description Обновить человека по его идентификатору. Принимает JSON с обновленными полями
// @Tags люди
// @Accept json
// @Param id path string true "Идентификатор человека"
// @Param input body request.UpdateHuman true "Тело запроса"
// @Success 200 {object} entity.Human
// @Router /v1/human/{id} [patch]
func (r *humanRoutes) updateCarByID(c *gin.Context) {
	userID := c.Param("id")
	var updateData request.UpdateHuman

	if err := c.ShouldBindJSON(&updateData); err != nil {
		r.log.Error("Invalid JSON payload", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON payload",
		})
		return
	}

	updatedHuman, err := r.humanService.UpdateByID(c, userID, &service.HumanInput{
		Name:        updateData.Name,
		Surname:     updateData.Surname,
		Potronymic:  updateData.Potronymic,
		Age:         updateData.Age,
		Gender:      updateData.Gender,
		Nationality: updateData.Nationality,
	})
	if err != nil {
		if errors.Is(err, service.ErrHumanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": service.ErrHumanNotFound.Error(),
			})
			return
		}

		r.log.Error("failed to update human", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	r.log.Info("successfully updated human", slog.String("user id", userID))
	c.JSON(http.StatusOK, gin.H{
		"human": updatedHuman,
	})
}
