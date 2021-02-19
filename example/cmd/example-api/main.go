package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/bxcodec/faker/v3"
	"github.com/phogolabs/log"
	"github.com/phogolabs/orm/example/database"
	"github.com/phogolabs/orm/example/database/model"
)

func main() {
	repository, err := database.NewUserRepository("sqlite3://prana.db")
	if err != nil {
		log.WithError(err).Fatal("failed to open database connection")
	}
	// close the connection
	defer repository.Close()

	if err = generate(repository); err != nil {
		log.WithError(err).Fatal("failed to generate all records")
	}

	users, err := repository.AllUsers(context.TODO())
	if err != nil {
		log.WithError(err).Fatal("failed to select all records")
	}

	if err := show(users); err != nil {
		log.WithError(err).Fatal("printing the records")
	}
}

func generate(repository *database.UserRepository) error {
	ctx := context.TODO()

	for i := 0; i < 10; i++ {
		record := &model.User{
			ID:        int(time.Now().UnixNano()),
			FirstName: faker.FirstName(),
			LastName:  pointer(faker.LastName()),
		}

		if err := repository.InsertUser(ctx, record); err != nil {
			return err
		}
	}

	return nil
}

func show(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", " ")
	return encoder.Encode(data)
}

func pointer(v string) *string {
	return &v
}
