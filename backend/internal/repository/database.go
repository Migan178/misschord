package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/Migan178/misschord-backend/internal/configs"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	"github.com/Migan178/misschord-backend/internal/repository/ent/migrate"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	client   *ent.Client
	Users    *UserRepository
	Messages *MessageRepository
}

var instance *Database
var once sync.Once

func GetDatabase() *Database {
	once.Do(func() {
		config := configs.GetConfig()
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", config.Database.Username, config.Database.Password, config.Database.Hostname, config.Database.Port, config.Database.Name)

		client, err := ent.Open("mysql", connectionString)
		if err != nil {
			panic(err)
		}

		if err := client.Schema.Create(context.Background(),
			migrate.WithDropColumn(true),
			migrate.WithDropIndex(true),
			migrate.WithForeignKeys(true),
		); err != nil {
			fmt.Printf("failed creating schema resources: %v\n", err)
		}

		instance = &Database{
			client:   client,
			Users:    newUserRepository(client),
			Messages: newMessageRepository(client),
		}
	})

	return instance
}
