package request

type CreateHuman struct {
	Name       string `json:"name" validate:"required" example:"Vasiliy"`
	Surname    string `json:"surname" validate:"required" example:"Vasiliev"`
	Potronymic string `json:"potronymic,omitempty" example:"Vasilevich"`
}

type UpdateHuman struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Potronymic  string `json:"potronymic"`
	Age         string `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
