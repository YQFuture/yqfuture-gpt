package shopTraining

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelPreSettingGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消预设商品
func NewCancelPreSettingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPreSettingGoodsLogic {
	return &CancelPreSettingGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelPreSettingGoodsLogic) CancelPreSettingGoods(req *types.CancelPreSettingGoodsReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	goodsIdInt, err := strconv.ParseInt(req.GoodsId, 10, 64)
	if err != nil {
		l.Logger.Error("转换商品id失败", err)
		return nil, err
	}
	_, err = l.svcCtx.ShopTrainingClient.CancelPreSettingGoods(l.ctx, &training.CancelPreSettingGoodsReq{
		GoodsId: goodsIdInt,
		UserId:  userId,
	})
	if err != nil {
		l.Logger.Error("取消预设商品失败", err)
		return nil, err
	}
	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "取消预设商品成功",
	}, nil
}
