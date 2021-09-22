package handlers

import (
	"net/http"

	"github.com/GMcD/telar-web/micros/collectives/database"
	service "github.com/GMcD/telar-web/micros/profile/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/pkg/parser"
	utils "github.com/red-gold/telar-core/utils"
)

type CollectiveQueryModel struct {
	Search     string      `query:"search"`
	Page       int64       `query:"page"`
	NotInclude []uuid.UUID `query:"nin"`
}

// QueryUserProfileHandle handle queru on userProfile
func QueryCollectiveHandle(c *fiber.Ctx) error {

	// Create service
	collectiveService, serviceErr := service.NewCollectiveService(database.Db)
	if serviceErr != nil {
		log.Error("NewCollectiveService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/collectiveService", "Error happened while creating collectiveService!"))
	}

	query := new(CollectiveQueryModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[QueryCollectiveHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	collectiveList, err := collectiveService.QueryCollective(query.Search, "created_date", query.Page, query.NotInclude)
	if err != nil {
		log.Error("[QueryCollective] %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/collectiveService", "Error happened while creating collectiveService!"))
	}

	return c.JSON(collectiveList)
}
