package handlers

import (
	"fmt"
	"net/http"

	"github.com/GMcD/cognito-jwt/verify"
	"github.com/GMcD/telar-web/constants"
	"github.com/GMcD/telar-web/micros/auth/database"
	"github.com/GMcD/telar-web/micros/auth/dto"
	"github.com/GMcD/telar-web/micros/auth/models"
	"github.com/GMcD/telar-web/micros/auth/provider"
	service "github.com/GMcD/telar-web/micros/auth/services"
	"github.com/gofiber/fiber/v2"
	coreConfig "github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/utils"
	"github.com/sethvargo/go-password/password"
)

type User struct {
	userid   string `json:"fullname" xml:"fullname" form:"fullname" query:"fullname"`
	fullname string `json:"fullname" xml:"fullname" form:"fullname" query:"fullname"`
	email    string `json:"email" xml:"email" form:"email" query:"email"`
	password string `json:"password" xml:"password" form:"password" query:"password"`
}

func Signup2Handle(c *fiber.Ctx) error {
	config := coreConfig.AppConfig

	// take Cognito JWT token from Authorization:
	jwt = c.get("Authorization")
	claims = verify.VerifyJWT(jwt)

	// take parameters from Request
	p := new(User)
	p.userid = claims["cognito:username"]
	p.fullname = claims["name"]
	p.email = claims["email"]
	p.password, _ = password.Generate(12, 4, 4, false, false)

	log.Info(fmt.Sprintf("%+v\n", p))

	if p.userid == "" {
		log.Error("Signup2Handle: missing JWT claim : cognito:username")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("missingFullname", "Missing fullname"))
	}

	if p.fullname == "" {
		log.Error("Signup2Handle: missing form value fullname")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("missingFullname", "Missing fullname"))
	}

	if p.email == "" {
		log.Error("Signup2Handle: missing form value email")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("missingEmail", "Missing email"))
	}

	if p.password == "" {
		log.Error("Signup2Handle: missing form value password")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("missingPassword", "Missing password"))
	}

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
		log.Error("Signup2Handle: missing model fullname")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("missingFullname", "Missing fullname"))
	}

	if model.User.Email == "" {
		log.Error("Signup2Handle: missing model email")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("missingEmail", "Missing email"))
	}

	if model.User.Password == "" {
		log.Error("Signup2Handle: missing model password")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("missingPassword", "Missing password"))
	}

	// Create service
	userAuthService, serviceErr := service.NewUserAuthService(database.Db)
	if serviceErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/newUserAuthService", serviceErr.Error()))
	}

	// Check user exist
	userAuth, findError := userAuthService.FindByUsername(model.User.Email)
	if findError != nil {
		errorMessage := fmt.Sprintf("Error while finding user by user name : %s",
			findError.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userAuthService", errorMessage))

	}

	if userAuth != nil {
		err := utils.Error("userAlreadyExists", "User already exists - "+model.User.Email)
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	// Create New User Auth
	createdDate := utils.UTCNowUnix()
	userUUID := p.userid
	hashPassword, hashErr := utils.Hash(p.password)
	//	remoteIpAddress := c.IP()
	if hashErr != nil {
		errorMessage := fmt.Sprintf("Cannot hash the password! error: %s", hashErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("hashPassword", "Cannot hash the password!"))
	}

	// Create UserAuth pre-verified with default token
	token := "000000"
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

	// Create User Profile
	socialName := generateSocialName(p.fullname, userUUID.String())
	newUserProfile := &models.UserProfileModel{
		ObjectId:    userUUID,
		FullName:    p.fullname,
		SocialName:  socialName,
		CreatedDate: createdDate,
		LastUpdated: createdDate,
		Email:       p.email,
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

	// return c.SendStatus(http.StatusOK)

	tokenModel := &TokenModel{
		token:            ProviderAccessToken{},
		oauthProvider:    nil,
		providerName:     *config.AppName,
		profile:          &provider.Profile{Name: p.fullname, ID: userUUID.String(), Login: p.email},
		organizationList: *config.OrgName,
		claim: UserClaim{
			DisplayName: p.fullname,
			SocialName:  socialName,
			Email:       p.email,
			UserId:      userUUID.String(),
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

	return c.JSON(fiber.Map{
		"token":   session,
		"version": utils.PkgVersion("telar-web"),
	})
}
