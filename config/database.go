package config

import (
	"context"
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	WriteDSN string
	ReadDSN  []string
}

var (
	writeDB  *gorm.DB
	readDBs  []*gorm.DB
	dbIndex1 uint32
	dbMutex  sync.RWMutex
)

func initDB() {
	dsn := AppConfig.Database.Dsn

	writeDB, err := initDatabase(dsn, "write")
	if err != nil {
		utils.Fatal("Failed to initialize write database: %v", err)
	}

	global.Db = writeDB

	err = writeDB.AutoMigrate(
		&models.Exchangerate{},
		&models.User{},
		&models.Article{},
	)
	if err != nil {
		utils.Fatal("Failed to migrate database schema: %v", err)
	}

	utils.Info("数据库连接成功，Schema迁移完成")
}

func initDatabase(dsn, dbType string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(AppConfig.Database.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	utils.Info("数据库连接成功 (%s)", dbType)
	return db, nil
}

func GetWriteDB() *gorm.DB {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return writeDB
}

func GetReadDB() *gorm.DB {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	if len(readDBs) == 0 {
		return writeDB
	}

	index := atomic.AddUint32(&dbIndex1, 1) % uint32(len(readDBs))
	return readDBs[index]
}

func WithTransaction(fn func(tx *gorm.DB) error) error {
	return writeDB.Transaction(fn)
}

func WithContext(ctx context.Context) *gorm.DB {
	return writeDB.WithContext(ctx)
}

func HealthCheck() error {
	sqlDB, err := writeDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func Close() error {
	var errs []error

	if writeDB != nil {
		if sqlDB, err := writeDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	for _, db := range readDBs {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
