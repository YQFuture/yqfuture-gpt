package user

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"
	"yufuture-gpt/common/consts"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFileUploadLogic 文件上传
func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(r *http.Request) (resp *types.FileUploadResp, err error) {
	path := r.FormValue("path")
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		l.Logger.Error("读取上传文件失败 ", err)
		return &types.FileUploadResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "上传失败",
			},
		}, nil
	}

	bucketName := l.svcCtx.Config.OssConf.BucketName
	client, err := oss.New(l.svcCtx.Config.OssConf.Endpoint, l.svcCtx.Config.OssConf.AccessKeyId, l.svcCtx.Config.OssConf.AccessKeySecret)
	if err != nil {
		l.Logger.Error("获取访问凭证失败 ", err)
		return &types.FileUploadResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "上传失败",
			},
		}, nil
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		l.Logger.Error("获取存储空间失败 ", err)
		return &types.FileUploadResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "上传失败",
			},
		}, nil
	}

	fileUuid := uuid.New().String()
	fileName := fileUuid + "_" + fileHeader.Filename
	objectKey := path + "/" + fileName

	err = bucket.PutObject(objectKey, file)
	if err != nil {
		l.Logger.Error("上传文件到阿里云OSS失败 ", err)
		return &types.FileUploadResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "上传失败",
			},
		}, nil
	}

	return &types.FileUploadResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "上传成功",
		},
		Data: fileName,
	}, nil
}
