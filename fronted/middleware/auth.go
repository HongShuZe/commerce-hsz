package middleware

import "github.com/kataras/iris/v12"

func HTTPInterceptor(ctx iris.Context)  {
	uid := ctx.GetCookie("uid")
	if uid == "" {
		ctx.Application().Logger().Debug("必须先登录")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("已经登录")
	ctx.Next()
}