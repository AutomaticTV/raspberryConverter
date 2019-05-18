package player

import (
	"errors"
	"fmt"
	"raspberryConverter/models"

	"gopkg.in/validator.v2"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //gorm requires sqlite
)

// dbPath represents the filepath where sqlite files are stored
const dbPath = "/var/lib/raspberryConverter/playerConfig.db"

// initStorage creates the DB table (if it's not done yet) and set a default entry
func initStorage() error {
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	db.AutoMigrate(&models.PlayerConfig{})
	return createDefaultConfig()
}

// GetConfig is a function that retrieves the player config from persistance
func GetConfig() (models.PlayerConfig, error) {
	// connect to db
	db, err := gorm.Open("sqlite3", dbPath)
	var c models.PlayerConfig
	defer db.Close()
	if err != nil {
		return c, err
	}
	// get user with given username
	if dbc := db.Last(&c); dbc.Error != nil {
		return c, dbc.Error
	}
	return c, nil
}

// SetConfig is a function that stores to persistance
func SetConfig(c models.PlayerConfig) error {
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	err = validator.Validate(c)
	if err != nil {
		fmt.Println(err)
		return errors.New("The player configuration is not valid")
	}
	db.Create(&c)
	return nil
}

// createDefaultConfig set the first record on the DB
func createDefaultConfig() error {
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	registeredValues := 0
	db.Find(&models.PlayerConfig{}).Count(&registeredValues)
	if registeredValues == 0 {
		db.Create(&models.PlayerConfig{
			Video:         "1080p59.94",
			AudioDecoding: "Hardware",
			URL:           "https://some.example",
			Buffer:        300,
			Username:      "",
			Password:      "",
			Volume:        0,
			Autoplay:      "Yes",
		})
	}
	return nil
}
