package controllers

import (
	"commerce-hsz/datamodels"
	"commerce-hsz/services"
	"commerce-hsz/tool"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"strconv"
	"commerce-hsz/encrypt"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
	Session *sessions.Session
}

// 注册页面
func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

func (c *UserController) PostRegister() {
	nickName := c.Ctx.FormValue("nick_name")
	userName := c.Ctx.FormValue("user_name")
	password := c.Ctx.FormValue("password")

	user := &datamodels.User{
		UserName: userName,
		NickName: nickName,
		Password: password,
	}

	_, err := c.Service.AddUser(user)
	if err != nil {
		c.Ctx.Application().Logger().Debug(err)
		return
	}
	c.Ctx.Redirect("/user/login")
	return

}

//登录页面
func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (c *UserController) PostLogin() mvc.Response {

	userName := c.Ctx.FormValue("user_name")
	password := c.Ctx.FormValue("password")

	user, isOK := c.Service.PwdSuccess(userName, password)
	if !isOK {
		return mvc.Response{
			Path: "/user/login",
		}
	}

	// 将用户id写入cookie中
	tool.GlobalCookie(c.Ctx, "uid", strconv.Itoa(user.ID))
	// 加密用户id
	uidByte := []byte(strconv.Itoa(user.ID))
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		c.Ctx.Application().Logger().Debug(err)
	}
	// todo 不是很清楚c.Session.Set
	//c.Session.Set("uid", strconv.Itoa(user.ID))
	// 将加密的用户id写入cookie中
	tool.GlobalCookie(c.Ctx, "sign", uidString)

	return mvc.Response{
		Path: "/product/detail",
	}
}
