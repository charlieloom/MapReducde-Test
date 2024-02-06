package controller

import (
	"bytes"
	model2 "dockermysql/model"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

type ExportData struct {
}

func (ExportData) Setup(sarama.ConsumerGroupSession) error {

	return nil
}
func (ExportData) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (ExportData) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("export msg rev: topic = %s, partition = %d,offset= %d, timestamp = %v, ", msg.Topic, msg.Partition, msg.Offset, msg.Timestamp)
		var exportMsg model2.ExportMsg
		err := jsoniter.Unmarshal([]byte(msg.Value), &exportMsg)
		if err != nil {
			fmt.Println("unmarshal error:", err)
			return err
		}
		b := make([]byte, 64)
		b = b[:runtime.Stack(b, false)]
		b = bytes.TrimPrefix(b, []byte("goroutine "))
		b = b[:bytes.IndexByte(b, ' ')]
		n, _ := strconv.ParseUint(string(b), 10, 64)

		log.Printf("协程ID %d", n)

		//打开文件

		file, err := os.OpenFile(exportMsg.File, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		if err != nil {
			fmt.Println("open file error :", err)
			return err
		}
		defer file.Close()
		for _, product := range exportMsg.Productlist {
			c := []string{cast.ToString(product.ID), product.Name, product.Description, product.Category, cast.ToString(product.Price), cast.ToString(product.StockQuantity), product.CountryOfManufacture, cast.ToString(product.DateAdded), cast.ToString(product.LastUpdated), cast.ToString(product.UnitsSold), cast.ToString(product.NumberOfReviews), cast.ToString(product.AverageRating)}
			row := strings.Join(c, ",")
			row += "\n"
			file.WriteString(row)
		}

		//打开文件
		// f, err := excelize.OpenFile(exportMsg.File)
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// 	return err
		// }
		// defer func() {
		// 	// Close the spreadsheet.
		// 	if err := f.Close(); err != nil {
		// 		fmt.Println(err)
		// 	}
		// }()
		// for idx, product := range exportMsg.Productlist {
		// 	row := []interface{}{product.ID, product.Name, product.Description, product.Category, product.Price, product.StockQuantity, product.CountryOfManufacture, product.DateAdded, product.LastUpdated, product.UnitsSold, product.NumberOfReviews, product.AverageRating}
		// 	// log.Printf("store &Row sizeof : %d", unsafe.Sizeof(&row))
		// 	if err := f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", idx+1+exportMsg.Row), &row); err != nil {
		// 		fmt.Println("Error setting sheet row:", err)
		// 		break
		// 	}
		// }
		// if err = f.Save(); err != nil {
		// 	fmt.Println("Error saving file", err)
		// 	return err
		// }
		// f.Close()
		// runtime.GC()
		log.Printf("保存成功: topic = %s, partition = %d, offset= %d, timestamp = %v, ", msg.Topic, msg.Partition, msg.Offset, msg.Timestamp)
		session.MarkMessage(msg, "")
	}
	return nil
}
