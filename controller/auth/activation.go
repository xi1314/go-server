package auth

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"net/http"
)

type ActivationParams struct {
	Code string `valid:"Required;";json:"code"`
}

func Activation(context *gin.Context) {
	var (
		input   ActivationParams
		err     error
		session *xorm.Session
		tx      bool
	)

	defer func() {

		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("unknown error")
			}
		}

		if tx {
			if err != nil {
				_ = session.Rollback()
			} else {
				err = session.Commit()
			}
		}

		if session != nil {
			session.Close()
		}

		if err != nil {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})
		} else {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusSuccess,
				Message: "",
				Data:    nil,
			})
		}
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	// TODO: 校验参数

	var (
		uid string
	)

	if uid, err = redis.ActivationCode.Get(input.Code).Result(); err != nil {
		err = exception.InvalidActiveCode
		return
	}

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	defer func() {
		if err != nil {
			_ = session.Rollback()
		} else {
			_ = session.Commit()
		}
	}()

	user := model.User{Id: uid}

	var isExist bool

	if isExist, err = session.Get(&user); err != nil {
		return
	}

	if isExist == false {
		err = exception.UserNotExist
		return
	}

	if user.Status != model.UserStatusInactivated {
		err = exception.UserHaveActive
		return
	}

	user.Status = model.UserStatusInit

	// 指定更新这个字段
	if _, err = orm.Db.Id(user.Id).Cols("status").Update(&user); err != nil {
		return
	}

	// delete code from redis
	if err = redis.ActivationCode.Del(input.Code).Err(); err != nil {
		return
	}
}
