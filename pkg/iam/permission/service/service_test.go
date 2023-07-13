/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import (
	"testing"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/glebarez/sqlite"
	"github.com/jacklv111/common-sdk/database"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/optimus/infra/casbin/adapter"
	"github.com/jacklv111/optimus/infra/casbin/config"
	"github.com/jacklv111/optimus/pkg/iam/constant"
	pvb "github.com/jacklv111/optimus/pkg/iam/permission/value-object"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"
)

func Test_permission_service(t *testing.T) {
	log.ValidateAndApply(log.LogConfig)
	// Initialize the Gorm adapter with an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to open database: %v", err)
	}
	adapterMem, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule")
	if err != nil {
		t.Errorf("Failed to create adapter: %v", err)
	}
	adapter.Adapter = adapterMem
	config.CasbinConfig.ModelPath = "./test-data/model"

	database.Db = db

	Convey("create admin permission, add role for users, delete role for users, enforce, delete admin permission, enforce", t, func() {
		Convey("Success", func() {
			// create admin permission
			PermissionSvc.CreatePermission(pvb.Permission{
				Domain:       "optimus",
				ResourceType: "dataset",
				ResourceId:   "dataset-1",
				Action:       []string{"*"},
				Effect:       "allow",
			})
			// add role for users
			PermissionSvc.AddRoleForUserInDomain("optimus", constant.ADMIN, "Lydia")
			PermissionSvc.AddRoleForUserInDomain("optimus", constant.ADMIN, "Glebarez")

			// delete role for users
			PermissionSvc.DeleteRoleForUserInDomain("optimus", constant.ADMIN, "Glebarez")

			var rules []gormadapter.CasbinRule
			err := adapterMem.GetDb().Table(adapterMem.GetDb().Statement.Table).Find(&rules).Error
			So(err, ShouldBeNil)
			for _, rule := range rules {
				log.Infow(rule.Ptype, rule.V0, rule.V1, rule.V2, rule.V3, rule.V4, rule.V5)
			}

			// enforce 1
			hasAuth, err := PermissionSvc.Enforce(pvb.PermissionEnforce{
				Domain:       "optimus",
				Name:         "Lydia",
				ResourceType: "dataset",
				ResourceId:   "dataset-1",
				Action:       "get",
			})
			So(err, ShouldBeNil)
			So(hasAuth, ShouldBeTrue)

			// enforce 2
			hasAuth, err = PermissionSvc.Enforce(pvb.PermissionEnforce{
				Domain:       "optimus",
				Name:         "Glebarez",
				ResourceType: "dataset",
				ResourceId:   "dataset-1",
				Action:       "get",
			})
			So(err, ShouldBeNil)
			So(hasAuth, ShouldBeFalse)

			// delete admin permission
			PermissionSvc.DeletePermission(pvb.Permission{
				Domain:       "optimus",
				ResourceType: "dataset",
				ResourceId:   "dataset-1",
				Action:       []string{"*"},
				Effect:       "allow",
			})

			// enforce 3
			hasAuth, err = PermissionSvc.Enforce(pvb.PermissionEnforce{
				Domain:       "optimus",
				Name:         "Lydia",
				ResourceType: "dataset",
				ResourceId:   "dataset-1",
				Action:       "get",
			})

			So(err, ShouldBeNil)
			So(hasAuth, ShouldBeFalse)
		})
	})
}
