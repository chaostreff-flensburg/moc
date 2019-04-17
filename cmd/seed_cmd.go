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

var seedCmd = cobra.Command{
	Use:   "seed",
	Short: "Seed database. Don't start the server",
	Long:  "Seed database data. This will create random seed data to test the setup.",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, seed)
	},
}

// seed database with testdata
func seed(config *config.Config) {
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
	// Add Test Data
	// ======================================
	models.Seed(db)
}
