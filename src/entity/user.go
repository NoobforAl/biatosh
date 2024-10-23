package entity

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}
