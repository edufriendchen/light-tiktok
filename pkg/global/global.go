package global

import (
	configtype "github.com/edufriendchen/light-tiktok/pkg/config"
	"github.com/edufriendchen/light-tiktok/pkg/jwt"
	"github.com/minio/minio-go/v7"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sony/sonyflake"
	"gorm.io/gorm"
)

var (
	Config       configtype.Config
	NacosClient  naming_client.INamingClient
	Neo4jDriver  neo4j.DriverWithContext
	Jwt          *jwt.JWT
	DB           *gorm.DB
	MinioClient  *minio.Client
	ID_GENERATOR *sonyflake.Sonyflake
)
