package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置Course表名
func (a *Pics) TableName() string {
	return PicsTBName()
}

// PicsQueryParam 用于搜索的类
type PicsQueryParam struct {
	BaseQueryParam
	TitleLike string
}

// Pics 实体类
type Pics struct {
	Id        int
	Title      string `orm:"size(128)"`
	Source string `orm:"size(64)"`
	Desc     string `orm:"size(512)"`
	Img       string `orm:"size(256)"`
	StartTime time.Time
	EndTime   time.Time
	CreateTime        time.Time
	Creator   *User `orm:"rel(fk)"` //设置一对多关系
}

// PicsPageList 获取分页数据
func PicsPageList(params *PicsQueryParam) ([]*Pics, int64) {
	query := orm.NewOrm().QueryTable(PicsTBName())
	data := make([]*Pics, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "CreateTime":
		sortorder = "CreateTime"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("title__istartswith", params.TitleLike)
	total, _ := query.Count()
	_, _ = query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// PicsBatchDelete 批量删除
func PicsBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(PicsTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// PicsOne 获取单条
func PicsOne(id int) (*Pics, error) {
	o := orm.NewOrm()
	m := Pics{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
