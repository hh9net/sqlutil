package main

import (
	"fmt"
	"time"

	"github.com/c9s/sqlutils"
)

type OperationLog struct {
	Id         int64     `json:"id"`
	LoginId    int64     `json:"loginid" field:"login_id"`
	UserName   string    `orm:"size(20)" json:"username" field:"user_name"`
	Domain     string    `orm:"size(300)" json:"domain" field:"domain"`
	LoginPhone string    `orm:"size(20)" json:"loginphone" field:"login_phone"`
	Role       string    `orm:"size(20)" json:"role" field:"role"`
	Oper       string    `orm:"size(20)" json:"oper" field:"oper"`
	Service    string    `orm:"size(20)" json:"service" field:"service"`
	ServiceId  string    `orm:"size(20)" json:"serviceid" field:"service_id"`
	Tabpage    string    `orm:"size(100)" json:"tabpage" field:"tabpage"`
	Ipaddr     string    `orm:"size(16)" json:"ipaddr" field:"ipaddr"`
	Useragent  string    `orm:"size(200)" json:"useragent" field:"useragent"`
	Typeid     int       `json:"typeid" field:"typeid"`
	OperTime   time.Time `orm:"type(datetime)" json:"opertime" field:"Oper_time"`
	Details    string    `orm:"size(1024)" json:"details" field:"details"`
	CreatedAt  time.Time `orm:"auto_now;type(datetime)" json:"createdAt" field:"created_at"`
}

func main() {
	dd := make([]OperationLog, 0)
	for i := 0; i < 100; i++ {
		cc := OperationLog{}
		cc.Id = int64(i)
		cc.LoginId = 1002
		cc.UserName = "wangyi"
		cc.Domain = "www.qq.com"
		cc.LoginPhone = "12323456543"
		cc.Role = "经理"
		dd = append(dd, cc)
	}
	//	for _, eee := range dd {
	//	if eee != nil {
	//		fmt.Println(eee)
	//	}
	//	}
	sql, _ := sqlutils.BuildInsertClause2(dd)
	fmt.Println(sql)
}
