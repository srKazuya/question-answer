//Pacakge auth provides
package auth

type User struct {
	ID           uint64 `json:"id"`
	Username     string `json:"username" validate:"required,min=3,max=32"`
	PasswordHash string `json:"-"`
}
