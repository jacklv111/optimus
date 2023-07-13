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

	"github.com/golang/mock/gomock"
	"github.com/jacklv111/optimus/pkg/iam/login/do"
	"github.com/jacklv111/optimus/pkg/iam/login/repo"
	loginrepomock "github.com/jacklv111/optimus/pkg/iam/login/repo/mock"
	iamvb "github.com/jacklv111/optimus/pkg/iam/login/value-object"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_login_service(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := loginrepomock.NewMockUserRepoInterface(ctrl)
	repo.UserRepo = mockUserRepo

	Convey("signup, login, parse token", t, func() {
		Convey("Success", func() {
			// signup
			signupInfo := iamvb.SignupInfo{
				Domain:      "gddi",
				Name:        "jack_test",
				Password:    "123456",
				DisplayName: "jacklv",
				Email:       "jacklv@gddi.com.cn",
				Phone:       "12345678901",
				Gender:      "male",
				IsAdmin:     true,
				IsForbidden: false,
			}
			var userDo do.UserDo
			mockUserRepo.EXPECT().CreateUser(gomock.Any()).Do(func(input do.UserDo) {
				userDo = input
			}).Return(nil) // create success

			// signup success
			sigupRes, err := LoginSvc.Signup(signupInfo)
			So(err, ShouldBeNil)
			So(sigupRes.ID, ShouldEqual, userDo.ID)
			So(sigupRes.Name, ShouldEqual, userDo.Name)
			So(sigupRes.Domain, ShouldEqual, userDo.Domain)

			So(userDo.Domain, ShouldEqual, signupInfo.Domain)
			So(userDo.Name, ShouldEqual, signupInfo.Name)
			So(userDo.DisplayName, ShouldEqual, signupInfo.DisplayName)

			// login success
			mockUserRepo.EXPECT().GetUserByUK(gomock.Eq(signupInfo.Domain), gomock.Eq(signupInfo.Name)).Return(userDo, nil)
			tokenString, err := LoginSvc.Login(signupInfo.Domain, signupInfo.Name, signupInfo.Password)
			So(err, ShouldBeNil)

			// parse token
			userInfo, err := LoginSvc.ParseUserInfoFromToken(tokenString)
			So(err, ShouldBeNil)
			So(userInfo.ID, ShouldEqual, userDo.ID)
			So(userInfo.Name, ShouldEqual, userDo.Name)
			So(userInfo.Domain, ShouldEqual, userDo.Domain)
			So(userInfo.IsAdmin, ShouldEqual, userDo.IsAdmin)
			So(userInfo.IsForbidden, ShouldEqual, userDo.IsForbidden)
		})
	})
}
