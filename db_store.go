package go_file_uploader

import (
	"github.com/jinzhu/gorm"
)

type dbStore struct {
	db *gorm.DB
}

func (is *dbStore) FileLoad(hash string) (fileModel *FileModel, err error) {
	fileModel = &FileModel{}
	err = is.db.Where(FileModel{Hash: hash}).First(fileModel).Error
	return
}

func (is *dbStore) FileCreate(fileModel *FileModel) error {
	return is.db.Create(fileModel).Error
}

func (is *dbStore) FileExist(hash string) (bool, error) {
	var count uint
	err := is.db.Model(FileModel{}).Where(FileModel{Hash: hash}).Count(&count).Error
	return count > 0, err
}

//
//// ImageSave 方法会判断图片在数据库中是否已经存在如果存在直接放回 否则创建它
//func (is *dbStore) ImageSave(image *Image) (storedImage *Image, exist bool, err error) {
//	originDB := is.db
//	is.db = is.db.Begin()
//	defer func() {
//		if r := recover(); r != nil {
//			is.db.Rollback()
//		}
//		is.db = originDB
//	}()
//	exist, err = is.ImageExist(image.Hash)
//
//	// tx.Commit().Error
//	if err != nil {
//		return nil, false, err
//	} else {
//		if exist {
//			// 图片已经存在
//			storedImage, err = is.ImageLoad(image.Hash)
//			return storedImage, true, err
//		} else {
//			// 图片不存在
//			is.ImageCreate(image)
//		}
//	}
//}

func NewDBStore(db *gorm.DB) Store {
	return &dbStore{db}
}
