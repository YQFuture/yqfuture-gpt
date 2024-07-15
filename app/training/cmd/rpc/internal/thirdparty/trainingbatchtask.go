package thirdparty

import (
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/model/common"
	"yufuture-gpt/common/utils"
)

type ImageInfo struct {
	ID   string   `json:"id"`
	URLs []string `json:"urls"`
}

type CreateBatchTaskReq struct {
	SystemPrompt string       `json:"system_prompt"`
	BatchImages  []*ImageInfo `json:"batch_image_urls"`
}

func CreateBatchTask(log logx.Logger, svcCtx *svc.ServiceContext, goodsDocumentList []*common.PddGoodsDocument, createBatchTaskResp *string) (string, error) {
	var batchImages []*ImageInfo
	for _, goodsDocument := range goodsDocumentList {
		batchImages = append(batchImages, &ImageInfo{
			ID:   goodsDocument.PlatformGoodsId,
			URLs: goodsDocument.PictureUrlList,
		})
	}
	err := utils.HTTPPostAndParseJSON(svcCtx.Config.TrainingGoodsConf.CreateBatchTaskUrl, CreateBatchTaskReq{
		SystemPrompt: "what do you see ？ reply in Chinese",
		BatchImages:  batchImages,
	}, createBatchTaskResp)
	if err != nil {
		log.Error("发送创建店铺训练批处理请求失败", err)
		return "", err
	}
	batchId := gjson.Get(*createBatchTaskResp, "data.response.batch_info.id")
	return batchId.String(), nil
}

func GetBatchTaskStatus(log logx.Logger, svcCtx *svc.ServiceContext, batchId string) (string, error) {
	var fileId string
	for {
		var batchTaskStatusResp string
		err := utils.HTTPGetAndParseJSON(svcCtx.Config.TrainingGoodsConf.QueryBatchTaskStatusUrl+"?batch_id="+batchId, &batchTaskStatusResp)
		if err != nil {
			log.Error("发送获取批处理状态请求失败", err)
			return "nil", err
		}
		status := gjson.Get(batchTaskStatusResp, "data.status")
		if status.String() == "completed" {
			fileId = gjson.Get(batchTaskStatusResp, "data.output_file_id").String()
			break
		}
		time.Sleep(time.Minute * 2)
	}
	return fileId, nil
}
