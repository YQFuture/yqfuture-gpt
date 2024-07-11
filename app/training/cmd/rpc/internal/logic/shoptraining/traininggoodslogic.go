package shoptraininglogic

import (
	"context"
	"github.com/tidwall/gjson"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/app/training/model/common"
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

// 训练商品
func (l *TrainingGoodsLogic) TrainingGoods(in *training.TrainingGoodsReq) (*training.TrainingGoodsResp, error) {
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
	tsShop, err := l.svcCtx.TsShopModel.FindOne(l.ctx, tsGoods.ShopId)
	if err != nil {
		l.Logger.Error("根据ShopId查找店铺失败", err)
		return nil, err
	}

	var trainingGoodsList []*orm.TsGoods
	trainingGoodsList = append(trainingGoodsList, tsGoods)

	// 请求获取商品JSON
	err = ApplyGoodsJson(l.svcCtx, trainingGoodsList)
	if err != nil {
		l.Logger.Info("发送获取商品JSON请求失败", err)
		return nil, err
	}

	// 等待2分钟
	time.Sleep(time.Minute * 2)

	// 每6分钟调用一次接口 连续10次失败则结束
	FetchAndSaveGoodsJson(l.Logger, l.ctx, l.svcCtx, trainingGoodsList)

	// 最终保存到ES的结果文档
	var goodsDocumentList []*common.PddGoodsDocument
	// 获取并解析商品JSON到结果文档列表
	GetAndParseGoodsJson(l.Logger, tsShop, goodsDocumentList, trainingGoodsList)

	// 发起店铺训练批处理 获取返回的batchId
	var createBatchTaskResp string
	batchId, err := CreateBatchTask(l.Logger, l.svcCtx, goodsDocumentList, &createBatchTaskResp)
	if err != nil {
		l.Logger.Error("发送创建店铺训练批处理请求失败", err)
		return nil, err
	}

	// 等待2分钟
	time.Sleep(time.Minute * 2)

	// 轮询等待批处理完成 获取返回的fileId
	fileId, err := GetBatchTaskStatus(l.Logger, l.svcCtx, batchId)
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
