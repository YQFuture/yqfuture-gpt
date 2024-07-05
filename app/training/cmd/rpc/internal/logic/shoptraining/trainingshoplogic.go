package shoptraininglogic

import (
	"context"
	"time"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/utills"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	yqmongo "yufuture-gpt/app/training/model/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type TrainingShopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTrainingShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrainingShopLogic {
	return &TrainingShopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 训练店铺
func (l *TrainingShopLogic) TrainingShop(in *training.TrainingShopReq) (*training.TrainingShopResp, error) {
	//TODO 根据uuid和userid从mongo中查找出店铺和商品列表
	var yqfutrueShop yqmongo.YqfutureShop

	// 根据uuid和userid查找出店铺
	shop, err := l.svcCtx.TsShopModel.FindOneByUuidAndUserId(l.ctx, in.UserId, in.Uuid)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找店铺失败", err)
		return nil, err
	}

	//TODO 如果mysql中店铺为空，说明是首次训练，需要从接口中获取店铺和商品，保存到mysql和mongo中，再开启训练
	if shop.Id == 0 {
		tsShop := &orm.TsShop{
			Id:           l.svcCtx.SnowFlakeNode.Generate().Int64(),
			UserId:       in.UserId,
			Uuid:         in.Uuid,
			ShopName:     yqfutrueShop.ShopName,
			PlatformType: yqfutrueShop.Platform,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
			CreateBy:     in.UserId,
			UpdateBy:     in.UserId,
		}
		_, err = l.svcCtx.TsShopModel.Insert(l.ctx, tsShop)
		if err != nil {
			l.Logger.Error("保存店铺到mysql失败", err)
			return nil, err
		}
		for _, goods := range yqfutrueShop.GoodsList {
			//TODO 通过商品id列表 调用第三方接口 获取商品json 解析图片数组 并一起保存到mongo 在这里轮询等待所有商品数据返回

			tsGoods := &orm.TsGoods{
				Id:         l.svcCtx.SnowFlakeNode.Generate().Int64(),
				ShopId:     tsShop.Id,
				PlatformId: goods.PlatformId,
				GoodsUrl:   goods.Url,
			}
			_, err = l.svcCtx.TsGoodsModel.Insert(l.ctx, tsGoods)
		}
	} else {
		//TODO 如果mysql中店铺不为空，那么将mongo中的数据更新到mysql中
		//TODO 首先将店铺状态置为训练中，然后开启训练

	}

	// 根据店铺shopId查找出商品列表，需要筛选出enabled字段为1的商品
	shopId := shop.Id
	goodsList, err := l.svcCtx.TsGoodsModel.FindEnabledListByShopId(l.ctx, shopId)
	if err != nil {
		return nil, err
	}

	//TODO 将商品列表推到消息队列
	for _, goods := range *goodsList {
		var goodsString string
		goodsString, err = utills.AnyToString(goods)
		if err != nil {
			return nil, err
		}
		err := l.svcCtx.KqPusherClient.Push(goodsString)
		if err != nil {
			l.Logger.Error("推送商品到kafka失败", goodsList)
			return nil, err
		}
	}

	// 修改店铺状态, 添加训练次数
	shop.TrainingStatus = 1
	shop.TrainingTimes += 1
	shop.UpdateTime = time.Now()
	shop.UpdateBy = in.UserId
	err = l.svcCtx.TsShopModel.Update(l.ctx, shop)
	if err != nil {
		l.Logger.Error("修改店铺状态失败", goodsList)
		return nil, err
	}
	// 修改商品状态, 添加训练次数
	for _, goods := range *goodsList {
		goods.TrainingStatus = 1
		goods.TrainingTimes += 1
		goods.UpdateTime = time.Now()
		goods.UpdateBy = in.UserId
		err = l.svcCtx.TsGoodsModel.Update(l.ctx, goods)
		if err != nil {
			l.Logger.Error("修改商品状态失败", goodsList)
			return nil, err
		}
	}

	// 返回正常
	return &training.TrainingShopResp{}, nil
}
