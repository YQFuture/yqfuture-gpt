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

type GetGoodsTrainingProgressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商品训练进度
func NewGetGoodsTrainingProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsTrainingProgressLogic {
	return &GetGoodsTrainingProgressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGoodsTrainingProgressLogic) GetGoodsTrainingProgress(req *types.GetGoodsTrainingProgressReq) (resp *types.GetGoodsTrainingProgressResp, err error) {
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
	result, err := l.svcCtx.ShopTrainingClient.GetGoodsTrainingProgress(l.ctx, &training.GetGoodsTrainingProgressReq{
		UserId:  userId,
		GoodsId: goodsIdInt,
	})
	if err != nil {
		l.Logger.Error("获取商品训练进度失败", err)
		return nil, err
	}
	return &types.GetGoodsTrainingProgressResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "获取商品训练进度成功",
		},
		Data: &types.GoodsTrainingProgress{
			TrainingStatus: result.TrainingStatus,
			GoodsNum:       result.GoodsNum,
			Token:          result.Token,
			FileSize:       result.FileSize,
			Power:          result.Power,
			Time:           result.Time,
		},
	}, nil
}
