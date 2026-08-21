package category

type Model struct {
	ID       int64  `json:"id"`
	Category string `json:"category"`
	IsActive bool   `json:"is_active"`
}
