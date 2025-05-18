package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/P1xart/effective_mobile_service/internal/entity"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

func limit(c *gin.Context) (int, error) {
	l := c.Query("limit")
	if l == "" {
		return defaultLimit, nil
	}

	return strconv.Atoi(l)
}

func offset(c *gin.Context) (int, error) {
	l := c.Query("offset")
	if l == "" {
		return defaultOffset, nil
	}

	return strconv.Atoi(l)
}

func checkLimitOffset(limit, offset int) bool {
	return limit >= 1 && offset >= 0
}

func buildFilters(c *gin.Context) *entity.HumanFilters {
	limit, errLimit := limit(c)
	offset, errOffset := offset(c)
	ok := checkLimitOffset(limit, offset)
	if !ok || errLimit != nil || errOffset != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Bad limit: %d, offest: %d", limit, offset),
		})
		return nil
	}

	ageFromStr := c.Query("age_from")
	ageToStr := c.Query("age_to")
	genderStr := c.Query("gender")
	nationalyStr := c.Query("nationaly")
	queryUrl := c.Request.URL.RawQuery

	var (
		ageFrom, ageTo uint64
		err error
	)
	if ageFromStr != "" {
		ageFrom, err = strconv.ParseUint(ageFromStr, 10, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to parse filter age from",
			})
			return nil
		}
	}

	if ageToStr != "" {
		ageTo, err = strconv.ParseUint(ageToStr, 10, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to parse filter age to",
			})
			return nil
		}
	}

	var genders []string
	if genderStr != "" {
		gendersArray := strings.Split(genderStr, ",")
		genders = append(genders, gendersArray...)

	}

	var nationalites []string
	if nationalyStr != "" {
		nationalitesArray := strings.Split(nationalyStr, ",")
		nationalites = append(nationalites, nationalitesArray...)

	}

	filters := entity.HumanFilters{
		Limit:  uint64(limit),
		Offset: uint64(offset),

		AgeFrom:   uint8(ageFrom),
		AgeTo:     uint8(ageTo),
		Gender:    genders,
		Nationaly: nationalites,

		QueryUrl: queryUrl,
	}

	return &filters
}
