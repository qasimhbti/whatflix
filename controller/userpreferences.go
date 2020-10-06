package controller

import (
	"github.com/whatflix/entity"
	"github.com/whatflix/model"
)

type userPreferencesManagerImpl struct{}

func (m *userPreferencesManagerImpl) get(userID int32) (*entity.UserPreferences, error) {
	userPrefs, err := model.UserPrefsGet(userID)
	if err != nil {
		return nil, err
	}
	return userPrefs, nil
}

func (m *userPreferencesManagerImpl) getAll() ([]*entity.UserPreferences, error) {
	userPrefs, err := model.UserPrefsGetAll()
	if err != nil {
		return nil, err
	}
	return userPrefs, nil
}
