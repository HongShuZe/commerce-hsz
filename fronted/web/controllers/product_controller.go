package controllers

import (
	"commerce-hsz/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
	"commerce-hsz/datamodels"
	"os"
	"text/template"
	"path/filepath"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
}

var (
	htmlOutPath = "./fronted/web/htmlProductShow"
	templatePath = "./fronted/web/views/template"
)

func (p *ProductController)GetGenerateHtml()  {
	productString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 获取模板
	contentTemp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	// 获取html生成路径
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	// 获取模板渲染数据
	product, err := p.ProductService.GetProduct(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	// 生成静态文件
	generateStaticHtml(p.Ctx, contentTemp, fileName, product)
}

// 生成html静态文件
func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product)  {
	// 1.判断静态文件是否存在
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
	}

	// 2.生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()
	template.Execute(file, &product)
}

// 判断文件是否存在
func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
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