package controllers

import (
	"github.com/kataras/iris/v12"
	"commerce-hsz/services"
	"github.com/kataras/iris/v12/mvc"
	"commerce-hsz/datamodels"
	"commerce-hsz/common"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
}

// 查询全部商品
func (p *ProductController)GetAll() mvc.View {
	productArray,_ := p.ProductService.GetAllProduct()
	return mvc.View{
		Name:"product/view.html",
		Data:iris.Map{
			"productArray": productArray,
		},
	}
}

// 修改商品
func (p *ProductController)PostUpdate() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{
		TagName: "goods",
	})
	err := dec.Decode(p.Ctx.Request().Form, product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err = p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

// 新建商品的输入页面
func (p *ProductController)GetAdd() mvc.View {
	return mvc.View{
		Name:"product/add.html",
	}
}

// 新建商品
func (p *ProductController)PostAdd() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{
		TagName: "goods",
	})

	err := dec.Decode(p.Ctx.Request().Form, product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	_, err = p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	p.Ctx.Redirect("/product/all")
}

// 根据id查询商品
func (p *ProductController)GetManager() mvc.View {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	product, err := p.ProductService.GetProduct(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// 删除商品
func (p *ProductController)GetDelete() {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	ok := p.ProductService.DeleteProduct(id)
	if ok {
		p.Ctx.Application().Logger().Debug("删除商品成功, 商品id为" + idString)
	} else {
		p.Ctx.Application().Logger().Debug("删除商品失败, 商品id为" + idString)
	}

	p.Ctx.Redirect("/product/all")
}














