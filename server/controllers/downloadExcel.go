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
	Week int64 `form:"week" binding:"required"`
}

// @Summary DownloadExcel
// @Description Download Excel
// @Tags After Authorization
// @accept plain
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @param week query int64 false "Week Count"
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

	store, ok := database.GetStore(claims.StoreID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Store error!",
		})
		return
	}

	var timesheet *models.Timesheet
	list, err := timesheet.GetListForExcel(*store)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Timesheet error!",
		})
		return
	}

	f := excelize.NewFile()
	makeExcel(f, list, form.Week)
	if err := f.SaveAs("export.xlsx"); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "file create error!",
		})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=export.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.File("export.xlsx")
}

type Subtotal struct {
	Name string
	Sum  float64
}

func makeExcel(f *excelize.File, list []*models.Timesheet, week int64) {
	f.SetCellValue("Sheet1", "B1", "Username")
	f.SetCellValue("Sheet1", "C1", "Sign In Time")
	f.SetCellValue("Sheet1", "D1", "Sign Out Time")
	f.SetCellValue("Sheet1", "E1", "Daily Total")
	var current int64 = 2
	var isBeignningSundayInit bool = false
	var theBeignningSunday time.Time
	var subtotals = make(map[uint]*Subtotal)
	for _, row := range list {
		if week < 0 {
			break
		}
		if !isBeignningSundayInit {
			isBeignningSundayInit = true
			theBeignningSunday = now.With(*row.SigninTime).BeginningOfWeek()
			f.SetCellValue("Sheet1", "A"+strconv.FormatInt(current, 10), "Week #"+strconv.FormatInt(week, 10))
			week--
			current++
		}
		if row.SigninTime.Before(theBeignningSunday) {
			theBeignningSunday = now.With(*row.SigninTime).BeginningOfWeek()
			f.SetCellValue("Sheet1", "A"+strconv.FormatInt(current, 10), "Week #"+strconv.FormatInt(week, 10))
			week--
			current++
		}
		f.SetCellValue("Sheet1", "A"+strconv.FormatInt(current, 10), row.SigninTime.Weekday())
		f.SetCellValue("Sheet1", "B"+strconv.FormatInt(current, 10), row.Username)
		f.SetCellValue("Sheet1", "C"+strconv.FormatInt(current, 10), row.SigninTime)
		f.SetCellValue("Sheet1", "D"+strconv.FormatInt(current, 10), row.SignoutTime)
		f.SetCellValue("Sheet1", "E"+strconv.FormatInt(current, 10), row.Total)
		if _, ok := subtotals[row.UserID]; !ok {
			subtotals[row.UserID] = &Subtotal{
				Name: row.Username,
				Sum:  row.Total,
			}
		} else {
			subtotals[row.UserID].Sum += row.Total
		}
		current++
	}
	current += 1
	f.SetCellValue("Sheet1", "C"+strconv.FormatInt(current, 10), "Total")
	current += 1
	for _, subtotal := range subtotals {
		f.SetCellValue("Sheet1", "B"+strconv.FormatInt(current, 10), subtotal.Name)
		f.SetCellValue("Sheet1", "C"+strconv.FormatInt(current, 10), subtotal.Sum)
		current++
	}
}
