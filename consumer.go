package main

import (
	"commerce-hsz/common"
	"fmt"
	"commerce-hsz/repositories"
	"commerce-hsz/services"
	"commerce-hsz/rabbitmq"
)

func main()  {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}

	product := repositories.NewProductManager("tbl_product", db)
	proService := services.NewProductService(product)

	order := repositories.NewOrderManager("tbl_order", db)
	orderService := services.NewOrderService(order)

	//启动消费者
	rmqConsumeSimple := rabbitmq.NewRabbitMQSimple("order_product")
	rmqConsumeSimple.ConsumeSimple(orderService, proService)
}