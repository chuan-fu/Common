package cdao

import (
	"github.com/chuan-fu/Common/zlog"
	"gorm.io/gorm"
)

// 创建
// 传入指针
func Create(db *gorm.DB, model interface{}) (err error) {
	err = db.Create(model).Error
	if err != nil {
		log.Error(err)
	}
	return
}

// 更新
// 传入指针
// 如columns不传，默认只更新非零值，如传入columns，则更新columns里的列
// model中主键必须存在，且非0，不可全局修改
func CommonUpdate(db *gorm.DB, model interface{}, columns ...string) (err error) {
	if len(columns) > 0 {
		db = db.Select(columns)
	}
	err = db.Model(model).Limit(1).Updates(model).Error
	if err != nil {
		log.Error(err)
	}
	return
}

// 更新
func Update(db *gorm.DB, tableName string, conditions map[string]interface{}, updates map[string]interface{}) (err error) {
	db = db.Table(tableName)
	for k, v := range conditions {
		db = db.Where(k, v)
	}
	err = db.Limit(1).UpdateColumns(updates).Error
	if err != nil {
		log.Error(err)
	}
	return
}

// model 为指针
// id为主键
func FindById(db *gorm.DB, id int64, model interface{}) (err error) {
	err = db.First(model, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
		return
	}
	return nil
}

// 保存所有字段
func Save(db *gorm.DB, model interface{}) (err error) {
	err = db.Limit(1).Save(model).Error
	if err != nil {
		log.Error(err)
	}
	return
}

func Delete(db *gorm.DB, model interface{}) (err error) {
	err = db.Delete(model).Error
	if err != nil {
		log.Error(err)
	}
	return
}
