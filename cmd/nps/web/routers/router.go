package routers

import (
	"github.com/wanghonggao007/nps/vender/github.com/astaxie/beego"
	"github.com/wanghonggao007/nps/web/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{}, "*:Index")
	beego.AutoRouter(&controllers.IndexController{})
	beego.AutoRouter(&controllers.LoginController{})
	beego.AutoRouter(&controllers.ClientController{})
	beego.AutoRouter(&controllers.AuthController{})
}
