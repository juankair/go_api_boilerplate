package entity

type Document struct {
	ID         string `json:"id" db:"pk"`
	PatientId  string `json:"patient_id"`
	DocsNumber string `json:"docs_number"`
	Remarks    string `json:"remarks"`
	Needs      string `json:"needs"`
	IsActive   int    `json:"is_active"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
