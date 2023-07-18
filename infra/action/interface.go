/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package action

import (
	"github.com/jacklv111/common-sdk/database"
	"gorm.io/gorm"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type ActionMgrInterface interface {
	Create(resourceType, resourceId, aName, params string) (res ActionDo, err error)
	Updates(res ActionDo) (err error)
	Delete(resourceType, resourceId string) (err error)
	GetUpdateAtLessThan(time int64) (res []ActionDo, err error)
	ExistsByResourceTypeAndResourceId(resourceType, resourceId string) (exists bool, err error)
}

type ActionMgrImpl struct {
}

func (repo *ActionMgrImpl) Create(resourceType, resourceId, aName, params string) (res ActionDo, err error) {
	res = ActionDo{
		ResourceType: resourceType,
		ResourceId:   resourceId,
		Name:         aName,
		Params:       params,
	}

	err = database.Db.Create(&res).Error
	return res, err
}

func (repo *ActionMgrImpl) Updates(res ActionDo) (err error) {
	err = database.Db.Model(&res).Updates(&res).Error
	return err
}

func (repo *ActionMgrImpl) Delete(resourceType, resourceId string) (err error) {
	err = database.Db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceId).Delete(&ActionDo{}).Error
	return err
}

func (repo *ActionMgrImpl) GetUpdateAtLessThan(time int64) (res []ActionDo, err error) {
	err = database.Db.Where("update_at < ?", time).Find(&res).Error
	return res, err
}

func (repo *ActionMgrImpl) ExistsByResourceTypeAndResourceId(resourceType, resourceId string) (exists bool, err error) {
	var res ActionDo
	err = database.Db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceId).First(&res).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

var ActionMgr ActionMgrInterface

func init() {
	ActionMgr = &ActionMgrImpl{}
}
