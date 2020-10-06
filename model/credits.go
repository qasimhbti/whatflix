package model

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/whatflix/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const creditsCollection = "credits"

func CreditsGet(searchTexts *entity.SearchText) ([]*entity.CreditsData, error) {
	var creditsDatas []*entity.CreditsData

	opts := options.Find().SetSort(bson.D{
		{"title", 1},
	}).
		SetProjection(bson.D{
			{"title", 1},
			//{"movie_id", 1},
			{"_id", 0},
		})

	// search crdits collection for given search texts
	if len(searchTexts.Texts) != 0 {
		for _, text := range searchTexts.Texts {
			filter := bson.D{
				{"$or",
					bson.A{
						bson.D{{"title", text}},
						bson.D{{"cast", bson.D{{"$elemMatch", bson.D{{"name", text}}}}}},
						bson.D{{"crew", bson.D{{"$elemMatch", bson.D{{"job", "Director"}, {"name", text}}}}}},
					},
				},
			}
			results, err := executeQuery(filter, opts)
			if err != nil {
				log.Printf("execute query :: filter %v, options %v", filter, opts)
				continue
			}
			creditsDatas = append(creditsDatas, results...)
		}
		return creditsDatas, nil
	}
	//} else {
	if len(searchTexts.Texts) == 0 {
		// --- len text === 0
		// search crdits collection for fav actors
		for _, actor := range searchTexts.Actors {
			filter := bson.D{
				{"cast", bson.D{{"$elemMatch", bson.D{{"name", actor}}}}},
			}

			results, err := executeQuery(filter, opts)
			if err != nil {
				log.Printf("execute query :: filter %v, options %v", filter, opts)
				continue
			}

			creditsDatas = append(creditsDatas, results...)
		}

		// search crdits collection for fav directors
		for _, director := range searchTexts.Directors {
			filter := bson.D{
				{"crew", bson.D{{"$elemMatch", bson.D{{"job", "Director"}, {"name", director}}}}},
			}
			results, err := executeQuery(filter, opts)
			if err != nil {
				log.Printf("execute query :: filter %v, options %v", filter, opts)
				continue
			}
			creditsDatas = append(creditsDatas, results...)
		}
		return creditsDatas, nil
	}
	return creditsDatas, nil
}

func executeQuery(filter primitive.D, opts *options.FindOptions) ([]*entity.CreditsData, error) {
	var results []*entity.CreditsData
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection(creditsCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.WithMessage(err, "credits collection")
	}

	for cursor.Next(context.TODO()) {
		var data entity.CreditsData
		err := cursor.Decode(&data)
		if err != nil {
			log.Printf("error while decoding : %v", err)
			continue
		}
		results = append(results, &data)
	}
	defer cursor.Close(context.TODO())
	return results, nil
}
