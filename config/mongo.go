package config

import (
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"

	"time"
)

var DB *mongo.Client

func InitMongo() {

	// Check if the environment variable is set
	if viper.Get("MONGO_URI") == "" {
		log.Fatal("MONGO_URI is not set")
	}

	// choose mongo local or atlas
	mongoEnvironment := "development"
	if Viper.GetString("APP_ENV") == "production" {
		mongoEnvironment = "atlas"
	}

	log.Println("Connecting to mongo " + mongoEnvironment + "...")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFunc()

	counts := 0

	for {
		LoadConfig()
		uri := Viper.Get("MONGO_URI").(string)
		connect, err := OpenDB(ctx, uri)

		if counts >= 10 {
			log.Fatal("Failed to connect to mongodb in " + mongoEnvironment + " mode")
		}

		counts++

		if err != nil {
			log.Printf("Failed to connect to mongodb in %s mode, trying again in 5 seconds : %v, count: %d/10", mongoEnvironment, err, counts)
			time.Sleep(5 * time.Second)
			continue
		}

		err = connect.Ping(ctx, nil)
		if err != nil {
			log.Printf("Failed to ping mongodb in %s mode, trying again in 5 seconds : %v, count: %d/10", mongoEnvironment, err, counts)
			time.Sleep(5 * time.Second)
			continue
		}

		DB = connect
		log.Printf("Connected to mongodb %s", mongoEnvironment)
		break
	}

}

func OpenDB(ctx context.Context, uri string) (*mongo.Client, error) {
	connect, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = connect.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return connect, nil
}
