// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

// Get Message detail
func GetByAdmin(c helper.Context, id string) (res schema.Response) {
	var (
		err  error
		data schema.MessageAdmin
		tx   *gorm.DB
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	tx = database.Db.Begin()

	adminInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	MessageInfo := model.Message{
		Id: id,
	}

	if err = tx.Last(&MessageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(MessageInfo, &data.MessagePureAdmin); err != nil {
		return
	}

	if MessageInfo.ReadAt != nil {
		readAt := MessageInfo.ReadAt.Format(time.RFC3339Nano)
		data.Read = true
		data.ReadAt = &readAt
	}

	data.CreatedAt = MessageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = MessageInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

// 管理员获取个人消息详情
var GetAdminRouter = router.Handler(func(c router.Context) {
	id := c.Param(ParamsIdName)

	c.ResponseFunc(nil, func() schema.Response {
		return GetByAdmin(helper.NewContext(&c), id)
	})
})
