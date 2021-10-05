package models

type CollectiveUpdateModel struct {
	Name        string `json:"Name" bson:"Name"`
	Avatar      string `json:"avatar" bson:"avatar"`
	Banner      string `json:"banner" bson:"banner"`
	Tagline     string `json:"tagLine" bson:"tagLine"`
	Title       string `json:"title" bson:"title"`
	Body        string `json:"body" bson:"body"`
	CreatedDate int64  `json:"created_date" bson:"created_date"`
	LastUpdated int64  `json:"last_updated" bson:"last_updated"`
}

type CollectiveGeneralUpdateModel struct {
	Avatar      string `json:"avatar" bson:"avatar"`
	Banner      string `json:"banner" bson:"banner"`
	Name        string `json:"Name" bson:"Name"`
	Tagline     string `json:"tagLine" bson:"tagLine"`
	Title       string `json:"title" bson:"title"`
	Body        string `json:"body" bson:"body"`
	LastUpdated int64  `json:"last_updated" bson:"last_updated"`
}
