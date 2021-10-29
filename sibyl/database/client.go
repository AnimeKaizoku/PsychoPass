package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/AnimeKaizoku/sibylapi-go/sibyl/core/sibylConfig"
	sv "github.com/AnimeKaizoku/sibylapi-go/sibyl/core/sibylValues"
	"github.com/AnimeKaizoku/sibylapi-go/sibyl/core/utils/logging"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var SESSION *gorm.DB

func StartDatabase() {
	// check if `SESSION` variable is already established or not.
	// if yes, check if we have got any error from it or not.
	// if there is an error in the session, it mean we have to establish
	// a new connection again.
	if SESSION != nil && SESSION.Error == nil {
		return
	}

	var db *gorm.DB
	var err error
	if sibylConfig.SibylConfig.UseSqlite {
		db, err = gorm.Open(sqlite.Open(
			fmt.Sprintf("%s.db", sibylConfig.SibylConfig.DbName)),
			&gorm.Config{})
	} else {
		db, err = gorm.Open(postgres.Open(sibylConfig.SibylConfig.DbUrl), &gorm.Config{})
	}

	if err != nil {
		logging.Fatal("failed to connect to the database:", err)
	}

	SESSION = db
	logging.Info("Database connected")

	// Create tables if they don't exist
	err = SESSION.AutoMigrate(&sv.User{}, &sv.Token{})
	if err != nil {
		logging.Fatal(err)
	}

	if sibylConfig.SibylConfig.UseSqlite {
		dbMutex = &sync.Mutex{}
	}
	tokenMapMutex = &sync.Mutex{}
	tokenDbMap = make(map[int64]*sv.Token)
	userMapMutex = &sync.Mutex{}
	userDbMap = make(map[int64]*sv.User)
	go cleanMaps()
	logging.Info("Auto-migrated database schema")
}

func cleanMaps() {
	mtime := sibylConfig.GetMaxCacheTime()
	for {
		time.Sleep(mtime)
		if tokenDbMap == nil || userDbMap == nil {
			return
		}
		tokenMapMutex.Lock()
		for key, value := range tokenDbMap {
			if value == nil || time.Since(value.GetCacheDate()) > mtime {
				delete(tokenDbMap, key)
			}
		}
		tokenMapMutex.Unlock()
		userMapMutex.Lock()
		for key, value := range userDbMap {
			if value == nil || time.Since(value.GetCacheDate()) > mtime {
				delete(tokenDbMap, key)
			}
		}
		userMapMutex.Unlock()
	}
}

func IsFirstTime() bool {
	return SESSION.Find(&sv.Token{}).RowsAffected == 0
}
