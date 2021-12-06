package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Store struct {
	Name string
	Dsn  string
	DB   *gorm.DB
}

var MAIN *Store = &Store{
	Name: "Main DB",
	Dsn:  "root:root@tcp(127.0.0.1:7531)/loveshop0?charset=utf8mb4&parseTime=True&loc=Local",
}

var STORES = []*Store{
	{
		Name: "Store 0",
		Dsn:  "root:root@tcp(127.0.0.1:7531)/loveshop3?charset=utf8mb4&parseTime=True&loc=Local",
	},
	{
		Name: "Store 1",
		Dsn:  "root:root@tcp(127.0.0.1:7531)/loveshop1?charset=utf8mb4&parseTime=True&loc=Local",
	},
}

var WAREHOUSE_INDEX = 0     // e.g. index of warehouse for addUser
var OFFICE_EXCEPT_INDEX = 1 // e.g. office users cannot access to Store 1

func GetStore(index int) (store *Store, ok bool) {
	if index < 0 || index > len(STORES)-1 {
		ok = false
		return
	}
	store = STORES[index]
	ok = true
	return
}

func init() {
	var err error

	MAIN.DB, err = gorm.Open(mysql.Open(MAIN.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connect error %v\n", err)
	}
	if MAIN.DB.Error != nil {
		log.Fatalf("database error %v\n", MAIN.DB.Error)
	}

	for _, store := range STORES {
		store.DB, err = gorm.Open(mysql.Open(store.Dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("database connect error %v\n", err)
		}
		if store.DB.Error != nil {
			log.Fatalf("database error %v\n", store.DB.Error)
		}
	}
}
