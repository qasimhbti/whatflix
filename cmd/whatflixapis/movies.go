package main

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const moviesCollection = "movies"

type moviesManagerImpl struct{}

type moviesCollRecord struct {
	Title       string  `bson:"title"`
	Language    string  `bson:"original_language"`
	VoteAverage float64 `bson:"vote_average"`
	VoteCount   int     `bson:"vote_count"`
}

func (m *moviesManagerImpl) get(title string, prefLangSF []string, db *mongo.Database) ([]*moviesCollRecord, error) {
	var moviesRecords []*moviesCollRecord
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetProjection(bson.D{
		{"_id", 0},
		{"title", 1},
		{"original_language", 1},
		{"vote_average", 1},
		{"vote_count", 1},
	})

	filter := bson.D{
		{"title", title},
	}

	cursor, err := db.
		Collection(moviesCollection).
		Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.WithMessage(err, "movies collection")
	}
	for cursor.Next(context.TODO()) {
		var movieRecord moviesCollRecord
		err := cursor.Decode(&movieRecord)
		if err != nil {
			log.Printf("error while decoding %v", err)
			continue
		}

		for _, lang := range prefLangSF {
			if movieRecord.Language == lang {
				moviesRecords = append(moviesRecords, &movieRecord)
			}
		}
	}
	defer cursor.Close(context.TODO())
	return moviesRecords, nil
}

// ByVoteAverage --sort by vote average of given movie.
type ByVoteAverage []*moviesCollRecord

func (v ByVoteAverage) Len() int           { return len(v) }
func (v ByVoteAverage) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByVoteAverage) Less(i, j int) bool { return v[i].VoteAverage > v[j].VoteAverage }

// ByTitle --sort by title of given movie.
type ByTitle []*moviesCollRecord

func (t ByTitle) Len() int           { return len(t) }
func (t ByTitle) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTitle) Less(i, j int) bool { return t[i].Title < t[j].Title }
