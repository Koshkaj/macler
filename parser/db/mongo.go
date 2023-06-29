package db

import (
	"context"
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

func (m *Mongo) QuerySeenIDs(ctx context.Context, filters bson.M) map[string]struct{} { // todo add data here as parameter
	col := m.Db.Collection("seen")
	cursor, err := col.Find(ctx, filters)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	seen_ids := make(map[string]struct{})
	var result bson.M
	for cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		seen_ids[result["productid"].(string)] = struct{}{}
	}
	return seen_ids
}

func (m *Mongo) InsertSeenIDs(ctx context.Context, data []interface{}) {
	col := m.Db.Collection("seen")
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
		_, err := col.InsertMany(ctx, data)
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
