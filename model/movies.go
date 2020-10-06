package model

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/whatflix/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const moviesCollection = "movies"

func Get(title string, prefLangSF []string) ([]*entity.MoviesCollRecord, error) {
	var moviesRecords []*entity.MoviesCollRecord
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
		var movieRecord entity.MoviesCollRecord
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
