package dao

import (
	"context"
	"dockermysql/dal"
	model "dockermysql/dal/model"
	"dockermysql/dal/query"
	model2 "dockermysql/model"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func GetAllproducts(condition *model2.Condition, offest int) ([]*model.Product, error) {
	p := query.Product
	query := p.WithContext(context.Background()).Where()
	if condition.Id != 0 {
		query = query.Where(p.ID.Eq(int32(condition.Id)))
	}
	if condition.Name != "" {
		query = query.Where(p.Name.Eq(condition.Name))
	}
	if condition.Category != "" {
		query = query.Where(p.Category.Eq(string(condition.Category)))
	}
	products, err := query.Debug().Limit(50000).Offset(offest).Find()
	if err != nil {
		return nil, err
	}
	return products, nil
}
func randProduct(idx int, out chan struct{}, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	var products []*model.Product
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Product%d", idx+1)
		description := fmt.Sprintf("Description for Product%d", idx+1)
		category := fmt.Sprintf("Category%d", rand.Intn(5))
		price := rand.Float64() * 10000000
		stockQuantity := rand.Int31()
		countryOfManufacture := fmt.Sprintf("Country%d", rand.Intn(100))
		dateAdded := time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour)
		lastUpdated := time.Now()
		unitsSold := rand.Int31()
		numberOfReviews := rand.Int31()
		averageRating := rand.Float64() * 5
		product := model.Product{
			Name:                 name,
			Description:          description,
			Category:             category,
			Price:                price,
			StockQuantity:        stockQuantity,
			CountryOfManufacture: countryOfManufacture,
			DateAdded:            dateAdded,
			LastUpdated:          lastUpdated,
			UnitsSold:            unitsSold,
			NumberOfReviews:      numberOfReviews,
			AverageRating:        averageRating,
		}
		products = append(products, &product)
	}

	err := query.Product.WithContext(context.Background()).CreateInBatches(products, len(products))
	if err != nil {
		fmt.Printf("create product fail, err: %v", err)
		return
	}
	<-out
}
func Createproduct() {
	query.SetDefault(dal.DB)

	waitGroup := new(sync.WaitGroup)
	in := make(chan struct{}, 40)

	for i := 0; i < 800000; i++ {
		in <- struct{}{}
		waitGroup.Add(1)
		go randProduct(i, in, waitGroup)
	}

	waitGroup.Wait()

	//更新
	// ret, err := query.Product.WithContext(context.Background()).
	// 	Where(query.Product.ID.Eq(1)).
	// 	Update(query.Product.Price, 59.9)
	// if err != nil {
	// 	fmt.Printf("update product fail, err: %v", err)
	// 	return
	// }
	// fmt.Printf("RowsAffected: %v\n", ret.RowsAffected)

	// //查询
	// product2, err := query.Product.WithContext(context.Background()).First()
	// if err != nil {
	// 	fmt.Printf("query product fail, err: %v", err)
	// 	return
	// }
	// fmt.Printf("product %v\n", product2)

	// //删除
	// ret, err = query.Product.WithContext((context.Background())).
	// 	Where(query.Product.ID.Eq(1)).
	// 	Delete()
	// if err != nil {
	// 	fmt.Printf("delete product fail, err: %v", err)
	// 	return
	// }
	// fmt.Printf("RowsAffected: %v\n", ret.RowsAffected)

	// //使用自定义的GetProductsByCategory方法
	// rets, err := query.Product.WithContext(context.Background()).GetProductsByCategory("食品")
	// if err != nil {
	// 	fmt.Printf("GetProductsByCategory  fail, err: %v", err)
	// 	return
	// }
	// for i, b := range rets {
	// 	fmt.Printf("%d:%v\n", i, b)
	// }
}
