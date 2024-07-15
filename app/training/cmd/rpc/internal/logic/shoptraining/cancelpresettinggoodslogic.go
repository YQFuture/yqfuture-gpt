package shoptraininglogic

import (
	"context"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelPreSettingGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelPreSettingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPreSettingGoodsLogic {
	return &CancelPreSettingGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CancelPreSettingGoods 取消预设商品
func (l *CancelPreSettingGoodsLogic) CancelPreSettingGoods(in *training.CancelPreSettingGoodsReq) (*training.CancelPreSettingGoodsResp, error) {
	tsGoods, err := l.svcCtx.TsGoodsModel.FindOne(l.ctx, in.GoodsId)
	if err != nil {
		l.Logger.Error("查询商品失败", err)
		return nil, err
	}
	if tsGoods.TrainingStatus != consts.TrainingComplete {
		l.Logger.Error("商品状态不是预设状态")
		return nil, err
	}
	if err := CancelGoodsPreSetting(l.ctx, l.svcCtx, tsGoods, in.UserId); err != nil {
		l.Logger.Error("取消预设商品失败", err)
		return nil, err
	}
	return &training.CancelPreSettingGoodsResp{}, nil
}
