package response

import "github.com/P1xart/effective_mobile_service/internal/entity"

type GetAllHumans struct {
	Humans []entity.Human `json:"humans"`
}
