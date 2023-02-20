package initialize

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	util "github.com/edufriendchen/light-tiktok/pkg/minio"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Minio 对象存储初始化
func InitMinio() {
	client, err := minio.New(consts.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(consts.MinioAccessKeyId, consts.MinioSecretAccessKey, ""),
		Secure: consts.MinioUseSSL,
	})
	if err != nil {
		klog.Errorf("minio client init failed: %v", err)
	}
	// fmt.Println(client)
	klog.Debug("minio client init successfully")
	global.MinioClient = client
	if err := util.CreateBucket(consts.MinioVideoBucketName); err != nil {
		klog.Errorf("minio client init failed: %v", err)
	}
}
