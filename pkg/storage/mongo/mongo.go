package mongo

import (
	"context"
	"news/pkg/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	client *mongo.Client
	dbName string
}

func New(uri string) (*Store, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}
	return &Store{
		client: client,
		dbName: "news",
	}, nil
}

func (s *Store) getPostsCollection() *mongo.Collection {
	return s.client.Database(s.dbName).Collection("posts")
}

func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.getPostsCollection()
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var posts []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, cur.Err()
}

func (s *Store) AddPost(p storage.Post) error {
	collection := s.getPostsCollection()
	_, err := collection.InsertOne(context.Background(), p)
	return err
}

func (s *Store) UpdatePost(p storage.Post) error {
	collection := s.getPostsCollection()
	filter := bson.M{"id": p.ID}
	update := bson.M{"$set": bson.M{
		"authorname": p.Author,
		"title":      p.Title,
		"content":    p.Content,
		"created_at": p.CreatedAt,
	}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (s *Store) DeletePost(p storage.Post) error {
	collection := s.getPostsCollection()
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": p.ID})
	return err
}
