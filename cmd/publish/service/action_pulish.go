package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"strconv"

	"github.com/edufriendchen/light-tiktok/cmd/publish/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/publish"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type ActionPulishService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
}

func NewActionPulishService(ctx context.Context, driver neo4j.DriverWithContext) *ActionPulishService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	return &ActionPulishService{ctx: ctx, session: session}
}

func (s *ActionPulishService) ActionPulish(req *publish.ActionRequest, user_id int64) error {
	videoData := []byte(req.Data)
	video_id, err := global.ID_GENERATOR.NextID()
	if err != nil {
		return err
	}
	name := strconv.FormatUint(video_id, 10)
	fileName := name + "." + "mp4"
	coverName := name + "." + "jpg"
	reader := bytes.NewReader(videoData)
	dataLen := int64(len(videoData))
	video_key, err := Upload(reader, fileName, dataLen)
	if err != nil {
		return err
	}
	localFile := "/home/test.jpg"
	cover_key, err := UploadCover(localFile, coverName, dataLen)
	if err != nil {
		return err
	}
	_, err = dal.CreateVideo(s.ctx, s.session, user_id, "http://rpstobjks.hb-bkt.clouddn.com/"+video_key, "http://rpstobjks.hb-bkt.clouddn.com/"+cover_key, req.Title)
	if err != nil {
		return err
	}
	return nil
}

// 从视频流中截取一帧并返回 需要在本地环境中安装ffmpeg并将bin添加到环境变量
func captureCover(filePath string) ([]byte, error) {
	reader := bytes.NewBuffer(nil)
	err := ffmpeg.Input(filePath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(reader, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	return buf.Bytes(), err
}

var videoFileExt = []string{"mp4", "flv"} //此处可根据需要添加格式

func IsVideoAllowed(suffix string) bool {
	for _, fileExt := range videoFileExt {
		if suffix == fileExt {
			return true
		}
	}
	return false
}

func Upload(file io.Reader, filename string, size int64) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: global.Config.GetString(consts.MINIO_BUCKET_NAME),
	}
	mac := qbox.NewMac(global.Config.GetString(consts.MINIO_ACCESS_KEY), global.Config.GetString(consts.MINIO_SECRET_ACCESS_KEY))
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	err := formUploader.Put(context.Background(), &ret, upToken, filename, file, size, &putExtra)
	if err != nil {
		return "", err
	}
	return ret.Key, nil
}

func UploadCover(localFile string, filename string, size int64) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: global.Config.GetString(consts.MINIO_BUCKET_NAME),
	}
	mac := qbox.NewMac(global.Config.GetString(consts.MINIO_ACCESS_KEY), global.Config.GetString(consts.MINIO_SECRET_ACCESS_KEY))
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	err := formUploader.PutFile(context.Background(), &ret, upToken, filename, localFile, &putExtra)
	if err != nil {
		return "", err
	}
	return ret.Key, nil
}
