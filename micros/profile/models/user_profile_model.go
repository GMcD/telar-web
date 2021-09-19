package models

import (
	"github.com/GMcD/telar-web/constants"
	uuid "github.com/gofrs/uuid"
)

type UserProfileModel struct {
	ObjectId       uuid.UUID                     `json:"objectId" bson:"objectId"`
	FullName       string                        `json:"fullName" bson:"fullName", validate:"max=50"`
	SocialName     string                        `json:"socialName" bson:"socialName", validate:"max=100"`
	Avatar         string                        `json:"avatar" bson:"avatar", validate:"max=5000"`
	Banner         string                        `json:"banner" bson:"banner"`
	TagLine        string                        `json:"tagLine" bson:"tagLine", validate:"max=100"`
	CreatedDate    int64                         `json:"created_date" bson:"created_date"`
	LastUpdated    int64                         `json:"last_updated" bson:"last_updated"`
	LastSeen       int64                         `json:"lastSeen" bson:"lastSeen"`
	Email          string                        `json:"email" bson:"email", validate:"max=200"`
	Birthday       int64                         `json:"birthday" bson:"birthday"`
	WebUrl         string                        `json:"webUrl" bson:"webUrl", validate:"max=150"`
	CompanyName    string                        `json:"companyName" bson:"companyName", validate:"max=100"`
	Country        string                        `json:"country" bson:"country", validate:"max=100"`
	Address        string                        `json:"address" bson:"address"`
	Phone          string                        `json:"phone" bson:"phone"`
	VoteCount      int64                         `json:"voteCount" bson:"voteCount"`
	ShareCount     int64                         `json:"shareCount" bson:"shareCount"`
	FollowCount    int64                         `json:"followCount" bson:"followCount"`
	FollowerCount  int64                         `json:"followerCount" bson:"followerCount"`
	PostCount      int64                         `json:"postCount" bson:"postCount"`
	FacebookId     string                        `json:"facebookId" bson:"facebookId"`
	InstagramId    string                        `json:"instagramId" bson:"instagramId"`
	TwitterId      string                        `json:"twitterId" bson:"twitterId"`
	LinkedInId     string                        `json:"linkedInId"`
	AccessUserList []string                      `json:"accessUserList" bson:"accessUserList"`
	Permission     constants.UserPermissionConst `json:"permission" bson:"permission"`
}
