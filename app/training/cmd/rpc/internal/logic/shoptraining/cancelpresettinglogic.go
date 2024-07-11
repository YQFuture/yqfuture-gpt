package shoptraininglogic

import (
	"context"
	"time"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelPreSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelPreSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPreSettingLogic {
	return &CancelPreSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消预训练
func (l *CancelPreSettingLogic) CancelPreSetting(in *training.CancelPreSettingReq) (*training.ShopTrainingResp, error) {
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
	if tsShop.TrainingStatus != consts.TrainingComplete {
		l.Logger.Error("店铺状态不是预训练状态")
		return nil, err
	}
	err = CancelShopPreSetting(l.ctx, l.svcCtx, tsShop, in.UserId)
	if err != nil {
		l.Logger.Error("取消预训练店铺成功", err)
		return nil, err
	}
	for _, tsGoods := range *tsGoodsList {
		// 排除掉不是预训练完成的商品
		if tsGoods.TrainingStatus != consts.TrainingComplete {
			continue
		}
		err = CancelGoodsPreSetting(l.ctx, l.svcCtx, tsGoods, in.UserId)
		if err != nil {
			l.Logger.Error("取消预训练店铺失败", err)
			return nil, err
		}
	}

	return &training.ShopTrainingResp{}, nil
}

func CancelShopPreSetting(ctx context.Context, svcCtx *svc.ServiceContext, tsShop *orm.TsShop, userId int64) error {
	tsShop.TrainingStatus = consts.Undefined
	tsShop.UpdateTime = time.Now()
	tsShop.UpdateBy = userId
	err := svcCtx.TsShopModel.Update(ctx, tsShop)
	if err != nil {
		return err
	}
	return nil
}

func CancelGoodsPreSetting(ctx context.Context, svcCtx *svc.ServiceContext, tsGoods *orm.TsGoods, userId int64) error {
	tsGoods.TrainingStatus = consts.Undefined
	tsGoods.UpdateTime = time.Now()
	tsGoods.UpdateBy = userId
	err := svcCtx.TsGoodsModel.Update(ctx, tsGoods)
	if err != nil {
		return err
	}
	return nil
}
