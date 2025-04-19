package parameter

import "time"

type Parameter struct {
	ID             int       `db:"id"`
	Name           string    `db:"name"`
	Unit           string    `db:"unit"`
	ParameterGroup string    `db:"parameter_group"`
	Status         string    `db:"status"` // enum: active, inactive, deleted
	Description    string    `db:"description"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
