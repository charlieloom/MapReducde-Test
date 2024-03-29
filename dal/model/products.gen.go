// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameProduct = "Products"

// Product mapped from table <Products>
type Product struct {
	ID                   int32     `gorm:"column:ID;primaryKey;autoIncrement:true;comment:商品ID，主键，自增" json:"ID"`   // 商品ID，主键，自增
	Name                 string    `gorm:"column:Name;comment:商品名称" json:"Name"`                                   // 商品名称
	Description          string    `gorm:"column:Description;comment:商品描述" json:"Description"`                     // 商品描述
	Category             string    `gorm:"column:Category;comment:商品类别" json:"Category"`                           // 商品类别
	Price                float64   `gorm:"column:Price;comment:商品价格，精确到小数点后两位" json:"Price"`                       // 商品价格，精确到小数点后两位
	StockQuantity        int32     `gorm:"column:StockQuantity;comment:商品库存数量" json:"StockQuantity"`               // 商品库存数量
	CountryOfManufacture string    `gorm:"column:CountryOfManufacture;comment:商品生产国家" json:"CountryOfManufacture"` // 商品生产国家
	DateAdded            time.Time `gorm:"column:DateAdded;comment:商品上架时间" json:"DateAdded"`                       // 商品上架时间
	LastUpdated          time.Time `gorm:"column:LastUpdated;comment:商品信息最后更新时间" json:"LastUpdated"`               // 商品信息最后更新时间
	UnitsSold            int32     `gorm:"column:UnitsSold;comment:商品销售数量" json:"UnitsSold"`                       // 商品销售数量
	NumberOfReviews      int32     `gorm:"column:NumberOfReviews;comment:商品评价数量" json:"NumberOfReviews"`           // 商品评价数量
	AverageRating        float64   `gorm:"column:AverageRating;comment:商品平均评分，精确到小数点后两位" json:"AverageRating"`     // 商品平均评分，精确到小数点后两位
}

// TableName Product's table name
func (*Product) TableName() string {
	return TableNameProduct
}
