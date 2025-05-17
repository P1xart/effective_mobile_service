package v1

import (
	"log/slog"
	"net/http"


	"github.com/P1xart/effective_mobile_service/internal/config"
	"github.com/P1xart/effective_mobile_service/internal/controller/v1/request"
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
}

// @Summary Создание нового человека
// @Description Создание нового человека
// @Tags люди
// @Accept json
// @Param input body request.CreateHuman true "Тело запроса"
// @Success 201
// @Router /v1/human/ [post]
func (r *humanRoutes) createNewHuman(c *gin.Context) {
	var human request.CreateHuman

	err := c.ShouldBindJSON(&human); if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = r.valid.Struct(&human); if err != nil {
		r.log.Info("error validating data", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = r.humanService.CreateHuman(c, &service.CreateHuman{
		Name: human.Name,
		Surname: human.Surname,
		Potronymic: human.Potronymic,
	}, config.ApiUrls{}); if err != nil {
		r.log.Error("failed to create human", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}