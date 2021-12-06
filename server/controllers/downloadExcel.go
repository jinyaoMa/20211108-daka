package controllers

import (
	"net/http"
	"strconv"
	"this/database"
	"this/database/models"
	"this/gate"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"github.com/xuri/excelize/v2"
)

type DownloadExcelQuery struct {
	StoreID   int    `form:"storeId"`
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
	IsAll     bool   `form:"all"`
}

// @Summary DownloadExcel
// @Description Download Excel
// @Tags After Authorization
// @accept plain
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @param storeId query int false "Store ID"
// @param startDate query string true "Start Date"
// @param endDate query string true "End Date"
// @param all query bool false "Is All"
// @Success 200 ".xlsx file"
// @Failure 401 "Auth failed"
// @Failure 404 "{ error }"
// @Router /timesheet/download [get]
func DownloadExcel(c *gin.Context) {
	temp, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Middleware Auth not exists!",
		})
		return
	}
	claims := temp.(*gate.Claims)

	var form DownloadExcelQuery
	if c.ShouldBind(&form) != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Check for empty fields!",
		})
		return
	}

	if claims.Usertype == "Office" && form.StoreID == database.OFFICE_EXCEPT_INDEX {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usertype Office denied by store!",
		})
		return
	}

	store, ok := database.GetStore(form.StoreID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Store error!",
		})
		return
	}

	if form.IsAll {
		form.StartDate = "1000-01-01"
		form.EndDate = "9999-12-31"
	}

	var timesheet *models.Timesheet
	excelTotals, err := timesheet.GetTotalForExcelWithRange(*store, form.StartDate, form.EndDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Timesheet error!",
		})
		return
	}
	list, err := timesheet.GetListForExcelWithRange(*store, form.StartDate, form.EndDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Timesheet error!",
		})
		return
	}

	f := excelize.NewFile()
	makeExcelTotal(f, excelTotals)
	makeExcelNewSheet(f, list)

	filename := database.STORES[form.StoreID].Name + "_" + form.StartDate + "_" + form.EndDate + ".xlsx"
	if form.IsAll {
		filename = database.STORES[form.StoreID].Name + "_" + "all.xlsx"
	}

	if err := f.SaveAs(filename); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "file create error!",
		})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filename)
}

type Subtotal struct {
	Name string
	Sum  float64
}

func makeExcelTotal(f *excelize.File, excelTotals []*models.ExcelTotal) (err error) {
	var styleTitleId int
	styleTitleId, err = f.NewStyle(`{"font":{"bold":true},"fill":{"type":"pattern","color":["#FFFF00"],"pattern":1}}`)
	var styleGridId int
	styleGridId, err = f.NewStyle(`{
		"border": [{
			"type": "left",
			"color": "000000",
			"style": 4
		},
		{
			"type": "top",
			"color": "000000",
			"style": 4
		},
		{
			"type": "bottom",
			"color": "000000",
			"style": 4
		},
		{
			"type": "right",
			"color": "000000",
			"style": 4
		}]
	}`)

	f.SetSheetName("Sheet1", "Total")
	f.SetCellValue("Total", "A1", "Name")
	f.SetCellValue("Total", "B1", "Hours")
	f.SetCellStyle("Total", "A1", "B1", styleTitleId)

	f.SetColWidth("Total", "A", "A", 30)
	f.SetColWidth("Total", "B", "B", 10)

	var cursor int64 = 2
	for _, row := range excelTotals {
		f.SetCellValue("Total", "A"+strconv.FormatInt(cursor, 10), row.Name)
		f.SetCellValue("Total", "B"+strconv.FormatInt(cursor, 10), row.Hours)
		f.SetCellStyle("Total", "A"+strconv.FormatInt(cursor, 10), "B"+strconv.FormatInt(cursor, 10), styleGridId)
		cursor++
	}

	return
}

