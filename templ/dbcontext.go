package templates

import (
	"fmt"

	"github.com/mates182/ginshot/models"
)

func GetDBContextTemplate(config *models.ProjectConfig, db string, name string) string {
	var dbContextTemplate string

	switch db {
	case "mongo":
		dbContextTemplate = fmt.Sprintf(`
package dbcontext

import (
	"context"
	"` + config.ProjectName + `/secrets"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var clientInstance *mongo.Client
var clientOnce sync.Once

func GetDBClient() *mongo.Client {
	clientOnce.Do(func() {
		endpoint := secrets.Get%sDBURI()
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(endpoint))
		if err != nil {
			panic(err)
		}
		fmt.Println("Connected to %s Database Server")

		err = client.Ping(context.Background(), readpref.Primary())
		if err != nil {
			panic(err)
		}
		fmt.Println("Pong")
		clientInstance = client
	})

	return clientInstance
}
		`, name, name)
	case "redis":
		dbContextTemplate = `
package dbcontext

import (
	"context"
	"` + config.ProjectName + `/secrets"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func GetDBClient() *redis.Client {

	dbURI := secrets.Get` + name + `DBURI()
	dbPassword := secrets.Get` + name + `DBPassword()
	dbOptions := &redis.Options{
		Addr: dbURI,
		DB:   0,
	}
	if dbPassword != "" {
		dbOptions.Password = dbPassword
	}

	client := redis.NewClient(dbOptions)
	ping, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to DB: %s\n", err.Error())
		return nil
	}
	fmt.Printf("Ping: %s\n", ping)
	return client
}
`
default:
	fmt.Println("Unsupported database type.")
	
}
return dbContextTemplate
}
