package handlers

import (
	"fmt"
	"net/http"

	authConfig "github.com/GMcD/telar-web/micros/auth/config"
	"github.com/GMcD/telar-web/micros/auth/database"
	models "github.com/GMcD/telar-web/micros/auth/models"
	"github.com/GMcD/telar-web/micros/auth/provider"
	service "github.com/GMcD/telar-web/micros/auth/services"
	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
	utils "github.com/red-gold/telar-core/utils"
)

// Login2Handle creates a handler for logging in telar social
func Login2Handle(c *fiber.Ctx) error {

	model := &models.LoginModel{
		Username:     c.FormValue("username"),
		Password:     c.FormValue("password"),
		ResponseType: SPAResponseType,
	}

	// Create service
	userAuthService, serviceErr := service.NewUserAuthService(database.Db)
	if serviceErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userAuthService", serviceErr.Error()))
	}

	if model.Username == "" {
		log.Error("Username is required!")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("usernameIsRequired", "Username is required!"))

	}

	if model.Password == "" {
		log.Error("Password is required!")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("passwordIsRequired", "Password is required!"))
	}

	foundUser, err := userAuthService.FindByUsername(model.Username)
	if err != nil || foundUser == nil {
		if err != nil {
			log.Error("User not found %s", err.Error())
		}
		log.Error("User not found!")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("findUserByUserName", "User not found!"))

	}

	if !foundUser.EmailVerified && !foundUser.PhoneVerified {
		log.Error("User is not verified!")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("userNotVerified", "User is not verified!"))
	}

	log.Info(" foundUser.Password: %s, model.Password: %s", foundUser.Password, model.Password)
	compareErr := utils.CompareHash(foundUser.Password, []byte(model.Password))
	if compareErr != nil {
		log.Error("Password doesn't match %s", compareErr.Error())
		return c.Status(http.StatusUnauthorized).JSON(utils.Error("passwordNotMatch", "Password doesn't match!"))
	}

	profileChannel := readProfileAsync(foundUser.ObjectId)
	langChannel := readLanguageSettingAsync(foundUser.ObjectId,
		&UserInfoInReq{UserId: foundUser.ObjectId, Username: foundUser.Username, SystemRole: foundUser.Role})

	profileResult, langResult := <-profileChannel, <-langChannel
	if profileResult.Error != nil || profileResult.Profile == nil {
		if profileResult.Error != nil {
			log.Error(" User profile  %s", profileResult.Error.Error())
		}
		return c.Status(http.StatusBadRequest).JSON(utils.Error("internal/getUserProfile", "Can not find user profile!"))
	}

	currentUserLang := "en"
	fmt.Println("langResult.settings", langResult.settings)
	langSettigPath := getSettingPath(foundUser.ObjectId, "lang", "current")
	if val, ok := langResult.settings[langSettigPath]; ok && val != "" {
		currentUserLang = val
	} else {
		go func() {
			userInfoReq := &UserInfoInReq{
				UserId:      foundUser.ObjectId,
				Username:    foundUser.Username,
				Avatar:      profileResult.Profile.Avatar,
				DisplayName: profileResult.Profile.FullName,
				SystemRole:  foundUser.Role,
			}
			createDefaultLangSetting(userInfoReq)
		}()
	}

	tokenModel := &TokenModel{
		token:            ProviderAccessToken{},
		oauthProvider:    nil,
		providerName:     "telar",
		profile:          &provider.Profile{Name: foundUser.Username, ID: foundUser.ObjectId.String(), Login: foundUser.Username},
		organizationList: "Red Gold",
		claim: UserClaim{
			DisplayName: profileResult.Profile.FullName,
			SocialName:  profileResult.Profile.SocialName,
			Email:       profileResult.Profile.Email,
			Avatar:      profileResult.Profile.Avatar,
			Banner:      profileResult.Profile.Banner,
			TagLine:     profileResult.Profile.TagLine,
			UserId:      foundUser.ObjectId.String(),
			Role:        foundUser.Role,
			CreatedDate: foundUser.CreatedDate,
		},
	}
	session, err := createToken(tokenModel)
	if err != nil {
		log.Error("Error creating session: %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/createToken", "Internal server error creating token"))
	}

	// Write session on cookie
	writeSessionOnCookie(c, session, &authConfig.AuthConfig)

	// Write user language on cookie
	writeUserLangOnCookie(c, currentUserLang)

	return c.JSON(fiber.Map{
		"user":    profileResult.Profile,
		"lang":    currentUserLang,
		"session": session,
	})

}
