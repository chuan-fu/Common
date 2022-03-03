package cdao

import (
	"fmt"
	"testing"

	"github.com/chuan-fu/Common/db/mysql"
	"github.com/chuan-fu/Common/zlog"
)

// 检测字典表
type CheckMap struct {
	ID        int64  `gorm:"column:id;primary_key" json:"id"` // id
	K         int    `gorm:"column:k" json:"k"`               // 1委托单位，2检测单位
	V         string `gorm:"column:v" json:"v"`               // value
	IsDeleted int    `json:"is_deleted" gorm:"column:is_deleted"`
}

// TableName sets the insert table name for this struct type
func (d *CheckMap) TableName() string {
	return "check_map"
}

func init() {
	err := mysql.ConnectGORM(mysql.MysqlConf{
		DataSourceName: "root:123456@tcp(0.0.0.0:3306)/dpm?charset-utf8mb4",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestCommonUpdate(t *testing.T) {
	err := CommonUpdate(mysql.GetGorm(), &CheckMap{
		ID:        1,
		K:         11,
		V:         "杭州老爸评测科技有限公司1",
		IsDeleted: 0,
	}, "k", "v", "is_deleted")
	fmt.Println(err)
}
