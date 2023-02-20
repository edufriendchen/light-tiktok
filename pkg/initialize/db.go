package initialize

import (
	"time"

	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	//"gorm.io/plugin/opentelemetry/logging/logrus"
)

// InitDB to init DB
func InitDB() {
	var err error
	gormlogrus := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond,
			Colorful:      false,
			LogLevel:      logger.Info,
		},
	)

	global.DB, err = gorm.Open(mysql.Open(global.Config.GetString(consts.MYSQL_DEFAULT_DSN)),
		&gorm.Config{
			PrepareStmt: true,
			Logger:      gormlogrus,
		},
	)
	if err != nil {
		panic(err)
	}
	// if err := global.DB.Use(tracing.NewPlugin()); err != nil {
	// 	panic(err)
	// }
}
