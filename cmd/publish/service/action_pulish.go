package service

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"

	"github.com/edufriendchen/light-tiktok/cmd/publish/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/publish"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	vieoSuffix  = ".mp4"
	imageSuffix = ".png"
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
	rawUniqueNum, err := global.ID_GENERATOR.NextID()
	if err != nil {
		return err
	}
	uniqueNum := strconv.FormatUint(rawUniqueNum, 10)
	video_key, err := dal.MinioUpload(bytes.NewReader(videoData), uniqueNum+vieoSuffix, int64(len(videoData)))
	if err != nil {
		return err
	}
	prefix := global.Config.GetString(consts.MINIO_PREFIX)
	iamgeData, err := getFrame(prefix, uniqueNum+imageSuffix)
	if err != nil {
		fmt.Println("Failed to extract frame:", err)
		return err
	}
	cover_key, err := dal.MinioUpload(bytes.NewReader(iamgeData), uniqueNum+imageSuffix, int64(len(iamgeData)))
	if err != nil {
		return err
	}
	_, err = dal.CreateVideo(s.ctx, s.session, user_id, prefix+video_key, prefix+cover_key, req.Title)
	if err != nil {
		return err
	}
	return nil
}

// Use ffmpeg to extract specified frames as image files
func getFrame(videoUri string, outputPath string) ([]byte, error) {
	cmd := exec.Command("ffmpeg", "-i", videoUri, "-vf", "select=eq(n\\,100)", "-vframes", "1", outputPath)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	content, err := ioutil.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}
	return content, nil
}
