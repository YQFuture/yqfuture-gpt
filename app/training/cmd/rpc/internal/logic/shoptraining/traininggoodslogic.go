package shoptraininglogic

import (
	"context"
	"github.com/tidwall/gjson"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/internal/thirdparty"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type TrainingGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTrainingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrainingGoodsLogic {
	return &TrainingGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// TrainingGoods 训练商品
func (l *TrainingGoodsLogic) TrainingGoods(in *training.TrainingGoodsReq) (*training.TrainingGoodsResp, error) {
	// 根据shopId从mongo中找到最新的一条预设商品数据
	dbpresettinggoodstitles, err := l.svcCtx.DbpresettinggoodstitlesModel.FindNewOneByGoodsId(l.ctx, in.GoodsId)
	if dbpresettinggoodstitles == nil || len(dbpresettinggoodstitles.GoodsDocumentList) == 0 {
		l.Logger.Error("mongo中没有可用的预设数据", err)
		return nil, err
	}
	// 预设保存的商品文档列表
	var goodsDocumentList = dbpresettinggoodstitles.GoodsDocumentList
	// 从mysql中查询出商品
	tsGoods, err := l.svcCtx.TsGoodsModel.FindOne(l.ctx, in.GoodsId)
	if err != nil {
		l.Logger.Error("查找商品数据失败", err)
		return nil, err
	}
	if tsGoods.TrainingStatus == 1 {
		return &training.TrainingGoodsResp{
			Result: "商品已经在训练中",
		}, nil
	}
	// 训练的商品列表
	var trainingGoodsList []*orm.TsGoods
	trainingGoodsList = append(trainingGoodsList, tsGoods)
	// 更新商品状态为训练中
	err = UpdateGoodsTraining(l.ctx, l.svcCtx, tsGoods, in.UserId)
	if err != nil {
		l.Logger.Error("修改商品状态失败", err)
	}

	// 发起店铺训练批处理 获取返回的batchId
	var createBatchTaskResp string
	batchId, err := thirdparty.CreateBatchTask(l.Logger, l.svcCtx, goodsDocumentList, &createBatchTaskResp)
	if err != nil {
		l.Logger.Error("发送创建店铺训练批处理请求失败", err)
		return nil, err
	}
	// 等待2分钟
	time.Sleep(time.Minute * 2)
	// 轮询等待批处理完成 获取返回的fileId
	fileId, err := thirdparty.GetBatchTaskStatus(l.Logger, l.svcCtx, batchId)
	if err != nil {
		return nil, err
	}
	// 获取批处理结果 对于识别失败的结果将不返回
	var batchTaskResultResp string
	err = utils.HTTPGetAndParseJSON(l.svcCtx.Config.TrainingGoodsConf.QueryBatchTaskResultUrl+"?file_id="+fileId, &batchTaskResultResp)
	if err != nil {
		l.Logger.Error("发送获取批处理结果请求失败", err)
		return nil, err
	}
	// 解析结果写入goodsDocument
	var batchTaskResultMap map[string]*gjson.Result
	for _, batchTaskResult := range gjson.Get(batchTaskResultResp, "data").Array() {
		batchTaskResultMap[batchTaskResult.Get("custom_id").String()] = &batchTaskResult
	}
	for _, goodsDocument := range goodsDocumentList {
		// 只有训练成功的商品才去获取训练结果
		if batchTaskResultMap[goodsDocument.PlatformGoodsId] != nil {
			// 保存训练结果和消耗的token
			goodsDocument.DetailGalleryDescription = batchTaskResultMap[goodsDocument.PlatformGoodsId].Get("content").String()
			goodsDocument.Token = batchTaskResultMap[goodsDocument.PlatformGoodsId].Get("token").Int()
		}
		// 保存训练结果到ES
		es := l.svcCtx.Elasticsearch
		res, err := es.Index().Index("training_goods").BodyJson(goodsDocument).Refresh("true").Do(context.Background())
		if err != nil {
			logx.Errorf("商品解析结果写入ES失败, err :%s", err.Error())
			continue
		}
		logx.Infof("商品解析结果写入ES成功, res :%v", res)
	}

	// 修改商品状态
	for _, trainingGoods := range trainingGoodsList {
		err = UpdateGoodsTrainingComplete(l.ctx, l.svcCtx, trainingGoods, in.UserId)
		if err != nil {
			l.Logger.Error("修改商品状态失败", trainingGoods)
			return nil, err
		}
	}
	// 返回正常
	return &training.TrainingGoodsResp{}, nil
}
