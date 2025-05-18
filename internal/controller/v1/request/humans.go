package request

type CreateHuman struct {
	Name       string `json:"name" validate:"required" example:"Vasiliy"`
	Surname    string `json:"surname" validate:"required" example:"Vasiliev"`
	Potronymic string `json:"potronymic,omitempty" example:"Vasilevich"`
}

type UpdateHuman struct {
	Name        string `json:"name,omitempty" example:"Igor"`
	Surname     string `json:"surname,omitempty" example:"Igorev"`
	Potronymic  string `json:"potronymic,omitempty" example:"Igorevich"`
	Age         int    `json:"age,omitempty" example:"22"`
	Gender      string `json:"gender,omitempty" example:"male"`
	Nationality string `json:"nationality,omitempty" example:"US"`
}
