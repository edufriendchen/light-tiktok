package global

import (
	"github.com/edufriendchen/light-tiktok/pkg/jwt"
	"github.com/minio/minio-go/v7"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sony/sonyflake"
	"gorm.io/gorm"
)

var (
	Neo4jDriver  neo4j.DriverWithContext
	Jwt          *jwt.JWT
	DB           *gorm.DB
	MinioClient  *minio.Client
	ID_GENERATOR *sonyflake.Sonyflake
)
