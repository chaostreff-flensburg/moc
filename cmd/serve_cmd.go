package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"

	"github.com/sas1024/gorm-loggable"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/chaostreff-flensburg/moc/api"
	"github.com/chaostreff-flensburg/moc/config"
)

var serveCmd = cobra.Command{
	Use:   "serve",
	Short: "Serve API",
	Long:  "Serve the api",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

// serve server
func serve(config *config.Config) {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	if executeMigrate {
		log.Info("Migrate...")
		migrate(config)
	}

	if executeSeed {
		log.Info("Seed...")
		seed(config)
	}

	// ======================================
	// Database
	// ======================================
	log.Info("Init Database...")

	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open(config.Database.Driver, config.Database.Path)
		if err != nil {
			log.Error("try to connect...")
		} else {
			break
		}
		time.Sleep(5 * time.Second)
	}
	defer db.Close()
	log.Info("Database connected!")

	// ======================================
	// Log
	// ======================================
	_, err = loggable.Register(db)
	if err != nil {
		log.Error(err)
	}

	// ======================================
	// Server
	// ======================================
	server := api.NewAPI(db, config)

	server.ListenAndServe(":80")
}
