package controller

import (
	"github.com/whatflix/entity"
	"github.com/whatflix/model"
)

type creditsManagerImpl struct{}

func (m *creditsManagerImpl) get(searchTexts *entity.SearchText) ([]*entity.CreditsData, error) {
	creditsDatas, err := model.CreditsGet(searchTexts)
	if err != nil {
		return nil, err
	}
	return creditsDatas, nil
}
