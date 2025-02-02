package entity

type Pekerjaan struct {
	ID          int     `json:"id" db:"pk"`
	Pekerjaan   string  `json:"pekerjaan"`
	Note        *string `json:"note"`
	Deleted     int     `json:"deleted"`
	CreatedBy   *string `json:"created_by"`
	CreatedDate *string `json:"created_date"`
	UpdatedBy   *string `json:"updated_by"`
	UpdatedDate *string `json:"updated_date"`
}
