// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"time"
)

func userTypeMessageImageHandler(userClient *ws.Client, msg ws.Message) (err error) {
	// 如果还没有认证
	if userClient.GetProfile() == nil {
		return exception.UserNotLogin
	}

	waiterId := ws.MatcherPool.GetMyWaiter(userClient.UUID)

	var body ws.MessageImagePayload

	if err = util.Decode(&body, msg.Payload); err != nil {
		return err
	}

	if err = validator.ValidateStruct(&body); err != nil {
		return err
	}

	// 如果这个客户端没有连接客服，那么消息不会发送
	if waiterId != nil {
		// 把收到的消息广播到客服池
		ws.WaiterPoll.Broadcast <- ws.Message{
			OpID:    msg.OpID,
			From:    userClient.UUID,
			Type:    msg.Type,
			To:      *waiterId,
			Payload: msg.Payload,
			Date:    time.Now().Format(time.RFC3339Nano),
		}
	} else {
		if err = userClient.WriteJSON(ws.Message{
			To:   userClient.UUID,
			Type: string(ws.TypeResponseUserNotConnect),
			Date: time.Now().Format(time.RFC3339Nano),
		}); err != nil {
			return
		}
	}

	return err
}
