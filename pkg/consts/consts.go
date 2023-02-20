package consts

const (
	ApiServiceName      = "api"
	ApiServiceAddr      = ":1060"
	UserServiceName     = "user"
	UserServiceAddr     = ":9000"
	RelationServiceName = "relation"
	RelationServiceAddr = ":9100"
	MessageServiceName  = "message"
	MessageServiceAddr  = ":9200"
	CommentServiceName  = "comment"
	CommentServiceAddr  = ":9300"
	FavoriteServiceName = "favorite"
	FavoriteServiceAddr = ":9400"
	FeedServiceName     = "feed"
	FeedServiceddr      = ":9500"
	PushlishServiceName = "puslish"
	PushlishServiceddr  = ":9600"
)

const (
	NoteTableName   = "note"
	UserTableName   = "user"
	SecretKey       = "secret key"
	IdentityKey     = "id"
	Total           = "total"
	Notes           = "notes"
	MySQLDefaultDSN = "root:Zty20001011!@tcp(localhost:3306)/tiktok?charset=utf8&parseTime=True&loc=Local"
	TCP             = "tcp"
	ExportEndpoint  = "localhost:14268"
	DefaultLimit    = 10
	Neo4jUri        = "neo4j://localhost:7687"
	Neo4jUsername   = "neo4j"
	Neo4jPassword   = "friendchen0429"
	NacosAddr       = "127.0.0.1"
	NacosPoint      = 8848
	NacosLogDir     = ""
	JWTSecretKey    = "chen"
	Limit           = 35
)

const (
	MinioEndpoint        = "http://portal.qiniu.com/cdn/domain/rpstobjks.hb-bkt.clouddn.com"
	MinioAccessKeyId     = "GWh-FLBRPXOY59v_bKHMQ3At7fX7sZgIs5kMcrOB"
	MinioSecretAccessKey = "pkgIULQ4BRdhuocutEzIhPDkTJfL5TLVFAgjBYIf"
	MinioUseSSL          = true
	MinioVideoBucketName = "friendchen"
	START_TIME           = "2022-05-21 00:00:01" // 固定启动时间，保证生成 ID 唯一性START_TIME           = "2022-05-21 00:00:01"  // 固定启动时间，保证生成 ID 唯一性
)
