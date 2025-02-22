package models

type ExpenseCategoryStatus string

var (
	ExpenseCategoryStatusActive   ExpenseCategoryStatus = "active"
	ExpenseCategoryStatusInactive ExpenseCategoryStatus = "inactive"
)

type ExpenseCategory struct {
	ID        string                `json:"id"`
	Name      string                `json:"name"`
	Status    ExpenseCategoryStatus `json:"status"`
	Priority  int                   `json:"priority"`
	VaultID   string                `json:"vaultID"`
	CreatedBy string                `json:"createdBy"`
	CreatedAt string                `json:"createdAt"`
}
