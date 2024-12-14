package entity

type Account struct {
	AccountId    string `json:"account_id" db:"pk"`
	RoleId       string `json:"role_id"`
	EmployeeCode string `json:"employee_code"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Password     string `json:"password"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	IsActive     int    `json:"is_active"`
	IsDeleted    int    `json:"is_deleted"`
	CreatedAt    string `json:"created_at"`
}

type AccountMinimalData struct {
	AccountId    string `json:"account_id"`
	RoleId       string `json:"role_id"`
	RoleName     string `json:"role_name"`
	EmployeeCode string `json:"employee_code"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Token        string `json:"token"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	IsActive     int    `json:"is_active"`
	CreatedAt    string `json:"created_at"`
}

func (u Account) GetID() string {
	return u.AccountId
}

func (u Account) GetName() string {
	return u.FullName
}
