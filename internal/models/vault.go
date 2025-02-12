package models

type Vault struct {
	ID        string `json:"ID"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type UserVaultWithRole struct {
	ID       string `json:"ID"`
	Name     string `json:"name"`
	UserRole string `json:"userRole"`
}
