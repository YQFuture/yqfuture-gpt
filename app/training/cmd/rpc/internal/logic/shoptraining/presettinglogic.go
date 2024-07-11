package shoptraininglogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/app/training/model/common"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/consts"
)

type PreSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingLogic {
	return &PreSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreSettingLogic) PreSetting(in *training.ShopTrainingReq) (*training.ShopTrainingResp, error) {
	// TODO 发送爬取商品ID请求

	// TODO 轮询等待爬取商品ID落库完成

	// TODO 将店铺和商品置为预训练状态

	// 根据uuid和userId从mongo中找到最新的一条店铺数据
	saveShop, err := l.svcCtx.ShoptrainingshoptitlesModel.FindNewOneByUuidAndUserId(l.ctx, in.Uuid, in.UserId)
	if err != nil {
		l.Logger.Error("根据uuid和userId从mongo中找到最新的一条店铺数据失败", err)
		return nil, err
	}

	// 根据uuid和userid从mysql中查找出店铺
	tsShop, err := l.svcCtx.TsShopModel.FindOneByUuidAndUserId(l.ctx, in.UserId, in.Uuid)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找店铺失败", err)
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
			tsGoodsMap[tsGoods.PlatformId] = tsGoods
		}
	}

	// 更新店铺状态为预训练中
	err = UpdateShopPreSetting(l.ctx, l.svcCtx, tsShop, in.UserId)
	if err != nil {
		l.Logger.Error("修改店铺状态失败", in)
		return nil, err
	}
	// 更新商品状态为预训练中 同时提取本次需要预训练的商品列表
	var preSettingGoodsList []*orm.TsGoods
	for _, saveGoods := range saveShop.GoodsList {
		// 同时在mongo中并且enabled字段为2启用的商品即为本次需要预训练的商品
		if tsGoods, ok := tsGoodsMap[saveGoods.PlatformId]; ok {
			// 排除掉状态已经在预训练中/训练中/训练完成的商品
			if tsGoods.TrainingStatus == consts.Presetting || tsGoods.TrainingStatus == consts.Training || tsGoods.TrainingStatus == consts.PresettingComplete {
				continue
			}
			// 更新商品状态为预训练中
			err = UpdateGoodsPreSetting(l.ctx, l.svcCtx, tsGoods, in.UserId)
			if err != nil {
				l.Logger.Error("修改商品状态失败", tsGoods)
				return nil, err
			}
			//将筛选出的商品添加到预训练商品列表
			preSettingGoodsList = append(preSettingGoodsList, tsGoods)
		}
	}

	// 请求获取商品JSON
	err = ApplyGoodsJson(l.svcCtx, preSettingGoodsList)
	if err != nil {
		l.Logger.Info("发送获取商品JSON请求失败", err)
		return nil, err
	}

	// 等待2分钟
	time.Sleep(time.Minute * 2)

	// 每6分钟调用一次接口 连续10次失败则结束
	FetchAndSaveGoodsJson(l.Logger, l.ctx, l.svcCtx, preSettingGoodsList)

	// 从jSON解析的商品列表文档
	var goodsDocumentList []*common.PddGoodsDocument
	// 获取并解析商品JSON到结果文档列表
	GetAndParseGoodsJson(l.Logger, tsShop, goodsDocumentList, preSettingGoodsList)

	// 构建获取训练时长的请求图片列表
	var goodPicList []string
	for _, goodsDocument := range goodsDocumentList {
		goodPicList = append(goodPicList, goodsDocument.PictureUrlList...)
	}

	// TODO 发送请求 获取商品训练所需时长

	// 设计结构化文档 预训练结果保存到mongo 正式训练时直接从mongo中取

	// 更新店铺和商品状态为预训练完成

	return &training.ShopTrainingResp{}, nil
}

func UpdateShopPreSetting(ctx context.Context, svcCtx *svc.ServiceContext, tsShop *orm.TsShop, userId int64) error {
	tsShop.TrainingStatus = consts.Presetting
	tsShop.UpdateTime = time.Now()
	tsShop.UpdateBy = userId
	err := svcCtx.TsShopModel.Update(ctx, tsShop)
	if err != nil {
		return err
	}
	return nil
}

func UpdateShopPreSettingComplete(ctx context.Context, svcCtx *svc.ServiceContext, tsShop *orm.TsShop, userId int64) error {
	tsShop.TrainingStatus = consts.PresettingComplete
	tsShop.UpdateTime = time.Now()
	tsShop.UpdateBy = userId
	err := svcCtx.TsShopModel.Update(ctx, tsShop)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGoodsPreSetting(ctx context.Context, svcCtx *svc.ServiceContext, tsGoods *orm.TsGoods, userId int64) error {
	tsGoods.TrainingStatus = consts.Presetting
	tsGoods.UpdateTime = time.Now()
	tsGoods.UpdateBy = userId
	err := svcCtx.TsGoodsModel.Update(ctx, tsGoods)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGoodsPreSettingComplete(ctx context.Context, svcCtx *svc.ServiceContext, tsGoods *orm.TsGoods, userId int64) error {
	tsGoods.TrainingStatus = consts.PresettingComplete
	tsGoods.UpdateTime = time.Now()
	tsGoods.UpdateBy = userId
	err := svcCtx.TsGoodsModel.Update(ctx, tsGoods)
	if err != nil {
		return err
	}
	return nil
}
