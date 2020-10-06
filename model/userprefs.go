package model

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/whatflix/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const userPreferencesCollection = "user_preferences"

func UserPrefsGet(userID int32) (*entity.UserPreferences, error) {
	var userPref *entity.UserPreferences
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

	modUserPref := creatShortAbrevation(*userPref)
	modUserPref.UserID = userID
	return &modUserPref, nil
}

func UserPrefsGetAll() ([]*entity.UserPreferences, error) {
	var userPrefs []*entity.UserPreferences
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	cursor, err := db.
		Collection(userPreferencesCollection).
		Find(ctx, filter)
	if err != nil {
		return nil, errors.WithMessage(err, "model:Get All")
	}

	for cursor.Next(ctx) {
		var data entity.UserPreferences
		err := cursor.Decode(&data)
		if err != nil {
			log.Printf("error while decoding : %v", err)
			continue
		}
		modData := creatShortAbrevation(data)
		userPrefs = append(userPrefs, &modData)
	}
	defer cursor.Close(context.TODO())
	return userPrefs, nil
}

func creatShortAbrevation(data entity.UserPreferences) entity.UserPreferences {
	for _, lang := range data.PreferredLanguages {
		sf, ok := languagesShortForm[lang]
		if !ok {
			log.Printf("short form of language :%s not found", lang)
			continue
		}
		data.PrefLangShortForm = append(data.PrefLangShortForm, sf)
	}
	return data
}
