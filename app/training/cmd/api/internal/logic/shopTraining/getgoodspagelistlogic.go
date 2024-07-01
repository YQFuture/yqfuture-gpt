package shopTraining

import (
	"context"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGoodsPageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGoodsPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsPageListLogic {
	return &GetGoodsPageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGoodsPageListLogic) GetGoodsPageList(req *types.GoodsPageListReq) (resp *types.GoodsPageListResp, err error) {

	result, err := l.svcCtx.ShopTrainingClient.GetGoodsPageList(l.ctx, &training.GoodsPageListReq{
		ShopId:     req.ShopId,
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		Query:      req.Query,
		Enabled:    req.Enabled,
		UpdateTime: req.UpdateTime,
	})
	if err != nil {
		l.Logger.Error("查询商品列表失败", err)
		return nil, err
	}

	var list []*types.GoodsResp

	for _, value := range result.List {
		list = append(list, &types.GoodsResp{
			Id:              value.Id,
			ShopId:          value.ShopId,
			GoodsName:       value.GoodsName,
			GoodsUrl:        value.GoodsUrl,
			TrainingSummary: value.TrainingSummary,
			TrainingStatus:  value.TrainingStatus,
			TrainingTimes:   value.TrainingTimes,
			UpdateTime:      value.UpdateTime,
			Enabled:         value.Enabled,
		})
	}
	return &types.GoodsPageListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "查询商品列表成功",
		},
		Data: &types.GoodsListResp{
			BasePageResp: types.BasePageResp{
				Total:    result.Total,
				PageNum:  req.PageNum,
				PageSize: req.PageSize,
			},
			List: list,
		},
	}, nil
}
