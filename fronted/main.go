package main

import (
	"github.com/kataras/iris/v12"
	"commerce-hsz/common"
	"github.com/opentracing/opentracing-go/log"
	"context"
	"github.com/kataras/iris/v12/sessions"
	"time"
	"commerce-hsz/repositories"
	"commerce-hsz/services"
	"github.com/kataras/iris/v12/mvc"
	"commerce-hsz/fronted/web/controllers"
	"commerce-hsz/fronted/middleware"
)

func main()  {
	app := iris.New()

	app.Logger().SetLevel("debug")

	tmplate := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)

	app.HandleDir("/public", "./fronted/web/public")
	app.HandleDir("/html", "./fronted/web/htmlProductShow")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message","访问页面出错!"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	sess := sessions.New(sessions.Config{
		Cookie:"AdminCookie",
		Expires:600*time.Minute,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userManager := repositories.NewUserRepositry("tbl_user", db)
	userService := services.NewService(userManager)
	user := mvc.New(app.Party("/user"))
	user.Register(userService, ctx, sess.Start)
	user.Handle(new(controllers.UserController))

	productManager := repositories.NewProductManager("tbl_product", db)
	proService := services.NewProductService(productManager)
	orderManager := repositories.NewOrderManager("tbl_order", db)
	orderService := services.NewOrderService(orderManager)
	pro := app.Party("/product")
	product := mvc.New(pro)
	pro.Use(middleware.HTTPInterceptor)
	product.Register(proService, orderService)
	product.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("localhost:8088"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
