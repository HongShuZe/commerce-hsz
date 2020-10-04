package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"commerce-hsz/common"
	"github.com/opentracing/opentracing-go/log"
	"context"
	"commerce-hsz/repositories"
	"commerce-hsz/services"
	"commerce-hsz/backend/web/controllers"
)

func main()  {
	app := iris.New()

	app.Logger().SetLevel("debug")

	tmplate := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)

	app.HandleDir("/assets", "./backend/web/assets")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message","访问页面出错!"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	proRepository := repositories.NewProductManager("tbl_product", db)
	proService := services.NewProductService(proRepository)
	pro := mvc.New(app.Party("/product"))
	pro.Register(ctx, proService)
	pro.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManager("tbl_order", db)
	orderService := services.NewOrderService(orderRepository)
	order := mvc.New(app.Party("/order"))
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}

