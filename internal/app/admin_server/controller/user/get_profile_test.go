// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package user_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetProfileWithInvalidAuth(t *testing.T) {

	header := mocker.Header{
		"Authorization": "Bearera 12312", // invalid Bearera
	}

	r := tester.HttpUser.Get("/v1/user/profile", nil, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, exception.InvalidToken.Code(), res.Status)
	assert.Equal(t, exception.InvalidAuth.Error(), res.Message)
}

func TestGetProfileWithInvalidToken(t *testing.T) {
	header := mocker.Header{
		"Authorization": token.Prefix + " 12312",
	}

	r := tester.HttpUser.Get("/v1/user/profile", nil, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, exception.InvalidToken.Code(), res.Status)
	assert.Equal(t, exception.InvalidToken.Error(), res.Message)
}

func TestGetProfileByAdminRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	r := tester.HttpAdmin.Get("/v1/user/"+userInfo.Id, nil, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := schema.Profile{}

	assert.Nil(t, res.Decode(&profile))
	assert.Equal(t, userInfo.Id, profile.Id)
	assert.Equal(t, userInfo.Username, profile.Username)
}
