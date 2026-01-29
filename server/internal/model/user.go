package model

type User struct {
	ID        int64
	Email     string
	Password  string
	Name      string
	Status    string // PENDING, APPROVED, SUSPENDED
	IsAdmin   bool
	CreatedAt string
	UserSalt  string
}
