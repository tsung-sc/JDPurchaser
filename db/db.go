package db

import (
	"JD_Purchase/config"
	"JD_Purchase/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"sync"
	"time"
)

var (
	dbSession *gorm.DB
	once      sync.Once
	lock      = new(sync.Mutex)
)

func Instance() *gorm.DB {
	lock.Lock()
	defer lock.Unlock()
	once.Do(func() {
		// gorm init
		orm, err := gorm.Open(config.Get().DB.Type, config.Get().DB.URL)
		if err != nil {
			panic(err)
		}
		if config.Get().DB.MaxLifetime != 0 {
			orm.DB().SetConnMaxLifetime(time.Second * time.Duration(config.Get().DB.MaxLifetime))
		}
		if config.Get().DB.MaxIdleConn != 0 {
			orm.DB().SetMaxIdleConns(config.Get().DB.MaxIdleConn)
		}
		if config.Get().DB.MaxOpenConn != 0 {
			orm.DB().SetMaxOpenConns(config.Get().DB.MaxOpenConn)
		}

		orm.LogMode(config.Get().DB.Debug)
		orm.AutoMigrate(new(models.SysConfig), new(models.Goods))
		dbSession = orm
	})
	return dbSession
}
