package entity

type Role struct {
	RoleId    string `json:"role_id" db:"pk"`
	Name      string `json:"name"`
	IsActive  int    `json:"is_active"`
	IsDeleted int    `json:"is_deleted"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

type RoleAccess struct {
	RoleAccessId        string `json:"role_access_id" db:"pk"`
	RoleId              string `json:"role_id"`
	DivisionId          string `json:"division_id"`
	AbleToManageTask    int    `json:"able_to_manage_task"`
	AsResponsiblePerson int    `json:"as_responsible_person"`
}

type ViewRole struct {
	RoleId    string       `json:"role_id" db:"pk"`
	Name      string       `json:"name"`
	Access    []RoleAccess `json:"access"`
	IsActive  int          `json:"is_active"`
	IsDeleted int          `json:"is_deleted"`
	CreatedBy string       `json:"created_by"`
	CreatedAt string       `json:"created_at"`
}
