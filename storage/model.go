package storage

import (
	"gorm.io/gorm"
)

type AVDB struct {
	gorm.Model        // 包含 ID, CreatedAt, UpdatedAt, DeletedAt
	NO         string `gorm:"uniqueIndex;size:64;not null"` // 作品番号，唯一索引
	Title      string `gorm:"size:2048;not null"`           // 原始标题
	ZhCnTitle  string `gorm:"size:2048"`                    // 中文标题
	Pretty     string `gorm:"size:255"`                     // 美化后的显示名称
}

// Insert 基础插入方法
func (avdb *AVDB) Insert() error {
	return GetSqlite().Create(avdb).Error
}

// Update 更新指定字段
func (avdb *AVDB) Update(values interface{}) error {
	return GetSqlite().Model(avdb).Updates(values).Error
}

// Delete 软删除记录
func (avdb *AVDB) Delete() error {
	return GetSqlite().Delete(avdb).Error
}

// GetByID 根据ID查询
func (avdb *AVDB) GetByID(id uint) error {
	return GetSqlite().First(avdb, id).Error
}

// GetByNO 根据番号查询
func (avdb *AVDB) GetByNO(no string) error {
	return GetSqlite().Where("no = ?", no).First(avdb).Error
}

// GetAVDBByID 根据ID查询（静态方法）
func GetAVDBByID(id uint) (*AVDB, error) {
	var avdb AVDB
	err := GetSqlite().First(&avdb, id).Error
	if err != nil {
		return nil, err
	}
	return &avdb, nil
}

// GetAVDBByNO 根据番号查询（静态方法）
func GetAVDBByNO(no string) (*AVDB, error) {
	var avdb AVDB
	err := GetSqlite().Where("no = ?", no).First(&avdb).Error
	if err != nil {
		return nil, err
	}
	return &avdb, nil
}

// SearchByTitle 根据标题模糊查询（支持原始标题和中文标题）
func SearchByTitle(keyword string, limit int) ([]AVDB, error) {
	var avdbs []AVDB
	searchPattern := "%" + keyword + "%"
	err := GetSqlite().Where("title LIKE ? OR zh_cn_title LIKE ?", searchPattern, searchPattern).
		Limit(limit).
		Find(&avdbs).Error
	return avdbs, err
}

// SearchByNO 根据番号模糊查询
func SearchByNO(keyword string, limit int) ([]AVDB, error) {
	var avdbs []AVDB
	searchPattern := "%" + keyword + "%"
	err := GetSqlite().Where("no LIKE ?", searchPattern).
		Limit(limit).
		Find(&avdbs).Error
	return avdbs, err
}
