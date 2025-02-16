package models

type VaultRole string

const (
	VaultRoleOwner VaultRole = "owner"
)

type Vault struct {
	ID        string `json:"ID"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type UserVaultWithRole struct {
	ID       string    `json:"ID"`
	Name     string    `json:"name"`
	UserRole VaultRole `json:"userRole"`
}
