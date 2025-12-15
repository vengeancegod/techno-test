package db

import (
	"context"
	"fmt"
	"log"
	"techno/internal/config/db"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(cfg db.DBConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host(),
		cfg.Port(),
		cfg.User(),
		cfg.Password(),
		cfg.DBName(),
		cfg.SSLMode(),
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Printf("failed to connect to db: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Printf("successfull connnect to DB: %s", cfg.DBName())
	return pool, nil
}

func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
