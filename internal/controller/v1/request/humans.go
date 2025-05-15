package request

type CreateHuman struct {
	Name       string `json:"name" validate:"required" example:"Vasiliy"`
	Surname    string `json:"surname" validate:"required" example:"Vasiliev"`
	Potronymic string `json:"potronymic,omitempty" example:"Vasilevich"`
}
