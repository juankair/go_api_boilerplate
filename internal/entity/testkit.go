package entity

type TestKit struct {
	ID          int     `json:"id" db:"pk"`
	Testkit     string  `json:"testkit"`
	Hasil       string  `json:"hasil"`
	Deleted     int     `json:"deleted"`
	CreatedBy   *string `json:"created_by"`
	CreatedDate *string `json:"created_date"`
	UpdatedBy   *string `json:"updated_by"`
	UpdatedDate *string `json:"updated_date"`
}
