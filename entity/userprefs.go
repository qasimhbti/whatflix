package entity

//UserPreferences -- schema of user_preferences collection
type UserPreferences struct {
	UserID             int32    `json:"userid" bson:"userid"`
	FavouriteActors    []string `json:"favourite_actors" bson:"favourite_actors"`
	PreferredLanguages []string `json:"preferred_languages" bson:"preferred_languages"`
	FavouriteDirectors []string `json:"favourite_directors" bson:"favourite_directors"`
	PrefLangShortForm  []string `json:"pref_lang_short_form"`
}
