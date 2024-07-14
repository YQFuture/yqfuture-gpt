package shopTraining

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelPreSettingShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCancelPreSettingShopLogic 取消店铺预设
func NewCancelPreSettingShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPreSettingShopLogic {
	return &CancelPreSettingShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelPreSettingShopLogic) CancelPreSettingShop(req *types.CancelPreSettingShopReq) (resp *types.CancelPreSettingShopResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	l.Logger.Info("userId", userId)
	_, err = l.svcCtx.ShopTrainingClient.CancelPreSettingShop(l.ctx, &training.CancelPreSettingShopReq{
		UserId: userId,
		Uuid:   req.Uuid,
	})
	if err != nil {
		l.Logger.Error("取消预设店铺失败", err)
		return nil, err
	}
	return &types.CancelPreSettingShopResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "取消预设成功",
		},
	}, nil
}
