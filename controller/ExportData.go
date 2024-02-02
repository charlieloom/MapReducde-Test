package controller

import (
	model2 "dockermysql/model"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/xuri/excelize/v2"
)

var rwMut sync.RWMutex

type ExportData struct {
}

func (ExportData) Setup(sarama.ConsumerGroupSession) error {

	return nil
}
func (ExportData) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (ExportData) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// log.Printf("Message claimed: topic = %s, partition = %d, value = %s, offset= %d, timestamp = %v, ", msg.Topic, msg.Partition, msg.Value, msg.Offset, msg.Timestamp)
		var exportMsg model2.ExportMsg
		err := json.Unmarshal([]byte(msg.Value), &exportMsg)
		if err != nil {
			fmt.Println("unmarshal error:", err)
			return err
		}
		rwMut.Lock()
		//打开文件
		f, err := excelize.OpenFile(exportMsg.File)
		// f := excelize.NewFile()
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		defer func() {
			// Close the spreadsheet.
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()
		//  file, fileHeader, err := req.FormFile(exportMsg.File)

		// if err != nil {
		// 	fmt.Println("Error:", err)
		// 	return err
		// }
		// 打开 Excel 文件
		// file, err := os.Open(exportMsg.File)
		// fmt.Println("msg_file", exportMsg.File)
		// if err != nil {
		// 	fmt.Println("Error opening file:", err)
		// 	return err
		// }
		// defer file.Close() // 释放资源
		// f, err := excelize.OpenReader()
		// fmt.Println("file:", f.Path)
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// 	return err
		// }

		for idx, product := range exportMsg.Productlist {
			row := []interface{}{product.ID, product.Name, product.Description, product.Category, product.Price, product.StockQuantity, product.CountryOfManufacture, product.DateAdded, product.LastUpdated, product.UnitsSold, product.NumberOfReviews, product.AverageRating}

			if err := f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", idx+1+exportMsg.Row), &row); err != nil {
				fmt.Println("Error setting sheet row:", err)
				break
			}
		}
		if err = f.Save(); err != nil {
			fmt.Println("Error saving file", err)
			return err
		}
		f.Close()
		rwMut.Unlock()
		log.Printf("保存成功: topic = %s, partition = %d, offset= %d, timestamp = %v, ", msg.Topic, msg.Partition, msg.Offset, msg.Timestamp)
		session.MarkMessage(msg, "")
	}
	return nil
}
