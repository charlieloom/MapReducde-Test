package controller

import (
	dao "dockermysql/Dao"
	model "dockermysql/dal/model"
	model2 "dockermysql/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func GetAllproducts(c *gin.Context) {
	var condition model2.Condition
	condition.Id, _ = strconv.Atoi(c.Query("id"))
	condition.Name = c.Query("name")
	condition.Category = c.Query("category")
	condition.Page, _ = strconv.Atoi(c.Query("page"))
	condition.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	condition.Sort = c.Query("sort")

	start := time.Now()
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	selectDone := 0
	var err error
	var productlists []*model.Product
	offest := (condition.Page - 1) * condition.PageSize
	for {
		productlists, err = dao.GetAllproducts(&condition, offest+selectDone)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(productlists) <= 0 {
			break
		}
		for idx, product := range productlists {
			row := []interface{}{product.ID, product.Name, product.Description, product.Category, product.Price, product.StockQuantity, product.CountryOfManufacture, product.DateAdded, product.LastUpdated, product.UnitsSold, product.NumberOfReviews, product.AverageRating}
			for j, value := range row {
				cell, err := excelize.CoordinatesToCellName(j+1, idx+1+selectDone)
				if err != nil {
					fmt.Println(err)
					return
				}
				err = f.SetCellValue("Sheet1", cell, value)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		if err = f.SaveAs("products.xlsx"); err != nil {
			fmt.Println(err)
		}
		selectDone += len(productlists)
		if selectDone >= condition.PageSize {
			break
		}
	}
	cost := time.Since(start)
	fmt.Printf("花费时间：[%s]", cost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": productlists,
	})
}
