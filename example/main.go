package main

import (
	"fmt"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	_ "github.com/mattn/go-sqlite3"
	"github.com/phogolabs/orm/example/database"
	"github.com/phogolabs/parcello"
	"github.com/phogolabs/schema"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/apex/log"
	"github.com/phogolabs/orm"
	"github.com/phogolabs/orm/example/database/model"
)

func main() {
	repository, err := connect()
	if err != nil {
		log.WithError(err).Fatal("Failed to open database connection")
	}
	defer repository.Gateway.Close()

	if err = random(repository); err != nil {
		log.WithError(err).Fatal("Failed to generate users")
	}

	users, err := repository.SelectAll()
	if err != nil {
		log.WithError(err).Fatal("Failed to select all users")
	}

	validation := validator.New()

	if err := validation.Struct(users[0]); err != nil {
		panic(err)
	}

	show(users)
}

func connect() (*database.UserRepository, error) {
	gateway, err := orm.Connect("sqlite3://prana.db")
	if err != nil {
		return nil, err
	}

	if err = gateway.Migrate(parcello.ManagerAt("migration")); err != nil {
		return nil, err
	}

	if err = gateway.ReadDir(parcello.ManagerAt("routine")); err != nil {
		return nil, err
	}

	repository := &database.UserRepository{
		Gateway: gateway,
	}

	return repository, nil
}

func random(repository *database.UserRepository) error {
	for i := 0; i < 10; i++ {
		var lastName string

		if i%2 == 0 {
			lastName = randomdata.LastName()
		}

		user := &model.User{
			ID:        int(time.Now().UnixNano()),
			FirstName: randomdata.FirstName(randomdata.Male),
			LastName:  schema.NullStringFrom(lastName),
		}

		if err := repository.Insert(user); err != nil {
			return err
		}
	}

	return nil
}

func show(users []*model.User) {
	validate := validator.New()

	for _, user := range users {
		if err := validate.Struct(user); err != nil {
			log.WithError(err).Error("Failed to validate user")
			continue
		}

		fmt.Printf("User ID: %v\n", user.ID)
		fmt.Printf("First Name: %v\n", user.FirstName)

		if user.LastName.Valid {
			fmt.Printf("Last Name: %v\n", user.LastName.String)
		} else {
			fmt.Println("Last Name: null")
		}

		fmt.Println("---")
	}

}
