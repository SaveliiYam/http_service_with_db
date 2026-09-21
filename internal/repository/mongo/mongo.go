package mongo

import (
	"context"
	"time"

	"github.com/SaveliiYam/http_service_with_db/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	userCollection    = "users"
	counterCollection = "counters"
)

type MongoRepository struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoRepository() (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		options.Client().ApplyURI("mongodb://localhost:27017"),
	)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoRepository{
		client: client,
		db:     client.Database("my_app"),
	}, nil
}

func (r *MongoRepository) GetByID(
	ctx context.Context,
	id int,
) (models.User, error) {
	var user userDocument

	if err := r.db.Collection(userCollection).FindOne(ctx, bson.M{"id": id}).Decode(&user); err != nil {
		return models.User{}, err
	}

	return user.toModels(), nil
}

func (r *MongoRepository) CreateUser(
	ctx context.Context,
	name string,
) (models.User, error) {
	var counter int

	err := r.db.Collection(counterCollection).FindOneAndUpdate(
		ctx,
		bson.M{"_id": "user_id"},         // фильтр
		bson.M{"$inc": bson.M{"seq": 1}}, // атомарный инкремент
		options.FindOneAndUpdate().
			SetUpsert(true).                  // создаст, если нет
			SetReturnDocument(options.After), // вернёт документ ПОСЛЕ инкремента
	).Decode(&counter)
	if err != nil {
		return models.User{}, err
	}

	user := userDocument{
		ID:   counter,
		Name: name,
	}

	_, err = r.db.Collection(userCollection).InsertOne(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return user.toModels(), nil
}

type userDocument struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
}

func (u userDocument) toModels() models.User {
	return models.User{
		ID:   u.ID,
		Name: u.Name,
	}
}
