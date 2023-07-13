/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package adapter

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/jacklv111/common-sdk/database"
)

var Adapter *gormadapter.Adapter

func InitCasbinGormAdapter() (err error) {
	Adapter, err = gormadapter.NewAdapterByDB(database.Db)
	return err
}
