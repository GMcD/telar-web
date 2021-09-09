package handlers

import (
	"fmt"
	"net/http"

	"github.com/GMcD/telar-web/constants"
	ac "github.com/GMcD/telar-web/micros/auth/config"
	"github.com/GMcD/telar-web/micros/auth/database"
	"github.com/GMcD/telar-web/micros/auth/dto"
	"github.com/GMcD/telar-web/micros/auth/models"
	"github.com/GMcD/telar-web/micros/auth/provider"
	service "github.com/GMcD/telar-web/micros/auth/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	coreConfig "github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/utils"
)

type User struct {
	fullname string `json:"fullname" xml:"fullname" form:"fullname" query:"fullname"`
	email    string `json:"email" xml:"email" form:"email" query:"email"`
	password string `json:"password" xml:"password" form:"password" query:"password"`
}

func Signup2Handle(c *fiber.Ctx) error {
	config := coreConfig.AppConfig
	authConfig := &ac.AuthConfig

	// take parameters from Request
	p := new(User)

	if err := c.BodyParser(p); err != nil {
		errorMessage := fmt.Sprintf("Unmarshal User Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("internal/userMarshal", "Can not parse body"))
	}
	log.Info(fmt.Sprintf("%+v\n", p))
	log.Info(p.fullname)
	log.Info(p.email)
	log.Info(p.password)
	log.Info(c.FormValue("fullname"))

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
		return c.Status(http.StatusBadRequest).JSON(utils.Error("missingFullname", "Missing fullname"))
	}

	if model.User.Email == "" {
		log.Error("Signup2Handle: missing email")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("missingEmail", "Missing email"))
	}

	if model.User.Password == "" {
		log.Error("Signup2Handle: missing password")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("missingPassword", "Missing password"))
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

	// Create New User Auth
	createdDate := utils.UTCNowUnix()
	userUUID := uuid.Must(uuid.NewV4())
	hashPassword, hashErr := utils.Hash(p.password)
	remoteIpAddress := c.IP()
	if hashErr != nil {
		errorMessage := fmt.Sprintf("Cannot hash the password! error: %s", hashErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("hashPassword", "Cannot hash the password!"))
	}

	// Create signup token
	newUserId := uuid.Must(uuid.NewV4())
	userVerificationService, serviceErr := service.NewUserVerificationService(database.Db)
	if serviceErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userVerificationService", serviceErr.Error()))
	}

	token := ""
	var tokenErr error
	token, tokenErr = userVerificationService.CreateEmailVerficationToken(service.EmailVerificationToken{
		UserId:          newUserId,
		HtmlTmplPath:    "views/email_code_verify.html",
		Username:        model.User.Email,
		EmailTo:         model.User.Email,
		EmailSubject:    "Your verification code",
		RemoteIpAddress: remoteIpAddress,
		FullName:        model.User.Fullname,
		UserPassword:    model.User.Password,
	}, &config)
	if tokenErr != nil {
		log.Error("Error on creating token: %s", tokenErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findUserAuth", "Error happened in creating token!"))
	}

	newUserAuth := &dto.UserAuth{
		ObjectId:      userUUID,
		Username:      p.email,
		Password:      hashPassword,
		AccessToken:   token,
		EmailVerified: true,
		Role:          "user",
		PhoneVerified: false,
		CreatedDate:   createdDate,
		LastUpdated:   createdDate,
	}
	userAuthErr := userAuthService.SaveUserAuth(newUserAuth)
	if userAuthErr != nil {

		errorMessage := fmt.Sprintf("Cannot save user authentication! error: %s", userAuthErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal", "Error happened during verification!"))
	}

	// Save into DB

	// if that worked we need to make the socialprofile

	socialName := generateSocialName(p.fullname, newUserId.String())
	newUserProfile := &models.UserProfileModel{
		ObjectId:    newUserId,
		FullName:    p.fullname,
		SocialName:  socialName,
		CreatedDate: createdDate,
		LastUpdated: createdDate,
		Email:       p.email,
		Avatar:      "https://util.telar.dev/api/avatars/" + newUserId.String(),
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

	// return c.SendStatus(http.StatusOK)

	tokenModel := &TokenModel{
		token:            ProviderAccessToken{},
		oauthProvider:    nil,
		providerName:     *config.AppName,
		profile:          &provider.Profile{Name: p.fullname, ID: newUserId.String(), Login: p.email},
		organizationList: *config.OrgName,
		claim: UserClaim{
			DisplayName: p.fullname,
			SocialName:  socialName,
			Email:       p.email,
			UserId:      newUserId.String(),
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
	webURL := authConfig.ExternalRedirectDomain
	return c.Render("redirect", fiber.Map{
		"URL": webURL,
	})
	// change verify signup model

}
