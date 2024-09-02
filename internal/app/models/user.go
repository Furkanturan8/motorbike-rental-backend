package models

type UserRole int

const (
	UserRoleNormal UserRole = 1
	UserRoleAdmin  UserRole = 10
)

type User struct {
	BaseModel
	Name     string   `gorm:"column:name"`
	Surname  string   `gorm:"column:surname"`
	UserName string   `gorm:"column:username"`
	Email    string   `gorm:"column:email;unique"`
	Phone    string   `gorm:"column:phone;unique"`
	Password string   `gorm:"column:password"`
	Role     UserRole `gorm:"column:role"`
}

func (User) ModelName() string {
	return "user"
}

func (u User) String() string {
	return u.Name + " " + u.Surname
}

func (r UserRole) String() string {
	switch r {
	case UserRoleNormal:
		return "normal"
	case UserRoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}
