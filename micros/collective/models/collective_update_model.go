package models

type CollectiveUpdateModel struct {
	Name        string `json:"Name" bson:"Name"`
	Avatar      string `json:"avatar" bson:"avatar"`
	Banner      string `json:"banner" bson:"banner"`
	LastUpdated int64  `json:"last_updated" bson:"last_updated"`
}

type CollectiveGeneralUpdateModel struct {
	Avatar      string `json:"avatar" bson:"avatar"`
	Banner      string `json:"banner" bson:"banner"`
	Name        string `json:"Name" bson:"Name"`
	LastUpdated int64  `json:"last_updated" bson:"last_updated"`
}
