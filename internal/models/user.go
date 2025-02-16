package models

type User struct {
	ID          string `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	ActiveVault string `json:"activeVault"`
	CreatedAt   string `json:"createdAt"`
}
