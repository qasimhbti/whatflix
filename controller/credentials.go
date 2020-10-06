package controller

import (
	"github.com/whatflix/entity"
	"github.com/whatflix/model"
)

type signinCredentialManagerImpl struct{}

func (m *signinCredentialManagerImpl) get(userName string) (*entity.SigninCred, error) {
	cred, err := model.GetCredentials(userName)
	if err != nil {
		return nil, err
	}
	return cred, nil
}
