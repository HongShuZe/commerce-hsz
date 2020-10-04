package controllers

import (
	"commerce-hsz/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
	"commerce-hsz/datamodels"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
}

// 获取商品信息
func (p *ProductController)GetDetail() mvc.View {
	product, err := p.ProductService.GetProduct(1)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name: "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// 创建订单
func (p *ProductController)GetOrder() mvc.View {
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	product, err := p.ProductService.GetProduct(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	var orderID int64
	showMessage := "抢购失败"
	// 判断商品数量是否满足
	if product.ProductNum > 0 {
		// 更新商品数量
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		// 创建订单
		userID, err := strconv.Atoi(userString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}

		order := &datamodels.Order{
			UserID: userID,
			ProductID: productID,
			OrderNum: 1,
		}

		orderID, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功"
		}
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name: "product/result.html",
		Data: iris.Map{
			"orderID": orderID,
			"showMessage": showMessage,
		},
	}
}