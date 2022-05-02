package persistence

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetClient(host, user, password, dbname, port string) (*gorm.DB, error) {
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	dsn := fmt.Sprintf("host=localhost user=postgres password=admin dbname=NOVA port=5432 sslmode=disable")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
