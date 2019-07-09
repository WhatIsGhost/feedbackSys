package models

import "time"

// RoleBackendUserRel 角色与用户关系
type RoleUserRel struct {
	Id          int
	Role        *Role        `orm:"rel(fk)"`  //外键
	User *User `orm:"rel(fk)" ` // 外键
	Created     time.Time    `orm:"auto_now_add;type(datetime)"`
}

// TableName 设置表名
func (a *RoleUserRel) TableName() string {
	return RoleUserRelTBName()
}
