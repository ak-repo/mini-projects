package app

import (
	"log"

	"github.com/ak-repo/code-board/pkg/config"
	"github.com/ak-repo/code-board/pkg/db"
	"github.com/ak-repo/code-board/services/user-auth-service/internals/models"
	"gorm.io/gorm"
)

type AppContext struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewApp() *AppContext {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	database := db.Connect(cfg)
	db.Migrate(database, &models.User{}, &models.Address{})

	app := AppContext{
		DB:     database,
		Config: cfg,
	}
	return &app

}
