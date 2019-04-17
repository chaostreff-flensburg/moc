package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/chaostreff-flensburg/moc/config"
	"github.com/chaostreff-flensburg/moc/models"
)

var migrateCmd = cobra.Command{
	Use:   "migrate",
	Short: "Migrate database. Don't start the server",
	Long:  "Migrate database strucutures. This will create new tables and add missing collumns and indexes.",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, migrate)
	},
}

// migrate database
func migrate(config *config.Config) {
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
	// Migrate
	// ======================================
	log.Info("Migrate...")
	db.AutoMigrate(&models.Message{})
	log.Info("Finish...")
}
