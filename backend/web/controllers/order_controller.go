package controllers

import (
	"github.com/kataras/iris/v12"
	"commerce-hsz/services"
	"github.com/kataras/iris/v12/mvc"
	"commerce-hsz/datamodels"
	"commerce-hsz/common"
	"strconv"
)

type OrderController struct {
	Ctx iris.Context
	OrderService services.IOrderService
}

// 查询全部订单
func (o *OrderController) GetAll() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("查询订单信息失败")
	}
	o.Ctx.Application().Logger().Debug("查询全部订单成功")
	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray, //前端遍历order
		},
	}
}

// 修改订单
func (o *OrderController) PostUpdate() {
	order := &datamodels.Order{}
	o.Ctx.Request().ParseForm() // 填充r.Form and r.PostForm.
	dec := common.NewDecoder(&common.DecoderOptions{
		TagName: "sql",
	})
	err := dec.Decode(o.Ctx.Request().Form, order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	err = o.OrderService.UpdateOrder(order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	o.Ctx.Redirect("/order/all")
}

// 新建订单输入页面
func (o *OrderController)GetAdd() mvc.View {
	return mvc.View{
		Name:"order/add.html",
	}
}

// 新建订单
func (o *OrderController)PostAdd() {
	order := &datamodels.Order{}
	o.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{
		TagName:"sql",
	})

	err := dec.Decode(o.Ctx.Request().Form, order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	_, err = o.OrderService.InsertOrder(order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	o.Ctx.Redirect("/order/all")
}

// 根据id查询订单
func (o *OrderController)GetManager() mvc.View {
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 10)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	order, err := o.OrderService.GetOrder(id)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name:"order/manager.html",
		Data:iris.Map{
			"order": order,
		},
	}
}

// 删除订单
func (o *OrderController)GetDelete()  {
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 10)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	ok := o.OrderService.DeleteOrder(id)
	if ok {
		o.Ctx.Application().Logger().Debug("删除订单成功, 商品id为" + idString)
	} else {
		o.Ctx.Application().Logger().Debug("删除订单失败, 商品id为" + idString)
	}
	o.Ctx.Redirect("/order/all")
}