// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package transfer_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/user_server/controller/transfer"
	"github.com/axetroy/go-server/internal/app/user_server/controller/wallet"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDetail(t *testing.T) {
	var log schema.TransferLog
	userFrom, _ := tester.CreateUser()
	userTo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userFrom.Username)
	defer tester.DeleteUserByUserName(userTo.Username)

	// 给账户充钱
	{
		assert.Nil(t, database.Db.Table(wallet.GetTableName("CNY")).Where("id = ?", userFrom.Id).Update(model.Wallet{
			Balance:  100,
			Currency: model.WalletCNY,
		}).Error)
	}

	// 转账一次
	{
		input := transfer.ToParams{
			Currency: "CNY",
			To:       userTo.Id,
			Amount:   "20", // 转账 20
		}

		b, err := json.Marshal(input)

		assert.Nil(t, err)

		signature, err := util.Signature(string(b))

		assert.Nil(t, err)

		res2 := transfer.To(helper.Context{
			Uid: userFrom.Id,
		}, input, signature)

		assert.Equal(t, "", res2.Message)
		assert.Equal(t, schema.StatusSuccess, res2.Status)
		assert.Nil(t, res2.Decode(&log))
	}

	{
		r := transfer.GetDetail(helper.Context{
			Uid: userFrom.Id,
		}, log.Id)

		detail := schema.TransferLog{}

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, r.Decode(&detail))

		assert.Equal(t, log.Id, detail.Id)
		assert.Equal(t, log.From, detail.From)
		assert.Equal(t, log.To, detail.To)
		assert.Equal(t, log.Amount, detail.Amount)
		assert.Equal(t, log.Status, detail.Status)
	}
}
