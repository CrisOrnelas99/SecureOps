package utils

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"secureops/backend-go/api/config"
	"secureops/backend-go/api/model"
)

func Connect(ctx context.Context, cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func Close(database *gorm.DB) error {
	sqlDB, err := database.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func RunMigrations(ctx context.Context, database *gorm.DB) error {
	if err := ensureUserSchema(ctx, database); err != nil {
		return err
	}

	if err := database.WithContext(ctx).AutoMigrate(
		&model.Vulnerability{},
		&model.Asset{},
	); err != nil {
		return err
	}

	return ensureIndexes(ctx, database)
}

func ensureUserSchema(ctx context.Context, database *gorm.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			email VARCHAR NOT NULL,
			password_hash VARCHAR NOT NULL,
			role VARCHAR NOT NULL DEFAULT 'user'
		)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR NOT NULL DEFAULT 'user'`,
	}

	for _, statement := range statements {
		if err := database.WithContext(ctx).Exec(statement).Error; err != nil {
			return err
		}
	}

	return nil
}

func ensureIndexes(ctx context.Context, database *gorm.DB) error {
	statements := []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email)`,
		`ALTER TABLE users DROP CONSTRAINT IF EXISTS ukr43af9ap4edm43mmtq01oddj6`,
		`ALTER TABLE users DROP CONSTRAINT IF EXISTS uk6dotkott2kjsp8vw4d0m25fb7`,
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM vulnerabilities GROUP BY cve_id HAVING count(*) > 1
			) THEN
				CREATE UNIQUE INDEX IF NOT EXISTS idx_vulnerabilities_cve_id ON vulnerabilities (cve_id);
			END IF;
		END $$`,
		`CREATE INDEX IF NOT EXISTS idx_assets_user_id ON assets (user_id)`,
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_role'
			) THEN
				ALTER TABLE users ADD CONSTRAINT chk_users_role CHECK (role IN ('admin', 'user'));
			END IF;
		END $$`,
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'chk_vulnerabilities_severity'
			) THEN
				ALTER TABLE vulnerabilities ADD CONSTRAINT chk_vulnerabilities_severity CHECK (severity IN ('Low', 'Medium', 'High', 'Critical'));
			END IF;
		END $$`,
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'chk_vulnerabilities_status'
			) THEN
				ALTER TABLE vulnerabilities ADD CONSTRAINT chk_vulnerabilities_status CHECK (status IN ('Open', 'Fixed', 'In Progress'));
			END IF;
		END $$`,
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_assets_user'
			) THEN
				ALTER TABLE assets ADD CONSTRAINT fk_assets_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
			END IF;
		END $$`,
		`ALTER TABLE asset_vulnerabilities DROP CONSTRAINT IF EXISTS fkavovmmqdpqv6hacqhae27ngt1`,
		`ALTER TABLE asset_vulnerabilities DROP CONSTRAINT IF EXISTS fkpldrve7axqj2xnyb09ojqmd02`,
	}

	for _, statement := range statements {
		if err := database.WithContext(ctx).Exec(statement).Error; err != nil {
			return err
		}
	}

	return nil
}
