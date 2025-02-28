package dto

import (
	"github.com/GMcD/telar-web/constants"
	uuid "github.com/gofrs/uuid"
)

type Location struct {
	GeoJSONType string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

func NewLocation(lat, long float64) Location {
	return Location{
		"Point",
		[]float64{lat, long},
	}
}

type UserProfile struct {
	ObjectId       uuid.UUID                     `json:"objectId" bson:"objectId"`
	FullName       string                        `json:"fullName" bson:"fullName"`
	SocialName     string                        `json:"socialName" bson:"socialName"`
	Avatar         string                        `json:"avatar" bson:"avatar"`
	Banner         string                        `json:"banner" bson:"banner"`
	TagLine        string                        `json:"tagLine" bson:"tagLine"`
	CreatedDate    int64                         `json:"created_date" bson:"created_date"`
	LastUpdated    int64                         `json:"last_updated" bson:"last_updated"`
	LastSeen       int64                         `json:"lastSeen" bson:"lastSeen"`
	Email          string                        `json:"email" bson:"email"`
	Birthday       int64                         `json:"birthday" bson:"birthday"`
	WebUrl         string                        `json:"webUrl" bson:"webUrl"`
	Country        string                        `json:"country" bson:"country"`
	Address        string                        `json:"address" bson:"address"`
	School         string                        `json:"school" bson:"school"`
	LiveLocation   Location                      `json:"liveLocation" bson:"liveLocation"`
	Phone          string                        `json:"phone" bson:"phone"`
	Lang           string                        `json:"lang" bson:"lang"`
	CompanyName    string                        `json:"companyName" bson:"companyName"`
	VoteCount      int64                         `json:"voteCount" bson:"voteCount"`
	ShareCount     int64                         `json:"shareCount" bson:"shareCount"`
	FollowCount    int64                         `json:"followCount" bson:"followCount"`
	FollowerCount  int64                         `json:"followerCount" bson:"followerCount"`
	PostCount      int64                         `json:"postCount" bson:"postCount"`
	FacebookId     string                        `json:"facebookId" bson:"facebookId"`
	InstagramId    string                        `json:"instagramId" bson:"instagramId"`
	TwitterId      string                        `json:"twitterId" bson:"twitterId"`
	LinkedInId     string                        `json:"linkedInId" bson:"linkedInId"`
	AccessUserList []string                      `json:"accessUserList" bson:"accessUserList"`
	Permission     constants.UserPermissionConst `json:"permission" bson:"permission"`
}

type UserPublicProfile struct {
	ObjectId       uuid.UUID                     `json:"objectId" bson:"objectId"`
	FullName       string                        `json:"fullName" bson:"fullName"`
	SocialName     string                        `json:"socialName" bson:"socialName"`
	Avatar         string                        `json:"avatar" bson:"avatar"`
	Banner         string                        `json:"banner" bson:"banner"`
	TagLine        string                        `json:"tagLine" bson:"tagLine"`
	CreatedDate    int64                         `json:"created_date" bson:"created_date"`
	LastUpdated    int64                         `json:"last_updated" bson:"last_updated"`
	FollowCount    int64                         `json:"followCount" bson:"followCount"`
	FollowerCount  int64                         `json:"followerCount" bson:"followerCount"`
	AccessUserList []string                      `json:"accessUserList" bson:"accessUserList"`
	Permission     constants.UserPermissionConst `json:"permission" bson:"permission"`
}
