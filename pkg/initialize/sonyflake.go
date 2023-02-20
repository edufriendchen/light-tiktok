package initialize

import (
	"math/rand"
	"time"

	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/sony/sonyflake"
)

// InitDB to init Sonyflake
func InitSonyflake() {
	rand.Seed(time.Now().Unix())
	startTime, _ := time.Parse("2006-01-02 15:04:05", consts.START_TIME)
	global.ID_GENERATOR = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: startTime,
	})
}
