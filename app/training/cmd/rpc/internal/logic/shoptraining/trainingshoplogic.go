package shoptraininglogic

import (
	"context"
	"time"
	"yufuture-gpt/common/utills"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

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
	//根据uuid和userid查找出店铺
	shop, err := l.svcCtx.TsShopModel.FindOneByUuidAndUserId(l.ctx, in)
	if err != nil {
		return nil, err
	}

	//TODO 如果mysql中商品为空，说明是首次训练，需要从mongo中获取店铺和商品，保存到mysql中，再开启训练

	//根据店铺shopId查找出商品列表，需要筛选出enabled字段为1的商品
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

	//修改店铺状态, 添加训练次数
	shop.TrainingStatus = 1
	shop.TrainingTimes += 1
	shop.UpdateTime = time.Now()
	shop.UpdateBy = in.UserId
	err = l.svcCtx.TsShopModel.Update(l.ctx, shop)
	if err != nil {
		l.Logger.Error("修改店铺状态失败", goodsList)
		return nil, err
	}
	//修改商品状态, 添加训练次数
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

	//返回正常
	return &training.TrainingShopResp{}, nil
}
