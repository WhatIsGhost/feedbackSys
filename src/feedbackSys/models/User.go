package models

import (
	"github.com/astaxie/beego/orm"
)

// TableName 设置User表名
func (a *User) TableName() string {
	return UserTBName()
}

// UserQueryParam 用于查询的类
type UserQueryParam struct {
	BaseQueryParam
	UserNameLike string //模糊查询
	RealNameLike string //模糊查询
	Mobile       string //精确查询
	SearchStatus string //为空不查询，有值精确查询
}

// User 实体类
type User struct {
	Id                 int
	RealName           string `orm:"size(32)"`
	UserName           string `orm:"size(24)"`
	UserPwd            string `json:"-"`
	IsSuper            bool
	Status             int
	Mobile             string                `orm:"size(16)"`
	Email              string                `orm:"size(256)"`
	Avatar             string                `orm:"size(256)"`
	RoleIds            []int                 `orm:"-" form:"RoleIds"`
	RoleUserRel []*RoleUserRel `orm:"reverse(many)"` // 设置一对多的反向关系
	ResourceUrlForList []string              `orm:"-"`
}

// UserPageList 获取分页数据
func UserPageList(params *UserQueryParam) ([]*User, int64) {
	query := orm.NewOrm().QueryTable(UserTBName())
	data := make([]*User, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("username__istartswith", params.UserNameLike)
	query = query.Filter("realname__istartswith", params.RealNameLike)
	if len(params.Mobile) > 0 {
		query = query.Filter("mobile", params.Mobile)
	}
	if len(params.SearchStatus) > 0 {
		query = query.Filter("status", params.SearchStatus)
	}
	total, _ := query.Count()
	_, _ = query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// UserOne 根据id获取单条
func UserOne(id int) (*User, error) {
	o := orm.NewOrm()
	m := User{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// UserOneByUserName 根据用户名密码获取单条
func UserOneByUserName(username, userpwd string) (*User, error) {
	m := User{}
	err := orm.NewOrm().QueryTable(UserTBName()).Filter("username", username).Filter("userpwd", userpwd).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
