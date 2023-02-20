package minio

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/minio/minio-go/v7"
)

// CreateBucket bucket
func CreateBucket(bucketName string) error {
	if len(bucketName) <= 0 {
		klog.Error("bucketName invalid")
	}

	location := "beijing"
	ctx := context.Background()

	err := global.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := global.MinioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			klog.Debugf("bucket %s already exists", bucketName)
			return nil
		} else {
			return err
		}
	} else {
		klog.Infof("bucket %s create successfully", bucketName)
	}
	return nil
}

// UploadLocalFile 上传本地文件（提供文件路径）至 minio
func UploadLocalFile(bucketName string, objectName string, filePath string, contentType string) (int64, error) {
	ctx := context.Background()
	info, err := global.MinioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		klog.Errorf("localfile upload failed, %s", err)
		return 0, err
	}
	klog.Infof("upload %s of size %d successfully", objectName, info.Size)
	return info.Size, nil
}

// UploadFile 上传文件（提供reader）至 minio
func UploadFile(bucketName string, objectName string, reader io.Reader, objectsize int64) error {
	ctx := context.Background()
	n, err := global.MinioClient.PutObject(ctx, bucketName, objectName, reader, objectsize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		klog.Errorf("upload %s of size %d failed, %s", bucketName, objectsize, err)
		return err
	}
	klog.Infof("upload %s of bytes %d successfully", objectName, n.Size)
	return nil
}

// GetFileUrl 从 minio 获取文件Url
func GetFileUrl(bucketName string, fileName string, expires time.Duration) (*url.URL, error) {
	ctx := context.Background()
	reqParams := make(url.Values)
	if expires <= 0 {
		expires = time.Second * 60 * 60 * 24
	}
	presignedUrl, err := global.MinioClient.PresignedGetObject(ctx, bucketName, fileName, expires, reqParams)
	if err != nil {
		klog.Errorf("get url of file %s from bucket %s failed, %s", fileName, bucketName, err)
		return nil, err
	}
	// TODO: url可能要做截取
	return presignedUrl, nil
}