func makeExcelNewSheet(f *excelize.File, list []*models.ExcelTimesheet) (err error) {
	var weekNumber int64 = 0
	var currentSheet string
	var currentBeginningOfWeek time.Time
	var currentDate time.Time
	var currentRow int64 = 1
	var styleGridId int
	styleGridId, err = f.NewStyle(`{
		"border": [{
			"type": "left",
			"color": "000000",
			"style": 4
		},
		{
			"type": "top",
			"color": "000000",
			"style": 4
		},
		{
			"type": "bottom",
			"color": "000000",
			"style": 4
		},
		{
			"type": "right",
			"color": "000000",
			"style": 4
		}]
	}`)
	for _, row := range list {
		if currentBeginningOfWeek != now.With(*row.SigninTime).BeginningOfWeek() {
			weekNumber++

			currentBeginningOfWeek = now.With(*row.SigninTime).BeginningOfWeek()
			currentSheet = "Week " + strconv.FormatInt(weekNumber, 10)
			f.NewSheet(currentSheet)

			f.SetColWidth(currentSheet, "A", "D", 15)
			f.SetColWidth(currentSheet, "B", "B", 30)
			f.SetColWidth(currentSheet, "E", "E", 10)

			f.SetCellValue(currentSheet, "B1", "Username")
			f.SetCellValue(currentSheet, "C1", "Sign In Time")
			f.SetCellValue(currentSheet, "D1", "Sign Out Time")
			f.SetCellValue(currentSheet, "E1", "Daily Total")

			var styleId int
			styleId, err = f.NewStyle(`{
				"font":{
					"bold":true
				},
				"border": [{
					"type": "left",
					"color": "000000",
					"style": 4
				},
				{
					"type": "top",
					"color": "000000",
					"style": 4
				},
				{
					"type": "bottom",
					"color": "000000",
					"style": 4
				},
				{
					"type": "right",
					"color": "000000",
					"style": 4
				}]
			}`)
			f.SetCellStyle(currentSheet, "B1", "E1", styleId)

			currentRow = 3
		}

		if currentDate != now.With(*row.SigninTime).BeginningOfDay() {
			currentDate = now.With(*row.SigninTime).BeginningOfDay()

			if currentRow > 3 {
				currentRow++
			}

			cell := "A" + strconv.FormatInt(currentRow, 10)
			if row.SigninTime == nil {
				f.SetCellValue(currentSheet, cell, "NULL")
			} else {
				f.SetCellValue(currentSheet, cell, row.SigninTime.Format("2006-01-02"))
			}

			var styleId int
			styleId, err = f.NewStyle(`{"fill":{"type":"pattern","color":["#FFFF00"],"pattern":1}}`)
			f.SetCellStyle(currentSheet, cell, "E"+strconv.FormatInt(currentRow, 10), styleId)

			currentRow++
		}

		f.SetCellValue(currentSheet, "A"+strconv.FormatInt(currentRow, 10), row.SigninTime.Weekday())
		f.SetCellValue(currentSheet, "B"+strconv.FormatInt(currentRow, 10), row.Name)
		if row.SigninTime == nil {
			f.SetCellValue(currentSheet, "C"+strconv.FormatInt(currentRow, 10), "NULL")
		} else {
			f.SetCellValue(currentSheet, "C"+strconv.FormatInt(currentRow, 10), row.SigninTime.Format("2006-01-02"))
		}
		if row.SignoutTime == nil {
			f.SetCellValue(currentSheet, "D"+strconv.FormatInt(currentRow, 10), "NULL")
		} else {
			f.SetCellValue(currentSheet, "D"+strconv.FormatInt(currentRow, 10), row.SignoutTime.Format("2006-01-02"))
		}
		f.SetCellValue(currentSheet, "E"+strconv.FormatInt(currentRow, 10), row.Total)

		f.SetCellStyle(currentSheet, "A"+strconv.FormatInt(currentRow, 10), "E"+strconv.FormatInt(currentRow, 10), styleGridId)

		currentRow++
	}

	/*
		var replaceNumer int64 = 1
		for i := weekNumber; i > 0; i-- {
			f.SetSheetName("Week "+strconv.FormatInt(i, 10), "Week #"+strconv.FormatInt(replaceNumer, 10))
			replaceNumer++
		}
	*/

	return
}
