package initialize

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/edufriendchen/light-tiktok/pkg/global"

	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitNeo4j() {
	driver, err := neo4j.NewDriverWithContext(consts.Neo4jUri, neo4j.BasicAuth("neo4j", "friendchen0429", ""))
	global.Neo4jDriver = driver
	if err != nil {
		klog.Infof("Neo4j Init Error: %v", err)
	}
}
