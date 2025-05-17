package entity

type Human struct {
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Potronymic string `json:"potronymic,omitempty" db:"potronymic"`
	Age        uint8  `json:"age" db:"age"`
	Gender     string `json:"Gender" db:"gender"`
	Nationaly  string `json:"nationaly" db:"nationaly"`
}

type AgeResp struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type GenderResp struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type NationalyResp struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []Country `json:"country"`
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type HumanFilters struct {
	Limit  uint64
	Offset uint64

	AgeFrom   uint8
	AgeTo     uint8
	Gender    []string
	Nationaly []string

	QueryUrl string
}
