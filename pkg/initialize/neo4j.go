package initialize

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/edufriendchen/light-tiktok/pkg/global"

	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitNeo4j() {
	driver, err := neo4j.NewDriverWithContext(global.Config.GetString(consts.DB_NEO4J_URI), neo4j.BasicAuth(global.Config.GetString(consts.DB_NEO4J_USERNAME), global.Config.GetString(consts.DB_NEO4J_PASS), ""))
	global.Neo4jDriver = driver
	if err != nil {
		klog.Infof("Neo4j Init Error: %v", err)
	}
}
