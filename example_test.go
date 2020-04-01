package orm_test

import (
	"context"

	"github.com/phogolabs/orm"
	"github.com/phogolabs/orm/dialect/sql"
	"github.com/phogolabs/parcello"
)

func ExampleGateway_First() {
	gateway, err := orm.Open("sqlite3", "example.db")
	if err != nil {
		panic(err)
	}

	user := &User{}
	query := orm.SQL("SELECT * FROM users ORDER BY created_at")

	if err := gateway.First(context.TODO(), query, user); err != nil {
		panic(err)
	}
}

func ExampleGateway_Only() {
	gateway, err := orm.Open("sqlite3", "example.db")
	if err != nil {
		panic(err)
	}

	user := &User{}
	query := orm.SQL("SELECT * FROM users WHERE id = ?", "007")

	if err := gateway.Only(context.TODO(), query, user); err != nil {
		panic(err)
	}
}

func ExampleGateway_Exec() {
	gateway, err := orm.Open("sqlite3", "example.db")
	if err != nil {
		panic(err)
	}

	query :=
		sql.Insert("users").
			Columns("first_name", "last_name").
			Values("John", "Doe").
			Returning("id")

	if _, err := gateway.Exec(context.TODO(), query); err != nil {
		panic(err)
	}
}

func ExampleRoutine() {
	gateway, err := orm.Open("sqlite3", "example.db")
	if err != nil {
		panic(err)
	}

	if err = gateway.ReadDir(parcello.Dir("./database/command")); err != nil {
		panic(err)
	}

	users := []*User{}
	routine := orm.Routine("show-top-5-users")

	if err := gateway.All(context.TODO(), routine, &users); err != nil {
		panic(err)
	}
}

func ExampleSQL() {
	gateway, err := orm.Open("sqlite3", "example.db")
	if err != nil {
		panic(err)
	}

	users := []*User{}
	query := orm.SQL("SELECT name FROM users")

	if err := gateway.All(context.TODO(), query, &users); err != nil {
		panic(err)
	}
}
