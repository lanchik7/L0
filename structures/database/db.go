package db

import (
	"encoding/json"
	"fmt"
	"log"

	ls "github.com/niciki/go-NatsService/structures/localStore"
	so "github.com/niciki/go-NatsService/structures/structOrder"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseRecord struct {
	OrderUid  string `gorm:"primaryKey;"`
	OrderJson string
}

type Database struct {
	db *gorm.DB
}

func InitDb(port string) (Database, error) {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=localhost port=%s user=aslan password=1302 dbname=lanchik sslmode=disable", port)))
	if err != nil {
		return Database{}, err
	}
	err = db.AutoMigrate(&DatabaseRecord{})
	if err != nil {
		return Database{}, err
	}
	return Database{db}, nil
}

func (d *Database) UploadCache(cache *ls.Store) error {
	var databaseData []DatabaseRecord
	err := d.db.Find(&databaseData)
	for _, rec := range databaseData {
		var JsonOrder so.Order
		err := json.Unmarshal([]byte(rec.OrderJson), &JsonOrder)
		if err != nil {
			log.Print(err)
		}
		cache.Add(JsonOrder)
	}
	return err.Error
}

func (d *Database) AddRecord(rec ...so.Order) error {
	log.Printf("Add to db record: %v", rec)
	for _, r := range rec {
		order, err := json.MarshalIndent(r, "", "\t")
		if err != nil {
			log.Print(err)
		}
		errGorm := d.db.Create(DatabaseRecord{r.OrderUid, string(order)})
		if errGorm.Error != nil {
			return errGorm.Error
		}
	}
	return nil
}
