package migration

import (
	"beep/internal/config"
	"beep/internal/types"
	"log/slog"

	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB, config *config.Config) error {
	slog.Info("migrating database...")
	if !config.DB.AutoMigrate {
		slog.Info("skipped database migration")
		return nil
	}
	if err := db.AutoMigrate(types.User{},
		types.ModelFactory{},
		types.Model{},
		types.Workspace{},
		types.UserWorkspace{},
		types.KnowledgeBase{},
		types.Document{},
		types.MCPServer{},
		types.Agent{}); err != nil {
		return err
	}
	slog.Info("database migrated")
	return nil
}
