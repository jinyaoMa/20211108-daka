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

type User2Add struct {
	Username  string  `gorm:"column:username"`
	Password  string  `gorm:"column:password"`
	Usertype  uint    `gorm:"column:usertype"`
	Firstname string  `gorm:"column:firstname"`
	Lastname  string  `gorm:"column:lastname"`
	StoreID   uint    `gorm:"column:storeid"`
	City      string  `gorm:"column:city"`
	Province  string  `gorm:"column:province"`
	Barcode   string  `gorm:"column:barcode"`
	Rate      float64 `gorm:"column:commission_rate"`
	Deleted   int     `gorm:"column:deleted"`
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

func (t *User) AddUser(store database.Store, firstname string, lastname string, password string, usertype uint) (err error) {
	result := store.DB.Table("users").Create(&User2Add{
		Username:  firstname + " " + lastname[0:1],
		Password:  password,
		Usertype:  usertype,
		Firstname: firstname,
		Lastname:  lastname,
		StoreID:   0,
		City:      "1",
		Province:  "1",
		Barcode:   "0",
		Rate:      0,
		Deleted:   0,
	})
	err = result.Error
	return
}
