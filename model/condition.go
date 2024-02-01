package Model

type Condition struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	PageSize int    `json:"pagesize"`
	Page     int    `json:"page"`
	Sort     string `json:"sort"`
}
