package tester

import (
	"errors"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_schema"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/user/user_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/util"
)

// 创建一个测试用户
func CreateUser() (profile user_schema.ProfileWithToken, err error) {
	var (
		username  = "test-" + util.RandomString(6)
		password  = "123123"
		ip        = "0.0.0.0"
		userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3788.1 Safari/537.36"
	)

	// 创建用户
	if r := auth.SignUp(auth.SignUpParams{
		Username: &username,
		Password: password,
	}); r.Status != schema.StatusSuccess {
		err = errors.New(r.Message)
		return
	}

	// 登陆获取 token
	r := auth.SignIn(schema.Context{
		UserAgent: userAgent,
		Ip:        ip,
	}, auth.SignInParams{
		Account:  username,
		Password: password,
	})

	if r.Status != schema.StatusSuccess {
		err = errors.New(r.Message)
		return
	}

	if err = Decode(r.Data, &profile); err != nil {
		return
	}

	return
}

// 登陆超级管理员
func LoginAdmin() (profile admin_schema.AdminProfileWithToken, err error) {
	r := admin.Login(admin.SignInParams{
		Username: "admin",
		Password: "admin",
	})

	if r.Status != schema.StatusSuccess {
		err = errors.New(r.Message)
		return
	}

	if err = Decode(r.Data, &profile); err != nil {
		return
	}

	return
}
