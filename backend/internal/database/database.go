package database

import (
	"backend/internal/domain/domain"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var (
	// host     = os.Getenv("DB_HOST")
	// port     = os.Getenv("DB_PORT")
	// username     = os.Getenv("DB_USER")
	// password = os.Getenv("DB_PASSWORD")
	// database   = os.Getenv("DB_NAME")
	// sslmode  = "disable"
	dbInstance *service
	// schema = os.Getenv(" 0[]")
)

func ConnectDB() *gorm.DB {
	// define db config
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	sslmode := "disable"
	fmt.Printf("DB port is at: %s from DB\n", os.Getenv("DB_PORT"))

	fmt.Printf("value of host is : %s \n", host)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, username, password, database, port, sslmode)

	fmt.Printf("Valus of dsn is %s: \n", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
		os.Exit(1)
	}

	fmt.Println("Db is :", db)
	log.Println("connect to DB successfull")
	DB = db
	return db

}

type Service interface {
	Health() map[string]string
	Close() error
}

type service struct {
	db *gorm.DB
}

// func New() Service {
// 	if dbInstance != nil {
// 		return dbInstance
// 	}
// 	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
// 		host, username, password, database, port, sslmode)

// 	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	dbInstance = &service{
// 		db: db,
// 	}
// 	return dbInstance
// }

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&domain.User{},
		&domain.UserSkill{},
		&domain.Company{},
		&domain.RefreshToken{},
		&domain.ApplicantProfile{},
		&domain.WorkExperience{},
		&domain.Education{},
		&domain.Institute{},
		&domain.FieldOfStudy{},
		&domain.PortfolioItem{},
		&domain.CompanyMember{},
		&domain.Skill{},
		&domain.Notification{},
		&domain.OAuthAccount{},
		&domain.Project{},
		&domain.Application{},
		&domain.Conversation{},
		&domain.ConversationParticipant{},
		&domain.Message{},
		&domain.ProjectTask{},
		&domain.ProjectSkill{},
		&domain.Submission{},
	)
}

func (s *service) Health() map[string]string {

	stats := make(map[string]string)

	db, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down %v", err)
		return stats
	}

	ctx, cancle := context.WithTimeout(context.Background(), 1*time.Second)
	// make cancle() working after this function finished before return or end process
	defer cancle()

	// ping to db
	err = db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down %v", err)
		log.Fatal("db down %v", err)
		return stats
	}

	dbStats := db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// func (s *service) Close() error {
// 	db, err := s.db.DB()
// 	if err != nil {
// 		return  err
// 	}
// 	log.Printf("Disconnected from database: %s", database)
// 	return db.Close()
// }
