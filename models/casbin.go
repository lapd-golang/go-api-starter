package models

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"go-admin-starter/utils/config"
)

//权限结构
type CasbinModel struct {
	Ptype    string `json:"ptype"`
	RoleName string `json:"rolename"`
	Path     string `json:"path"`
	Method   string `json:"method"`
}

//添加权限
func (c *CasbinModel) AddCasbin(cm CasbinModel) bool {
	e := Casbin()
	return e.AddPolicy(cm.RoleName, cm.Path, cm.Method)

}

//持久化到数据库
func Casbin() *casbin.Enforcer {
	conf := config.New()

	a := gormadapter.NewAdapter(conf.Database.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Name), true)
	e := casbin.NewEnforcer("conf/auth_model.conf", a)
	e.LoadPolicy()

	return e
}
