package module

import (
	"github.com/aloxc/gappuser/config"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smallnest/rpcx/log"
	"os"
	"time"
)

func init() {
	dbName := os.Getenv(config.USER_MYSQL_DATABASE_NAME)
	host := os.Getenv(config.USER_MYSQL_HOST)
	user := os.Getenv(config.USER_MYSQL_USER)
	password := os.Getenv(config.USER_MYSQL_PASSWORD)
	if dbName == "" {
		dbName = config.USER_MYSQL_DATABASE_NAME_DEFAULT
	}
	if host == "" {
		host = config.USER_MYSQL_HOST_DEFAULT
	}
	if user == "" {
		user = config.USER_MYSQL_USER_DEFAULT
	}
	if password == "" {
		password = config.USER_MYSQL_PASSWORD_DEFAULT
	}

	logs.Info("等待10秒后准备初始化数据库表")
	time.Sleep(time.Second * 10)
	ds := user + ":" + password + "@tcp(" + host + ")/" + dbName + "?charset=utf8mb4&loc=Local"
	logs.Info(ds)
	// set default database
	orm.RegisterDataBase("default", "mysql", ds, 30)
	//orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:33066)/user?charset=utf8", 30)

	orm.RegisterModel(new(User))
	// create table
	err := orm.RunSyncdb("default", true, true)
	if err != nil {
		log.Info("启动创建数据连接异常", err)
	}
}
