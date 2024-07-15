package shoptraininglogic

import (
	"context"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/internal/thirdparty"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/consts"
	"yufuture-gpt/common/utils"
)

type TrainingShopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type ImageInfo struct {
	ID   string   `json:"id"`
	URLs []string `json:"urls"`
}

type CreateBatchTaskReq struct {
	SystemPrompt string       `json:"system_prompt"`
	BatchImages  []*ImageInfo `json:"batch_image_urls"`
}

func NewTrainingShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrainingShopLogic {
	return &TrainingShopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// TrainingShop 训练店铺
func (l *TrainingShopLogic) TrainingShop(in *training.TrainingShopReq) (*training.TrainingShopResp, error) {
	// 根据uuid和userId从mongo中找到最新的一条预设店铺数据
	shoppresettingshoptitles, err := l.svcCtx.ShoppresettingshoptitlesModel.FindNewOneByUuidAndUserId(l.ctx, in.Uuid, in.UserId)
	if shoppresettingshoptitles == nil || len(shoppresettingshoptitles.GoodsDocumentList) == 0 {
		l.Logger.Error("mongo中没有可用的预设数据", err)
		return nil, err
	}
	if err != nil {
		l.Logger.Error("根据uuid和userId从mongo中找到最新的一条预设店铺数据失败", err)
		return nil, err
	}
	// 预设保存的商品文档列表
	var goodsDocumentList = shoppresettingshoptitles.GoodsDocumentList
	// 根据uuid和userid从mysql中查找出店铺
	tsShop, err := l.svcCtx.TsShopModel.FindOneByUuidAndUserId(l.ctx, in.UserId, in.Uuid)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找店铺失败", err)
		return nil, err
	}
	if tsShop.TrainingStatus != consts.PresettingComplete {
		l.Logger.Error("只有预设完成的店铺才能进行训练", err)
		return nil, err
	}
	// 根据店铺shopId从mysql中查找出enabled字段为2启用商品列表
	tsGoodsList, err := l.svcCtx.TsGoodsModel.FindEnabledListByShopId(l.ctx, tsShop.Id)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找商品失败", err)
		return nil, err
	}
	// 商品列表转换为map 便于后续查找
	var tsGoodsMap map[string]*orm.TsGoods
	if tsGoodsList != nil {
		tsGoodsMap = make(map[string]*orm.TsGoods)
		for _, tsGoods := range *tsGoodsList {
			// 只保留预设完成状态的商品
			if tsGoods.TrainingStatus != consts.TrainingComplete {
				continue
			}
			tsGoodsMap[tsGoods.PlatformId] = tsGoods
		}
	}

	// 更新店铺状态为训练中 添加训练次数
	err = UpdateShopTraining(l.ctx, l.svcCtx, tsShop, in.UserId)
	if err != nil {
		l.Logger.Error("修改店铺状态失败", in)
		return nil, err
	}
	// 更新商品状态为训练中 同时提取本次需要训练的商品列表
	var trainingGoodsList []*orm.TsGoods
	for _, goodsDocument := range goodsDocumentList {
		// 同时在mongo中并且enabled字段为2启用的商品即为本次需要训练的商品
		if tsGoods, ok := tsGoodsMap[goodsDocument.PlatformGoodsId]; ok {
			// 排除掉并非预设完成状态的商品
			if tsGoods.TrainingStatus != consts.TrainingComplete {
				continue
			}
			// 更新商品状态为训练中 添加训练次数
			err = UpdateGoodsTraining(l.ctx, l.svcCtx, tsGoods, in.UserId)
			if err != nil {
				l.Logger.Error("修改商品状态失败", tsGoods)
				return nil, err
			}
			//将筛选出的商品添加到训练商品列表
			trainingGoodsList = append(trainingGoodsList, tsGoods)
		}
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
		l.Logger.Error("获取批处理状态失败", err)
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
			goodsDocument.Power = batchTaskResultMap[goodsDocument.PlatformGoodsId].Get("power").Int()
			goodsDocument.FileSize = batchTaskResultMap[goodsDocument.PlatformGoodsId].Get("filesize").Int()
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

	// 更新数据库状态为训练完成 同时保存训练历史
	err = UpdateShopTrainingComplete(l.ctx, l.svcCtx, tsShop, in.UserId)
	if err != nil {
		l.Logger.Error("修改店铺状态失败", in)
		return nil, err
	}
	for _, trainingGoods := range trainingGoodsList {
		err = UpdateGoodsTrainingComplete(l.ctx, l.svcCtx, trainingGoods, in.UserId)
		if err != nil {
			l.Logger.Error("修改商品状态失败", trainingGoods)
			return nil, err
		}
	}
	// 返回正常
	return &training.TrainingShopResp{}, nil
}

func UpdateShopTraining(ctx context.Context, svcCtx *svc.ServiceContext, tsShop *orm.TsShop, userId int64) error {
	tsShop.TrainingStatus = consts.Training
	tsShop.TrainingTimes += 1
	tsShop.UpdateTime = time.Now()
	tsShop.UpdateBy = userId
	err := svcCtx.TsShopModel.Update(ctx, tsShop)
	if err != nil {
		return err
	}
	return nil
}

func UpdateShopTrainingComplete(ctx context.Context, svcCtx *svc.ServiceContext, tsShop *orm.TsShop, userId int64) error {
	tsShop.TrainingStatus = consts.TrainingComplete
	tsShop.UpdateTime = time.Now()
	tsShop.UpdateBy = userId
	err := svcCtx.TsShopModel.Update(ctx, tsShop)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGoodsTraining(ctx context.Context, svcCtx *svc.ServiceContext, tsGoods *orm.TsGoods, userId int64) error {
	tsGoods.TrainingStatus = consts.Training
	tsGoods.TrainingTimes += 1
	tsGoods.UpdateTime = time.Now()
	tsGoods.UpdateBy = userId
	tsGoods.GoodsJsonUrl = "" //每次训练开始 把获取商品json的url字段置空
	err := svcCtx.TsGoodsModel.Update(ctx, tsGoods)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGoodsTrainingComplete(ctx context.Context, svcCtx *svc.ServiceContext, tsGoods *orm.TsGoods, userId int64) error {
	tsGoods.TrainingStatus = consts.TrainingComplete
	tsGoods.UpdateTime = time.Now()
	tsGoods.UpdateBy = userId
	err := svcCtx.TsGoodsModel.Update(ctx, tsGoods)
	if err != nil {
		return err
	}
	return nil
}
