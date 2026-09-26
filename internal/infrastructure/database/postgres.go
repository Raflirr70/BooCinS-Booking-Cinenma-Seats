package database

import (
	"fmt"
	"log"

	"github.com/rafli/boocins/config"
	"github.com/rafli/boocins/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	db.Exec("ALTER TABLE seats DROP CONSTRAINT IF EXISTS seats_room_id_label_key")
	db.Exec("ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_name_key")
	db.Exec("DROP INDEX IF EXISTS idx_rooms_name")
	db.Exec("DROP INDEX IF EXISTS idx_rooms_name_active")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_rooms_name_active ON rooms(name) WHERE deleted_at IS NULL")

	seedDefaults(db)

	return db, nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Role{},
		&entity.Permission{},
		&entity.RolePermission{},
		&entity.Membership{},
		&entity.User{},
		&entity.Film{},
		&entity.Genre{},
		&entity.Media{},
		&entity.Room{},
		&entity.Seat{},
		&entity.Schedule{},
		&entity.ScheduleSeat{},
		&entity.GuestOrder{},
		&entity.EmailVerification{},
		&entity.Transaction{},
		&entity.Ticket{},
		&entity.Bookmark{},
		&entity.Rating{},
		&entity.SearchHistory{},
		&entity.Promo{},
		&entity.AuditLog{},
	)

}

func seedDefaults(db *gorm.DB) {
	roles := []entity.Role{
		{Name: "super_admin", Description: "Full system access"},
		{Name: "manager", Description: "Monitoring & reporting"},
		{Name: "admin", Description: "Content & schedule management"},
		{Name: "staff", Description: "Offline ticket sales"},
		{Name: "member", Description: "Registered user"},
	}
	for _, r := range roles {
		db.FirstOrCreate(&r, entity.Role{Name: r.Name})
	}

	memberships := []entity.Membership{
		{Name: "Bronze", MinScore: 0, Discount: 0.00},
		{Name: "Silver", MinScore: 100, Discount: 5.00},
		{Name: "Gold", MinScore: 500, Discount: 10.00},
		{Name: "Platinum", MinScore: 1000, Discount: 15.00},
	}
	for _, m := range memberships {
		db.FirstOrCreate(&m, entity.Membership{Name: m.Name})
	}

	log.Println("[SEED] Default roles and memberships seeded")
}
