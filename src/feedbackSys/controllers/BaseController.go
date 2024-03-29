package controllers

import (
	"feedbackSys/enums"
	"feedbackSys/models"
	"feedbackSys/utils"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
)

type BaseController struct {
	beego.Controller
	controllerName string         //当前控制名称
	actionName     string         //当前action名称
	curUser        models.User //当前用户信息
}

func (c *BaseController) Prepare() {
	//附值
	c.controllerName, c.actionName = c.GetControllerAndAction()
	//从Session里获取数据 设置用户信息
	c.adapterUserInfo()
}

func (c *BaseController) adapterUserInfo() {
	a := c.GetSession("user")
	if a != nil {
		c.curUser = a.(models.User)
		c.Data["user"] = a
	}
}

func (c *BaseController) jsonResult(code enums.JsonResultCode, msg string, obj interface{}) {
	r := &models.JsonResult{code, msg, obj}
	c.Data["json"] = r
	c.ServeJSON()
	c.StopRun()
}

// 重定向
func (c *BaseController) redirect(url string) {
	c.Redirect(url, 302)
	c.StopRun()
}

// 重定向 去错误页
func (c *BaseController) pageError(msg string) {
	errorurl := c.URLFor("HomeController.Error") + "/" + msg
	c.Redirect(errorurl, 302)
	c.StopRun()
}

// 重定向 去登录页
func (c *BaseController) pageLogin() {
	url := c.URLFor("HomeController.Login")
	c.Redirect(url, 302)
	c.StopRun()
}

func (c *BaseController) checkLogin() {
	if c.curUser.Id == 0 {
		//登录页面地址
		urlstr := c.URLFor("HomeController.Login") + "?url="
		//登录成功后返回的址为当前
		returnURL := c.Ctx.Request.URL.Path
		//如果ajax请求则返回相应的错码和跳转的地址
		if c.Ctx.Input.IsAjax() {
			//由于是ajax请求，因此地址是header里的Referer
			returnURL := c.Ctx.Input.Refer()
			c.jsonResult(enums.JRCode302, "请登录", urlstr+returnURL)
		}
		c.Redirect(urlstr+returnURL, 302)
		c.StopRun()
	}
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if c.curUser.Id == 0 {
		return false
	}
	//从session获取用户信息
	user := c.GetSession("user")
	//类型断言
	v, ok := user.(models.User)
	if ok {
		//如果是超级管理员，则直接通过
		if v.IsSuper == true {
			return true
		}
		//遍历用户所负责的资源列表
		for i, _ := range v.ResourceUrlForList {
			urlfor := strings.TrimSpace(v.ResourceUrlForList[i])
			if len(urlfor) == 0 {
				continue
			}
			// TestController.Get,:last,xie,:first,asta
			strs := strings.Split(urlfor, ",")
			if len(strs) > 0 && strs[0] == (ctrlName+"."+ActName) {
				return true
			}
		}
	}
	return false
}

func (c *BaseController) checkAuthor(ignores ...string) {
	//先判断是否登录
	c.checkLogin()
	//如果Action在忽略列表里，则直接通用
	for _, ignore := range ignores {
		if ignore == c.actionName {
			return
		}
	}
	hasAuthor := c.checkActionAuthor(c.controllerName, c.actionName)
	if !hasAuthor {
		utils.LogDebug(fmt.Sprintf("author control: path=%s.%s userid=%v  无权访问", c.controllerName, c.actionName, c.curUser.Id))

		if c.Ctx.Input.IsAjax() {
			c.jsonResult(enums.JRCode403, "无权访问", "")
		} else {
			c.pageError("无权访问")
		}

	}
}

//SetBackendUser2Session 获取用户信息（包括资源UrlFor）保存至Session
func (c *BaseController) setUser2Session(userId int) error {
	m, err := models.UserOne(userId)
	if err != nil {
		return err
	}
	//获取这个用户能获取到的所有资源列表
	resourceList := models.ResourceTreeGridByUserId(userId, 1000)
	for _, item := range resourceList {
		m.ResourceUrlForList = append(m.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	}
	c.SetSession("user", *m)
	return nil
}

func (c *BaseController) setTpl(template ...string) {
	var tplName string
	layout := "shared/layout_page.html"
	switch {
	case len(template) == 1:
		tplName = template[0]
	case len(template) == 2:
		tplName = template[0]
		layout = template[1]
	default:
		//不要Controller这个10个字母
		ctrlName := strings.ToLower(c.controllerName[0 : len(c.controllerName)-10])
		actionName := strings.ToLower(c.actionName)
		tplName = ctrlName + "/" + actionName + ".html"
	}
	c.Layout = layout
	c.TplName = tplName
}