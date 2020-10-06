package entity

type MoviesCollRecord struct {
	Title       string  `bson:"title"`
	Language    string  `bson:"original_language"`
	VoteAverage float64 `bson:"vote_average"`
	VoteCount   int     `bson:"vote_count"`
}
