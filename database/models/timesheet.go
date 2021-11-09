package models

import (
	"this/database"
	"time"
)

type Timesheet struct {
	ID          uint       `gorm:"->;column:id;primaryKey"`
	UserID      uint       `gorm:"->;column:userid"`
	Username    string     `gorm:"->;column:username"`
	SigninTime  *time.Time `gorm:"column:signintime1"`
	SignoutTime *time.Time `gorm:"column:signouttime1"`
	Total       float64    `gorm:"column:total"`
}

func (t *Timesheet) GetListByOffsetLimit(store database.Store, offset int, limit int, orderby string, order string) (timesheets []*Timesheet, err error) {
	if orderby == "Username" {
		orderby = "username"
	} else if orderby == "SigninTime" {
		orderby = "signintime1"
	} else if orderby == "SignoutTime" {
		orderby = "signouttime1"
	}
	result := store.DB.
		Offset(offset).
		Limit(limit).
		Order(orderby + " " + order).
		Table("timesheet").
		Select("timesheet.id, timesheet.signintime1, timesheet.signouttime1, timesheet.total, users.username as username").
		Joins("JOIN users ON users.userid = timesheet.userid").
		Group("timesheet.userid").
		Find(&timesheets)
	err = result.Error
	return
}

func (t *Timesheet) GetListForExcel(store database.Store) (timesheets []*Timesheet, err error) {
	result := store.DB.
		Table("timesheet").
		Order("timesheet.signintime1 desc").
		Select("timesheet.signintime1, timesheet.signouttime1, timesheet.total, users.username as username, timesheet.userid").
		Joins("JOIN users ON users.userid = timesheet.userid").
		Find(&timesheets)
	err = result.Error
	return
}

func (t *Timesheet) GetCount(store database.Store) (count int64, err error) {
	result := store.DB.Table("timesheet").Count(&count)
	err = result.Error
	return
}

func (t *Timesheet) UpdateById(store database.Store, id uint, signintime time.Time, signouttime time.Time, total float64) (count int64, err error) {
	result := store.DB.Table("timesheet").
		Select("signintime1", "signouttime1", "total").
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"signintime1":  signintime,
			"signouttime1": signouttime,
			"total":        total,
		})
	err = result.Error
	return
}
