package db

import (

	// import source file

	"api/internal/pkg/domain/domain_model/entity"
	"errors"
	"fmt"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	DB *gorm.DB
}

func NewDB() (Database, error) {
	//dsn := "bac4178dc89368:292965a5@tcp(us-cdbr-east-05.cleardb.net)/heroku_560fb6556eff9f8?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "hiendv:password@tcp(mysql:3306)/usersdb?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return Database{
		DB: db,
	}, err
}
func (db *Database) MigrateDBWithGorm() {
	db.DB.Migrator().AutoMigrate(
		&entity.Users{},
		&entity.UserBanned{},
		&entity.Room{},
	)

}
func (db *Database) Begin() Database {
	return Database{
		DB: db.DB.Begin(),
	}
}
func (db *Database) RollBack() {
	db.DB.Rollback()
}
func (db *Database) Commit() error {
	return db.DB.Commit().Error
}
func (db *Database) First(value interface{}, condition interface{}) error {
	err := db.DB.Preload(clause.Associations).First(value, condition).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
func (db *Database) Find(value interface{}, condition interface{}) error {
	err := db.DB.Preload(clause.Associations).Find(value, condition).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (db *Database) FindWithPagination(value interface{}, offset int, pageSize int, condition interface{}) error {
	err := db.DB.Offset(offset).Limit(pageSize).Preload(clause.Associations).Find(value, condition).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
func (db *Database) Create(value interface{}) error {
	//	fmt.Println(value)
	err := db.DB.Create(value).Error
	return err
}
func (db *Database) CreateMany(value interface{}, model interface{}) error {
	fmt.Println(model)
	err := db.DB.Model(model).Create(value).Error
	return err
}
func (db *Database) Delete(value interface{}) error {
	return db.DB.Delete(value).Error
}
func (db *Database) Update(model interface{}, oldVal interface{}, newVal interface{}) error {
	return db.DB.Model(model).Where(oldVal).Updates(newVal).Error
}
func (db *Database) ExcQuery(out interface{}, query string, values ...interface{}) error {
	return db.DB.Exec(query, values...).Scan(out).Error
}
func (db *Database) RawQuery(out interface{}, query string, values ...interface{}) error {
	err := db.DB.Raw(query, values...).Scan(out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
