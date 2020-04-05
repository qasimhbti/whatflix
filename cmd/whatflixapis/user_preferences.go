package main

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const userPreferencesCollection = "user_preferences"

// UserPreferences -- schema of user_preferences collection
type userPreferences struct {
	UserID             int      `json:"userid" bson:"userid"`
	FavouriteActors    []string `json:"favourite_actors" bson:"favourite_actors"`
	PreferredLanguages []string `json:"preferred_languages" bson:"preferred_languages"`
	FavouriteDirectors []string `json:"favourite_directors" bson:"favourite_directors"`
	PrefLangShortForm  []string `json:"pref_lang_short_form"`
}

type userPreferencesManagerImpl struct{}

func (m *userPreferencesManagerImpl) get(userID int, db *mongo.Database) (*userPreferences, error) {
	var userPref *userPreferences
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userid": userID}
	err := db.
		Collection(userPreferencesCollection).
		FindOne(ctx, filter).
		Decode(&userPref)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("userID : %d does not found", userID)
			return nil, nil
		}
		return nil, errors.WithMessage(err, "user preference")
	}

	_ = userPref.creatShortAbrevation()
	return userPref, nil
}

func (m *userPreferencesManagerImpl) getAll(db *mongo.Database) ([]*userPreferences, error) {
	var userPrefs []*userPreferences
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	cursor, err := db.
		Collection(userPreferencesCollection).
		Find(ctx, filter)
	if err != nil {
		return nil, errors.WithMessage(err, "get all")
	}

	for cursor.Next(ctx) {
		var data userPreferences
		err := cursor.Decode(&data)
		if err != nil {
			log.Printf("error while decoding : %v", err)
			continue
		}
		_ = data.creatShortAbrevation()
		userPrefs = append(userPrefs, &data)
	}
	defer cursor.Close(context.TODO())
	return userPrefs, nil

}

func (u *userPreferences) creatShortAbrevation() error {
	for _, lang := range u.PreferredLanguages {
		sf, ok := languagesShortForm[lang]
		if !ok {
			log.Printf("short form of language :%s not found", lang)
			continue
		}
		u.PrefLangShortForm = append(u.PrefLangShortForm, sf)
	}
	return nil
}
