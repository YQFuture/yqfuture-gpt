package shoptraininglogic

import (
	"context"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/thirdparty"
	"yufuture-gpt/app/training/model/common"
	yqmongo "yufuture-gpt/app/training/model/mongo"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/consts"
	"yufuture-gpt/common/utils"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSettingGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreSettingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingGoodsLogic {
	return &PreSettingGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 预设商品
func (l *PreSettingGoodsLogic) PreSettingGoods(in *training.PreSettingGoodsReq) (*training.PreSettingGoodsResp, error) {
	tsGoods, err := l.svcCtx.TsGoodsModel.FindOne(l.ctx, in.GoodsId)
	if err != nil {
		l.Logger.Error("根据goodsId查找商品失败", err)
		return nil, err
	}
	// 排除掉状态已经在预设中/训练中/预设完成的商品
	if tsGoods.TrainingStatus == consts.Presetting || tsGoods.TrainingStatus == consts.Training || tsGoods.TrainingStatus == consts.PresettingComplete {
		l.Logger.Error("商品不是可预设状态", err)
		return nil, err
	}
	tsShop, err := l.svcCtx.TsShopModel.FindOne(l.ctx, tsGoods.ShopId)
	if err != nil {
		l.Logger.Error("根据ShopId查找店铺失败", err)
		return nil, err
	}

	// 更新商品状态为预设中
	err = UpdateGoodsPreSetting(l.ctx, l.svcCtx, tsGoods, in.UserId)

	var presettingGoods []*orm.TsGoods
	presettingGoods = append(presettingGoods, tsGoods)
	// 请求获取商品JSON
	err = thirdparty.ApplyGoodsJson(l.svcCtx, presettingGoods)
	if err != nil {
		l.Logger.Info("发送获取商品JSON请求失败", err)
		return nil, err
	}
	// 等待2分钟
	time.Sleep(time.Minute * 2)
	// 每6分钟调用一次接口 连续10次失败则结束
	thirdparty.FetchAndSaveGoodsJson(l.Logger, l.ctx, l.svcCtx, presettingGoods)
	// 最终保存到ES的结果文档
	var goodsDocumentList []*common.PddGoodsDocument
	// 获取并解析商品JSON到结果文档列表
	thirdparty.GetAndParseGoodsJson(l.Logger, tsShop, goodsDocumentList, presettingGoods)

	// 构建获取训练时长的请求图片列表
	var goodPicList []string
	for _, goodsDocument := range goodsDocumentList {
		goodPicList = append(goodPicList, goodsDocument.PictureUrlList...)
	}
	//发送请求 获取商品训练所需资源和时长
	var fetchEstimateResultResp FetchEstimateResultResp
	err = utils.HTTPPostAndParseJSON(l.svcCtx.Config.TrainingGoodsConf.FetchEstimateResultUrl, struct {
		Urls []string `json:"urls"`
	}{Urls: goodPicList}, &fetchEstimateResultResp)
	if err != nil {
		l.Logger.Error("获取商品训练所需资源和时长失败", err)
	}

	// 设计结构化文档 预设结果保存到mongo 正式训练时直接从mongo中取
	shoppresettinggoodstitles := &yqmongo.Shoppresettinggoodstitles{
		GoodsId:    tsGoods.Id,
		PlatformId: goodsDocumentList[0].PlatformMallId,
		UserID:     in.UserId,

		PreSettingToken:    fetchEstimateResultResp.Data.Token,
		PresettingPower:    fetchEstimateResultResp.Data.Power,
		PresettingFileSize: fetchEstimateResultResp.Data.FileSize,
		//PreSettingTime:
		GoodsDocumentList: goodsDocumentList,
	}
	err = l.svcCtx.ShoppresettinggoodstitlesModel.Insert(l.ctx, shoppresettinggoodstitles)
	// 更新商品状态为预设完成
	err = UpdateGoodsPreSettingComplete(l.ctx, l.svcCtx, tsGoods, in.UserId)
	if err != nil {
		l.Logger.Error("修改商品状态失败", err)
	}

	return &training.PreSettingGoodsResp{}, nil
}
