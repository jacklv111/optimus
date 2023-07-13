/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	. "github.com/jacklv111/common-sdk/test"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_MgmtRepo(t *testing.T) {
	defer DbSetUpAndTearDown()()
	Convey("Get first", t, func() {
		domain := "optimus"
		workspace := "test"

		Convey("Success", func() {
			// no record and create
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `resource_management` WHERE `resource_management`.`domain` = ? AND `resource_management`.`workspace` = ? ORDER BY `resource_management`.`id` LIMIT 1")).
				WithArgs(domain, workspace).
				WillReturnRows(sqlmock.NewRows([]string{"id", "domain", "workspace", "created_at"}).AddRow(uuid.New().String(), domain, workspace, 111))

			resMgmtDo, err := ResourceMgmtRepo.GetFirst(domain, workspace)
			So(err, ShouldEqual, nil)
			So(resMgmtDo.Domain, ShouldEqual, domain)
			So(resMgmtDo.Workspace, ShouldEqual, workspace)
			_, err = uuid.Parse(resMgmtDo.ID)
			So(err, ShouldEqual, nil)
		})
	})

	Convey("create", t, func() {
		Convey("Success", func() {
			domain := "optimus"
			workspace := "test"
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `resource_management` (`id`,`domain`,`workspace`,`created_at`) VALUES (?,?,?,?)")).
				WithArgs(sqlmock.AnyArg(), domain, workspace, sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectCommit()
			resMgmtDo, err := ResourceMgmtRepo.Create(domain, workspace)
			So(err, ShouldEqual, nil)
			So(resMgmtDo.Domain, ShouldEqual, domain)
			So(resMgmtDo.Workspace, ShouldEqual, workspace)
			_, err = uuid.Parse(resMgmtDo.ID)
			So(err, ShouldEqual, nil)
		})
	})
}
