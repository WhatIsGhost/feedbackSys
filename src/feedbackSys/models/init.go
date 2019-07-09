package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// init 初始化
func init() {
	orm.RegisterModel(new(User), new(Pics), new(RoleResourceRel), new(Resource), new(Role), new(RoleUserRel))
}

// TableName 下面是统一的表名管理
func TableName(name string) string {
	prefix := beego.AppConfig.String("db_dt_prefix")
	return prefix + name
}

// UserTBName 获取 User 对应的表名称
func UserTBName() string {
	return TableName("user")
}


// RoleTBName 获取 Role 对应的表名称
func RoleTBName() string {
	return TableName("role")
}

// RoleUserRelTBName 角色与用户多对多关系表
func RoleUserRelTBName() string {
	return TableName("role_user_rel")
}

// ResourceTBName 获取 Resource 对应的表名称
func ResourceTBName() string {
	return TableName("resource")
}

// RoleResourceRelTBName 角色与资源多对多关系表
func RoleResourceRelTBName() string {
	return TableName("role_resource_rel")
}

func PicsTBName() string {
	return TableName("pics")
}

