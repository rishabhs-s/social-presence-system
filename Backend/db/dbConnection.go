package db

import (
	"fmt"
	"log"
	"log/slog"

	// "os"

	"github.com/rishabhs-s/models"
	// "github.com/rishabhs-s/websockets"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectToDB()
	return Db
}

func connectToDB() *gorm.DB {
	var err error

	host := "postgres-db" //for docker
	// host := "localhost" //local dev
	username := "postgres"
	password := "root"
	dbname := "supergaming"
	port := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, username, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		slog.Error("Error connecting to Database: ", err)
		return nil
	}
	slog.Info("Successfully connected to Database")

	return db
}

func Migrate() {

	if err := Db.AutoMigrate(&models.OnlineUser{}); err != nil {
		panic("failed to migrate OnlineUser model")
	}
	if err := Db.AutoMigrate(&models.FriendRequest{}); err != nil {
		panic("failed to migrate FriendRequest model")
	}
	if err := Db.AutoMigrate(&models.PartyInvitation{}); err != nil {
		panic("failed to migrate PartyInvitation model")
	}

	if err := Db.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate User model")
	}

	// Auto migrate the Party model
	if err := Db.AutoMigrate(&models.Party{}); err != nil {
		panic("failed to migrate Party model")
	}

	// mockNames := websockets.GenerateMockData(10) // Generate 10 mock names
	// websockets.PopulateTable(mockNames,Db)
	log.Print("All tables created Successfully")
}
