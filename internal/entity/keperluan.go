package entity

type Keperluan struct {
	ID          int     `json:"id" db:"pk"`
	Keperluan   string  `json:"keperluan"`
	Note        *string `json:"note"`
	Deleted     int     `json:"deleted"`
	CreatedBy   *string `json:"created_by"`
	CreatedDate *string `json:"created_date"`
	UpdatedBy   *string `json:"updated_by"`
	UpdatedDate *string `json:"updated_date"`
}
