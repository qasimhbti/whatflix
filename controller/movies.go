package controller

import (
	"github.com/pkg/errors"
	"github.com/whatflix/entity"
	"github.com/whatflix/model"
)

type moviesManagerImpl struct{}

func (m *moviesManagerImpl) get(title string, prefLangSF []string) ([]*entity.MoviesCollRecord, error) {
	moviesRecords, err := model.Get(title, prefLangSF)
	if err != nil {
		return nil, errors.WithMessage(err, "movies collection")
	}
	return moviesRecords, nil
}

// ByVoteAverage --sort by vote average of given movie.
type ByVoteAverage []*entity.MoviesCollRecord

func (v ByVoteAverage) Len() int           { return len(v) }
func (v ByVoteAverage) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByVoteAverage) Less(i, j int) bool { return v[i].VoteAverage > v[j].VoteAverage }

// ByTitle --sort by title of given movie.
type ByTitle []*entity.MoviesCollRecord

func (t ByTitle) Len() int           { return len(t) }
func (t ByTitle) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTitle) Less(i, j int) bool { return t[i].Title < t[j].Title }
