package cache

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"

	"github.com/parthoshuvo/authsvc/cfg"
	log "github.com/parthoshuvo/authsvc/log4u"
)

// TokenDB database.
type TokenDB struct {
	rdb *redis.Client
	ctx context.Context
}

// NewTokenDB creates a DB handler.
func NewTokenDB(dbDef *cfg.TokenDBDef) *TokenDB {
	ctx := context.Background()
	return &TokenDB{openDatabase(dbDef, ctx), ctx}
}

func openDatabase(dbDef *cfg.TokenDBDef, ctx context.Context) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", dbDef.Host, dbDef.Port),
		Password: dbDef.Password,
		DB:       dbDef.Database,
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to open database %s:%d/%s: [%v]", dbDef.Host, dbDef.Port, dbDef.Database, err)
	}
	log.Infof("%s!!, successfully cache database is connected!!", pong)
	return rdb
}

// Close closes database connection.
func (td *TokenDB) Close() {
	td.rdb.Close()
}
