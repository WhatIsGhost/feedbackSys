package routers

import (
	"feedbackSys/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/pics/index", &controllers.PicsController{}, "*:Index")
	beego.Router("/pics/datagrid", &controllers.PicsController{}, "Get,Post:DataGrid")
	beego.Router("/pics/edit/?:id", &controllers.PicsController{}, "Get,Post:Edit")
	beego.Router("/pics/detail/?:id", &controllers.PicsController{}, "Get,Post:Detail")
	beego.Router("/pics/delete", &controllers.PicsController{}, "Post:Delete")
	beego.Router("/pics/updatedesc", &controllers.PicsController{}, "Post:UpdateDesc")
	beego.Router("/pics/uploadimage", &controllers.PicsController{}, "Post:UploadImage")

	//用户角色路由
	beego.Router("/role/index", &controllers.RoleController{}, "*:Index")
	beego.Router("/role/datagrid", &controllers.RoleController{}, "Get,Post:DataGrid")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "Get,Post:Edit")
	beego.Router("/role/delete", &controllers.RoleController{}, "Post:Delete")
	beego.Router("/role/datalist", &controllers.RoleController{}, "Post:DataList")
	beego.Router("/role/allocate", &controllers.RoleController{}, "Post:Allocate")
	beego.Router("/role/updateseq", &controllers.RoleController{}, "Post:UpdateSeq")

	//资源路由
	beego.Router("/resource/index", &controllers.ResourceController{}, "*:Index")
	beego.Router("/resource/treegrid", &controllers.ResourceController{}, "POST:TreeGrid")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "Get,Post:Edit")
	beego.Router("/resource/parent", &controllers.ResourceController{}, "Post:ParentTreeGrid")
	beego.Router("/resource/delete", &controllers.ResourceController{}, "Post:Delete")
	//快速修改顺序
	beego.Router("/resource/updateseq", &controllers.ResourceController{}, "Post:UpdateSeq")

	//通用选择面板
	beego.Router("/resource/select", &controllers.ResourceController{}, "Get:Select")
	//用户有权管理的菜单列表（包括区域）
	beego.Router("/resource/usermenutree", &controllers.ResourceController{}, "POST:UserMenuTree")
	beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	//后台用户路由
	beego.Router("/user/index", &controllers.UserController{}, "*:Index")
	beego.Router("/user/datagrid", &controllers.UserController{}, "POST:DataGrid")
	beego.Router("/user/edit/?:id", &controllers.UserController{}, "Get,Post:Edit")
	beego.Router("/user/delete", &controllers.UserController{}, "Post:Delete")


	beego.Router("/home/index", &controllers.HomeController{}, "*:Index")
	beego.Router("/home/login", &controllers.HomeController{}, "*:Login")
	beego.Router("/home/dologin", &controllers.HomeController{}, "Post:DoLogin")
	beego.Router("/home/logout", &controllers.HomeController{}, "*:Logout")

	beego.Router("/home/404", &controllers.HomeController{}, "*:Page404")
	beego.Router("/home/error/?:error", &controllers.HomeController{}, "*:Error")

	beego.Router("/", &controllers.HomeController{}, "*:Index")
}
