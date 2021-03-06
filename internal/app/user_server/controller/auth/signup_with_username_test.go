// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package auth_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/user_server/controller/auth"
	"github.com/axetroy/go-server/internal/app/user_server/controller/invite"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func TestSignUpWithEmptyBody(t *testing.T) {
	// empty body
	r := tester.HttpUser.Post("/v1/auth/signup", []byte(nil), nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

	assert.Equal(t, exception.InvalidParams.Code(), res.Status)
	assert.Equal(t, "unexpected end of JSON input", res.Message)
	assert.Nil(t, res.Data)
}

func TestSignUpWithNotFullBody(t *testing.T) {
	username := "username"

	// 没有输入密码
	body, _ := json.Marshal(&auth.SignUpWithUsernameParams{
		Username: username,
	})

	// empty body
	r := tester.HttpUser.Post("/v1/auth/signup", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

	assert.Equal(t, exception.InvalidParams.Code(), res.Status)
	assert.Nil(t, res.Data)
}

func TestSignUpSuccess(t *testing.T) {
	rand.Seed(99) // 重置随机码，否则随机数会一样

	username := "test-TestSignUpSuccess"

	res := auth.SignUpWithUsername(auth.SignUpWithUsernameParams{
		Username: username,
		Password: "123123",
	})

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	defer tester.DeleteUserByUserName(username)

	profile := schema.Profile{}

	assert.Nil(t, res.Decode(&profile))

	// 默认未激活状态
	assert.Equal(t, int(profile.Status), int(model.UserStatusInit))
	assert.Equal(t, profile.Username, username)
	assert.Equal(t, *profile.Nickname, username)
	assert.Equal(t, profile.Role, []string{model.DefaultUser.Name})
	assert.Nil(t, profile.Email)
	assert.Nil(t, profile.Phone)
}

func TestSignUpInviteCode(t *testing.T) {
	rand.Seed(133) // 重置随机码，否则随机数会一样

	testerUsername := "tester"
	testerUid := ""
	username := "test-TestSignUpInviteCode"

	inviteCode := ""

	// 动态创建一个测试账号
	{
		r := auth.SignUpWithUsername(auth.SignUpWithUsernameParams{
			Username: testerUsername,
			Password: "123123",
		})

		profile := schema.Profile{}

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, r.Decode(&profile))

		inviteCode = profile.InviteCode
		testerUid = profile.Id

		defer tester.DeleteUserByUserName(testerUsername)
	}

	rand.Seed(1111) // 重置随机码，否则随机数会一样

	res := auth.SignUpWithUsername(auth.SignUpWithUsernameParams{
		Username:   username,
		Password:   "123123",
		InviteCode: &inviteCode,
	})

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	defer tester.DeleteUserByUserName(username)

	profile := schema.Profile{}

	if !assert.Nil(t, res.Decode(&profile)) {
		return
	}

	// 默认未激活状态
	assert.Equal(t, int(model.UserStatusInit), int(profile.Status))
	assert.Equal(t, username, profile.Username)
	assert.Equal(t, username, *profile.Nickname)
	assert.Nil(t, profile.Email)
	assert.Nil(t, profile.Phone)

	// 获取我的邀请记录
	resInvite := invite.GetByStruct(&model.InviteHistory{Invitee: profile.Id})
	InviteeData := schema.Invite{}

	assert.Nil(t, resInvite.Decode(&InviteeData))
	assert.Equal(t, profile.Id, InviteeData.Invitee)
	assert.Equal(t, testerUid, InviteeData.Inviter)
}
