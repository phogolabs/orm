package main

import (
	"context"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
	"github.com/phogolabs/parcello"

	"github.com/apex/log"
	"github.com/phogolabs/orm"
	"github.com/phogolabs/orm/_example/database"
	"github.com/phogolabs/orm/_example/database/model"
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

	users, err := repository.AllUsers(context.TODO())
	if err != nil {
		log.WithError(err).Fatal("Failed to select all users")
	}

	spew.Dump(users)
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
		var (
			firstName = randomdata.FirstName(randomdata.Male)
			lastName  = randomdata.LastName()
		)

		user := &model.User{
			ID:        int(time.Now().UnixNano()),
			FirstName: firstName,
			LastName:  &lastName,
		}

		if err := repository.InsertUser(context.TODO(), user); err != nil {
			return err
		}
	}

	return nil
}
