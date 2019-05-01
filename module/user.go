package module

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smallnest/rpcx/log"
	"time"
)

type User_Level uint8

const (
	User_Level_OK         User_Level = iota //正常用户
	User_Level_DENY                         //被封杀用户
	User_LEVEL_BLACK_LIST                   //黑名单用户
)

type User struct {
	Id          int32
	UserName    string
	Password    string
	Age         uint8
	Level       User_Level
	Create_time time.Time
}

func GetUser(user *User) error {
	orm.Debug = true
	orm := orm.NewOrm()
	//orm.Using("default")
	var err error
	if user.Password != "" {
		orm.QueryTable(user).Filter("id", user.Id).Filter("password", user.Password).One(user)
	} else {
		err = orm.Read(user)
	}
	if err != nil {
		return err
	}
	return nil
}
func Register(user *User) error {
	orm := orm.NewOrm()
	id, err := orm.Insert(user)
	if err != nil {
		return err
	}
	user.Id = int32(id)
	log.Info("注册成功", user)
	return nil
}
