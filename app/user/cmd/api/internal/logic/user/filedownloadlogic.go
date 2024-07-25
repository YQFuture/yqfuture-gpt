package user

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"net/http"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFileDownloadLogic 文件下载
func NewFileDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDownloadLogic {
	return &FileDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadReq, w http.ResponseWriter) error {
	path := req.Path
	fileName := req.FileName

	bucketName := l.svcCtx.Config.OssConf.BucketName

	client, err := oss.New(l.svcCtx.Config.OssConf.Endpoint, l.svcCtx.Config.OssConf.AccessKeyId, l.svcCtx.Config.OssConf.AccessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 下载文件
	objectKey := path + "/" + fileName
	reader, err := bucket.GetObject(objectKey)
	if err != nil {
		l.Logger.Error("从阿里云OSS下载文件失败", err)
		return err
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			l.Logger.Error(err)
		}
	}(reader)

	// 设置HTTP响应头
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	// 将文件数据写入HTTP响应
	_, err = io.Copy(w, reader)
	if err != nil {
		l.Logger.Error("将从阿里云OSS下载的文件写入HTTP响应失败", err)
		return err
	}
	return nil
}
