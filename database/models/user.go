package models

import (
	"this/database"
)

type User struct {
	ID       uint   `gorm:"->;column:userid;primaryKey"`
	Account  string `gorm:"->;column:username"`
	Password string `gorm:"->;column:password"`
	TypeID   uint   `gorm:"->;column:usertype"`
}

func (t *User) GetAll(store database.Store) (users []*User, err error) {
	result := store.DB.Find(&users)
	err = result.Error
	return
}

func (t *User) GetUserByAccountPassword(store database.Store, account string, password string) (user *User, err error) {
	result := store.DB.
		Where("username = ? AND password = ?", account, password).
		Find(&user)
	err = result.Error
	return
}
