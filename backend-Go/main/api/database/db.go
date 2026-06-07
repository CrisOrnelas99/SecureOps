package database

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"secureops/backend-go/api/config"
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

func EnsureSchema(ctx context.Context, database *gorm.DB) error {
	return database.WithContext(ctx).Exec(`
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS vulnerabilities (
	id BIGSERIAL PRIMARY KEY,
	cve_id TEXT NOT NULL,
	title TEXT NOT NULL,
	severity TEXT NOT NULL,
	description TEXT NOT NULL,
	status TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS assets (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	ip_address TEXT NOT NULL,
	operating_system TEXT,
	owner TEXT NOT NULL,
	criticality TEXT NOT NULL,
	risk_score SMALLINT NOT NULL DEFAULT 0,
	risk_level TEXT NOT NULL DEFAULT 'Low',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS asset_vulnerabilities (
	asset_id BIGINT NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
	vulnerability_id BIGINT NOT NULL REFERENCES vulnerabilities(id) ON DELETE CASCADE,
	PRIMARY KEY (asset_id, vulnerability_id)
);

CREATE TABLE IF NOT EXISTS waf_events (
	id BIGSERIAL PRIMARY KEY,
	method TEXT NOT NULL,
	path TEXT NOT NULL,
	reason TEXT NOT NULL,
	source_ip TEXT,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`).Error
}
