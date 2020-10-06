package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var db DB

type DB interface {
	Collection(string) *mongo.Collection
	Find(context.Context, bson.M) (*mongo.Cursor, error)
	FindOne(context.Context, bson.M) *mongo.SingleResult
	//Decode(interface{}) error
	//QueryRow(string, ...interface{}) Row
	//Exec(string, ...interface{}) (Result, error)
}

type Collection interface {
}

/*type Row interface {
	Scan(...interface{}) error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}*/

type mongoDB struct {
	db *mongo.Database
}

func (m mongoDB) Collection(collName string) *mongo.Collection {
	return m.db.Collection(collName)
}

func (m mongoDB) Find(ctx context.Context, filter bson.M) (*mongo.Cursor, error) {
	return m.db.Collection("").Find(ctx, filter)
}

func (m mongoDB) FindOne(ctx context.Context, filter bson.M) *mongo.SingleResult {
	return m.db.Collection("").FindOne(ctx, filter)
}

/*func (m mongoDB) Decode() *mongo.SingleResult {
	return m.db.SingleResult().Decode("")
}*/

/*func (m mongoDB) QueryRow(query string, args ...interface{}) Row {
	//return m.db.QueryRow(query, args...)
}

func (m mongoDB) Exec(query string, args ...interface{}) (Result, error) {
	//return m.db.Exec(query, args...)
}*/

func SetDatabase(database *mongo.Database) {
	db = &mongoDB{database}
}
