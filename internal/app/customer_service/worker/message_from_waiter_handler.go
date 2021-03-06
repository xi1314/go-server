package worker

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"log"
	"time"
)

func textMessageFromWaiterHandler(msg ws.Message) (err error) {
	waiterClient := ws.WaiterPoll.Get(msg.From)
	userClient := ws.UserPoll.Get(msg.To)

	if waiterClient == nil || userClient == nil {
		return
	}

	// 发送成功，写入聊天记录
	tx := database.Db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	hash := util.MD5(userClient.UUID + waiterClient.UUID)

	session := model.CustomerSession{
		Id:       hash,
		Uid:      userClient.GetProfile().Id,
		WaiterID: waiterClient.GetProfile().Id,
	}

	// 获取会话
	if err := tx.Model(model.CustomerSession{}).Where(&session).First(&session).Error; err != nil {
		return err
	}

	var raw []byte

	if raw, err = json.Marshal(msg.Payload); err != nil {
		return
	}

	sessionItem := model.CustomerSessionItem{
		SessionID:  session.Id,
		Type:       model.SessionTypeText,
		ReceiverID: userClient.GetProfile().Id,
		SenderID:   waiterClient.GetProfile().Id,
		Payload:    string(raw),
	}

	// 讲聊天记录写入会话
	if err := tx.Create(&sessionItem).Error; err != nil {
		return err
	}

	// 推送给两个端口
	// 不管成功与否，因为服务端已经收到

	// 发送给客户端
	_ = userClient.WriteJSON(ws.Message{
		Id:      session.Id,
		From:    msg.From,
		To:      msg.To,
		Type:    string(ws.TypeResponseUserMessageText),
		Payload: msg.Payload,
		Date:    sessionItem.CreatedAt.Format(time.RFC3339Nano),
	})

	// 给客服端一个回执
	_ = waiterClient.WriteJSON(ws.Message{
		OpID:    msg.OpID,
		Id:      sessionItem.Id,
		Type:    string(ws.TypeResponseWaiterMessageTextSuccess),
		From:    waiterClient.UUID,
		To:      userClient.UUID,
		Payload: msg.Payload,
		Date:    sessionItem.CreatedAt.Format(time.RFC3339Nano),
	})

	return
}

func imageMessageFromWaiterHandler(msg ws.Message) (err error) {
	waiterClient := ws.WaiterPoll.Get(msg.From)
	userClient := ws.UserPoll.Get(msg.To)

	if waiterClient == nil || userClient == nil {
		return
	}

	// 发送成功，写入聊天记录
	tx := database.Db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	hash := util.MD5(userClient.UUID + waiterClient.UUID)

	session := model.CustomerSession{
		Id:       hash,
		Uid:      userClient.GetProfile().Id,
		WaiterID: waiterClient.GetProfile().Id,
	}

	// 获取会话
	if err := tx.Model(model.CustomerSession{}).Where(&session).First(&session).Error; err != nil {
		return err
	}

	var raw []byte

	if raw, err = json.Marshal(msg.Payload); err != nil {
		return
	}

	sessionItem := model.CustomerSessionItem{
		SessionID:  session.Id,
		Type:       model.SessionTypeImage,
		ReceiverID: userClient.GetProfile().Id,
		SenderID:   waiterClient.GetProfile().Id,
		Payload:    string(raw),
	}

	// 讲聊天记录写入会话
	if err := tx.Create(&sessionItem).Error; err != nil {
		return err
	}

	// 推送给两个端口
	// 不管成功与否，因为服务端已经收到

	// 发送给客户端
	_ = userClient.WriteJSON(ws.Message{
		Id:      session.Id,
		From:    msg.From,
		To:      msg.To,
		Type:    string(ws.TypeResponseUserMessageImage),
		Payload: msg.Payload,
		Date:    sessionItem.CreatedAt.Format(time.RFC3339Nano),
	})

	// 给客服端一个回执
	_ = waiterClient.WriteJSON(ws.Message{
		OpID:    msg.OpID,
		Id:      sessionItem.Id,
		Type:    string(ws.TypeResponseWaiterMessageImageSuccess),
		From:    waiterClient.UUID,
		To:      userClient.UUID,
		Payload: msg.Payload,
		Date:    sessionItem.CreatedAt.Format(time.RFC3339Nano),
	})

	return
}

// 处理来之客服端的消息
func MessageFromWaiterHandler() {
	for {
		// 从用户池中取消息
		msg := <-ws.UserPoll.Broadcast

	typeCondition:
		switch ws.TypeRequestWaiter(msg.Type) {
		// 发送文本消息给用户
		case ws.TypeRequestWaiterMessageText:
			if err := textMessageFromWaiterHandler(msg); err != nil {
				log.Println(err)
			}
			break typeCondition
		// 发送图片消息给用户
		case ws.TypeRequestWaiterMessageImage:
			if err := imageMessageFromWaiterHandler(msg); err != nil {
				log.Println(err)
			}
			break typeCondition
		default:
			break typeCondition
		}
	}
}
