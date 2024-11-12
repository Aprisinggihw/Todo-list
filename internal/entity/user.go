package entity

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
}

func (User) TableName() string {
	return "public.users"
}

type UserReg struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
}

func (UserReg) TableName() string {
	return "public.users"
}

