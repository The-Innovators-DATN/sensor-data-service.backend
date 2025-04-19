package domain

type Dashboard struct {
	ID          int32
	Name        string
	Description string
	LayoutJSON  string
	CreatedBy   int32
	CreatedAt   string
	UpdatedAt   string
	Version     int32
	Status      string
}
