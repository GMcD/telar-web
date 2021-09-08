package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/middleware/authcookie"
	"github.com/red-gold/telar-core/middleware/authhmac"
	"github.com/red-gold/telar-web/micros/auth/handlers"

	"github.com/GMcD/cognito-jwt/verify"
)

type User struct {
	fullname string `json:"fullname" xml:"fullname" form:"fullname"`
	email string `json:"email" xml:"email" form:"email"`
	password string `json:"password" xml:"password" form:"password"`
}

func Signup2Handle(c *fiber.Ctx) error {
	config := coreConfig.AppConfig
	authConfig := &ac.AuthConfig

	
// take parameters from a URL 
	p := new(User)

	if err := c.BodyParser(p); err != nil {
	  return error
	}
	log.Println(p.fullname)
	log.Println(p.email)
	log.Println(p.password)
  // check parameters from URL

  // pass params into usermodel
  model := &models.SignupTokenModel{
	User: models.UserSignupTokenModel{
		Fullname: p.fullname,
		Email:    p.email,
		Password: p.password,
		},
	}
	// check model isn't empty
	if model.User.Fullname == "" {
		log.Error("Signup2Handle: missing fullname")
	}

	if model.User.Email == "" {
		log.Error("Signup2Handle: missing email")
	}

	if model.User.Password == "" {
		log.Error("Signup2Handle: missing password")
	}

	// Create service
	userAuthService, serviceErr := service.NewUserAuthService(database.Db)
	if serviceErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/recaptcha", serviceErr.Error()))
	}

	// Check user exist
	userAuth, findError := userAuthService.FindByUsername(model.User.Email)
	if findError != nil {
		errorMessage := fmt.Sprintf("Error while finding user by user name : %s",
			findError.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findUserAuth", errorMessage))

	}

	if userAuth != nil {
		err := utils.Error("userAlreadyExist", "User already exist - "+model.User.Email)
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	
	// Save into DB

	// if that worked we need to make the socialprofile

	socialName := generateSocialName(fullName, userId)
	newUserProfile := &models.UserProfileModel{
		ObjectId:    userUUID,
		FullName:    fullName,
		SocialName:  socialName,
		CreatedDate: createdDate,
		LastUpdated: createdDate,
		Email:       email,
		Avatar:      "https://util.telar.dev/api/avatars/" + userUUID.String(),
		Banner:      fmt.Sprintf("https://picsum.photos/id/%d/900/300/?blur", generateRandomNumber(1, 1000)),
		Permission:  constants.Public,
	}

	userProfileErr := saveUserProfile(newUserProfile)
	if userProfileErr != nil {
		log.Error("Save user profile %s", userProfileErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("canNotSaveUserProfile", "Cannot save user profile!"))
	}

	setupErr := initUserSetup(newUserAuth.ObjectId, newUserAuth.Username, "", newUserProfile.FullName, newUserAuth.Role)
	if setupErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("initUserSetupError", fmt.Sprintf("Cannot initialize user setup! error: %s", setupErr.Error())))
	}

	return c.SendStatus(http.StatusOK)

	tokenModel := &TokenModel{
		token:            ProviderAccessToken{},
		oauthProvider:    nil,
		providerName:     *coreConfig.AppConfig.AppName,
		profile:          &provider.Profile{Name: fullName, ID: userId, Login: email},
		organizationList: *coreConfig.AppConfig.OrgName,
		claim: UserClaim{
			DisplayName: fullName,
			SocialName:  socialName,
			Email:       email,
			UserId:      userId,
			Role:        "user",
			Banner:      newUserProfile.Banner,
			TagLine:     newUserProfile.TagLine,
			CreatedDate: newUserProfile.CreatedDate,
		},
	}
	session, sessionErr := createToken(tokenModel)
	if sessionErr != nil {
		errorMessage := fmt.Sprintf("Error creating session error: %s",
			sessionErr.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("initUserSetupError", errorMessage))

	}

	log.Info("\nSession is created: %s \n", session)
	webURL := authConfig.AuthConfig.ExternalRedirectDomain
	return c.Render("redirect", fiber.Map{
		"URL": webURL,
	})
	// change vewrify signup model

}