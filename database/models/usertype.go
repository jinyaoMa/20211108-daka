package models

import (
	"this/database"
)

type Usertype struct {
	ID   uint   `gorm:"->;column:usertype;primaryKey"`
	Name string `gorm:"->;column:typename"`
}

func (t *Usertype) GetById(store database.Store, id uint) (usertype *Usertype, err error) {
	result := store.DB.Table("usertype").Where("usertype = ?", id).Scan(&usertype)
	err = result.Error
	return
}
