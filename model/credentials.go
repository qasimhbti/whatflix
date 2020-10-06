package model

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/whatflix/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const signinCredentialCollection = "users"

func GetCredentials(userName string) (*entity.SigninCred, error) {
	var cred *entity.SigninCred
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"username": userName}

	coll := db.Collection(signinCredentialCollection)
	singleRes := coll.FindOne(ctx, filter)
	err := singleRes.Decode(&cred)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("userName : %s does not found", userName)
		}
		return nil, errors.WithMessage(err, "signin credentials")
	}
	return cred, nil
}
