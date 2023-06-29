package db

import (
	"context"
	"gitlab.com/koshkaj/macler/backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Mongo struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func (m *Mongo) InsertCronJob(ctx context.Context, data interface{}) {
	col := m.Db.Collection("cron_jobs")
	session, err := m.Client.StartSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.EndSession(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}
		_, err := col.InsertOne(ctx, data)
		if err != nil {
			return err
		}
		err = session.CommitTransaction(ctx)
		if err != nil {
			return err
		}
		return nil
	})

}

func (m *Mongo) GetAllCronJobs(ctx context.Context) ([]types.CronMongoInput, error) {
	var data []types.CronMongoInput
	cur, err := m.Db.Collection("cron_jobs").Find(ctx, bson.M{})
	if err != nil {
		return data, err
	}
	err = cur.All(ctx, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func New(client *mongo.Client, db *mongo.Database) *Mongo {
	return &Mongo{Client: client, Db: db}
}

func InitDb(ctx context.Context) *Mongo {
	uri := os.Getenv("MONGO_URI")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient := New(client, client.Database("macler_main"))
	return mongoClient
}
