package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SaveliiYam/http_service_with_db/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultDSN = "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"

type Repository struct {
	db *gorm.DB
}

func NewRepository() (*Repository, error) {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = defaultDSN
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (models.User, error) {
	var row userRow
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		return models.User{}, err
	}

	return row.toModel(), nil
}

func (r *Repository) Create(ctx context.Context, name string) (models.User, error) {
	row := userRow{Name: name}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return models.User{}, err
	}

	return row.toModel(), nil
}

// Schema is owned by goose migrations, so AutoMigrate is intentionally not used.
type userRow struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:now()"`
}

func (userRow) TableName() string {
	return "users"
}

func (r userRow) toModel() models.User {
	return models.User{ID: r.ID, Name: r.Name}
}
