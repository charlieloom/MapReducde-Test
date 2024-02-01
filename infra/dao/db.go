package dao

import (
	"dockermysql/dal"
	"dockermysql/dal/query"
)

const MySQLDSN = "root:123456@tcp(127.0.0.1:3307)/Shop?charset=utf8mb4&parseTime=True"

func Init() {
	dal.DB = dal.ConnectDB(MySQLDSN).Debug()
	if dal.DB == nil {
		panic("connect db fail")
	}

	query.SetDefault(dal.DB)
}
