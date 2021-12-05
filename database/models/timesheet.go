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

type ExcelTotal struct {
	Name  string  `gorm:"column:fullname"`
	Hours float64 `gorm:"column:hours"`
}

type ExcelTimesheet struct {
	Name        string     `gorm:"column:fullname"`
	SigninTime  *time.Time `gorm:"column:signintime1"`
	SignoutTime *time.Time `gorm:"column:signouttime1"`
	Total       float64    `gorm:"column:total"`
}

func (t *Timesheet) GetListByOffsetLimit(store database.Store, offset int, limit int, orderby string, order string, date string) (timesheets []*Timesheet, err error) {
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
		Order(orderby+" "+order).
		Table("timesheet").
		Select("timesheet.id, timesheet.signintime1, timesheet.signouttime1, timesheet.total, users.username as username").
		Where("date(timesheet.signintime1) = ?", date).
		Joins("JOIN users ON users.userid = timesheet.userid").
		Find(&timesheets)
	err = result.Error
	return
}

func (t *Timesheet) GetListByOffsetLimitAll(store database.Store, offset int, limit int, orderby string, order string) (timesheets []*Timesheet, err error) {
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
		Find(&timesheets)
	err = result.Error
	return
}

func (t *Timesheet) GetTotalForExcelWithRange(store database.Store, startDate string, endDate string) (excelTotals []*ExcelTotal, err error) {
	result := store.DB.
		Table("timesheet").
		Order("timesheet.signintime1 desc").
		Select("CONCAT(users.firstname,' ',users.lastname) as fullname, SUM(timesheet.total) as hours").
		Where("date(timesheet.signintime1) between ? and ?", startDate, endDate).
		Joins("JOIN users ON users.userid = timesheet.userid").
		Group("users.userid").
		Scan(&excelTotals)
	err = result.Error
	return
}

func (t *Timesheet) GetListForExcelWithRange(store database.Store, startDate string, endDate string) (timesheets []*ExcelTimesheet, err error) {
	result := store.DB.
		Table("timesheet").
		Order("timesheet.signintime1 desc").
		Select("timesheet.signintime1, timesheet.signouttime1, timesheet.total, CONCAT(users.firstname,' ',users.lastname) as fullname").
		Where("date(timesheet.signintime1) between ? and ?", startDate, endDate).
		Joins("JOIN users ON users.userid = timesheet.userid").
		Scan(&timesheets)
	err = result.Error
	return
}

func (t *Timesheet) GetCountAll(store database.Store) (count int64, err error) {
	result := store.DB.Table("timesheet").Count(&count)
	err = result.Error
	return
}

func (t *Timesheet) GetCount(store database.Store, date string) (count int64, err error) {
	result := store.DB.Table("timesheet").Where("date(timesheet.signintime1) = ?", date).Count(&count)
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
