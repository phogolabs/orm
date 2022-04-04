package mock

import "time"

type Group struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type User struct {
	ID        string    `db:"id,primary_key,immutable"`
	Name      string    `db:"name"`
	Email     *string   `db:"email"`
	Group     *Group    `db:"group,foreign_key=group_id,reference_key=id"`
	CreatedAt time.Time `db:"created_at,auto,read_only"`
	UpdatedAt time.Time `db:"created_at,auto,read_only"`
}
