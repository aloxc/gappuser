package module

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	logs.Info("准备初始化数据库表")
	// set default database
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:33066)/test?charset=utf8", 30)
	// register model
	orm.RegisterModel(new(User))
	// create table
	orm.RunSyncdb("default", false, true)
}
