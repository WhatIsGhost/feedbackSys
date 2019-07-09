package controllers

import (
	"encoding/json"
	"feedbackSys/minio"
	"time"

	"feedbackSys/enums"
	"feedbackSys/models"
	"feedbackSys/utils"

	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

//PicsController 课程管理
type PicsController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *PicsController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList", "UpdateSeq", "UploadImage")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

//Index 角色管理首页
func (c *PicsController) Index() {
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "pics/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "pics/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("PicsController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("PicsController", "Delete")
}

// DataGrid 课程管理首页 表格获取数据
func (c *PicsController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.PicsQueryParam
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.PicsPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

//Edit 添加、编辑课程界面
func (c *PicsController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.Pics{Id: Id, Source: "Web", StartTime: time.Now(), EndTime: time.Now(), CreateTime: time.Now()}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["hasImg"] = len(m.Img) > 0
	c.Data["m"] = m
	c.setTpl("pics/edit.html", "shared/layout_page.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "pics/edit_headcssjs.html"
	c.LayoutSections["footerjs"] = "pics/edit_footerjs.html"

	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor("PicsController.Index")
}

//Save 添加、编辑页面 保存
func (c *PicsController) Save() {
	var err error
	m := models.Pics{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "提交表单数据失败，可能原因："+err.Error(), m.Id)
	}
	o := orm.NewOrm()
	if m.Id == 0 {
		m.Creator = &c.curUser
		m.CreateTime = time.Now()
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败，可能原因："+err.Error(), m.Id)
		}

	} else {
		oM, err := models.PicsOne(m.Id)
		oM.Title = m.Title
		oM.Source = m.Source
		oM.Desc = m.Desc
		oM.StartTime = m.StartTime
		oM.EndTime = m.EndTime
		oM.Img = m.Img
		if _, err = o.Update(oM); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			utils.LogDebug(err)
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}

}

//Delete 批量删除
func (c *PicsController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.PicsBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

func (c *PicsController) UpdateDesc() {
	Id, _ := c.GetInt("pk", 0)
	oM, err := models.PicsOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}
	value := c.GetString("value", "Northing")
	oM.Desc = value
	o := orm.NewOrm()
	if _, err := o.Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "修改失败", oM.Id)
	}
}
func (c *PicsController) UploadImage() {
	//这里type没有用，只是为了演示传值
	stype, _ := c.GetInt32("type", 0)
	if stype > 0 {
		f, h, err := c.GetFile("fileImg")
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "上传失败", "")
		}
		defer f.Close()
		filePath := "static/upload/" + h.Filename
		// 保存位置在 static/upload, 没有文件夹要先创建
		_ = c.SaveToFile("fileImg", filePath)
		url:=minio.UploadFile("feedbacksys",filePath)
		if url!=""{
			c.jsonResult(enums.JRCodeSucc, "上传成功", url)
		} else{
			c.jsonResult(enums.JRCodeFailed, "上传失败", "")
		}
	} else {
		c.jsonResult(enums.JRCodeFailed, "上传失败", "")
	}
}

//Detail 添加、编辑课程界面
func (c *PicsController) Detail() {
	Id, _ := c.GetInt(":id", 0)
	m := models.Pics{Id: Id, Source: "Web", StartTime: time.Now(), EndTime: time.Now(), CreateTime: time.Now()}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["hasImg"] = len(m.Img) > 0
	c.Data["m"] = m
	c.setTpl("pics/detail.html", "shared/layout_page.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "pics/detail_headcssjs.html"
	c.LayoutSections["footerjs"] = "pics/detail_footerjs.html"

	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor("PicsController.Index")
}
