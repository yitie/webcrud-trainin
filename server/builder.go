package server

import (
	"fmt"
	"gorm.io/gorm"
	"webOrder/dao"
	"webOrder/dao/config"
	"webOrder/server/router"
	"webOrder/service"
)

func NewServer(dst ...interface{}) *router.Router {
	orm := config.NewOrm()
	// create(orm, dst...)
	migrate(orm, dst...)

	orderDao := dao.NewOrderDAO(orm)

	orderService := service.NewOrderService(orderDao)

	return &router.Router{
		RouterServer: orderService,
	}
}

func create(db *gorm.DB, dst ...interface{}) {
	if err := db.Migrator().CreateTable(dst); err != nil {
		panic(err)
	}
	fmt.Println("hhhh")
}

func migrate(db *gorm.DB, dst ...interface{}) {
	fmt.Println("migrate", len(dst))
	if err := db.AutoMigrate(dst...); err != nil {
		panic(err)
	}
}