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

type PreSettingGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 预训练商品
func NewPreSettingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingGoodsLogic {
	return &PreSettingGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreSettingGoodsLogic) PreSettingGoods(req *types.PreSettingGoodsReq) (resp *types.BaseResp, err error) {
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
	_, err = l.svcCtx.ShopTrainingClient.PreSettingGoods(l.ctx, &training.PreSettingGoodsReq{
		GoodsId: goodsIdInt,
		UserId:  userId,
	})
	if err != nil {
		l.Logger.Error("预训练商品失败", err)
		return nil, err
	}
	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "预训练商品成功",
	}, nil
}
