// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package system

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"os/user"
	"runtime"
	"time"
)

type Info struct {
	Username         string         `json:"username"`            // 当前用户名
	Host             host.InfoStat  `json:"host"`                // 操作系统信息
	Avg              load.AvgStat   `json:"avg"`                 // 负载信息
	Arch             string         `json:"arch"`                // 系统架构, 32/64位
	CPU              []cpu.InfoStat `json:"cpu"`                 // CPU信息
	RAMAvailable     uint64         `json:"ram_available"`       // 系统内存是否可供程序使用
	RAMTotal         uint64         `json:"ram_total"`           // 总内存大小
	RAMFree          uint64         `json:"ram_free"`            // 目前可用内存
	RAMUsedBy        uint64         `json:"ram_used_by"`         // 程序占用的内存
	RAMUsedByPercent float64        `json:"ram_used_by_percent"` // 程序占用的内存百分比
	Time             string         `json:"time"`                // 系统当前时间
	Timezone         string         `json:"timezone"`            // 当前服务器所在的时区
}

func GetSystemInfo() (res schema.Response) {
	var (
		err      error
		data     Info
		hostInfo *host.InfoStat
		CPUInfo  []cpu.InfoStat
		avgStat  *load.AvgStat
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

		helper.Response(&res, data, nil, err)
	}()

	v, _ := mem.VirtualMemory()

	if CPUInfo, err = cpu.Info(); err != nil {
		return
	}

	if hostInfo, err = host.Info(); err != nil {
		return
	}

	if avgStat, err = load.Avg(); err != nil {
		return
	}

	var u *user.User

	if u, err = user.Current(); err != nil {
		return
	}

	t := time.Now()

	data = Info{
		Username:         u.Username,
		Host:             *hostInfo,
		Arch:             runtime.GOARCH,
		Avg:              *avgStat,
		CPU:              CPUInfo,
		RAMAvailable:     v.Available,
		RAMTotal:         v.Total,
		RAMFree:          v.Free,
		RAMUsedBy:        v.Used,
		RAMUsedByPercent: v.UsedPercent,
		Time:             t.Format(time.RFC3339Nano),
		Timezone:         t.Location().String(),
	}

	return
}

var GetSystemInfoRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return GetSystemInfo()
	})
})
