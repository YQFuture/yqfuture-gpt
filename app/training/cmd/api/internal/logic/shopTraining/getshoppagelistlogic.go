package shopTraining

import (
	"context"
	"strconv"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShopPageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShopPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopPageListLogic {
	return &GetShopPageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShopPageListLogic) GetShopPageList(req *types.ShopPageListReq) (resp *types.ShopPageListResp, err error) {
	userId := l.ctx.Value("userId")
	intNumber, err := strconv.ParseInt(userId.(string), 10, 64)
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	result, err := l.svcCtx.ShopTrainingClient.GetShopPageList(l.ctx, &training.ShopPageListReq{
		UserId:         intNumber,
		PageNum:        req.PageNum,
		PageSize:       req.PageSize,
		Query:          req.Query,
		PlatFormType:   req.PlatFormType,
		TrainingStatus: req.TrainingStatus,
		UpdateTime:     req.UpdateTime,
	})
	if err != nil {
		l.Logger.Error("查询店铺列表失败", err)
		return nil, err
	}

	var list []*types.ShopResp

	for _, value := range result.List {
		list = append(list, &types.ShopResp{
			Id:             value.Id,
			Uuid:           value.Uuid,
			ShopName:       value.ShopName,
			PlatFormType:   value.PlatformType,
			TrainingStatus: value.TrainingStatus,
			TrainingTimes:  value.TrainingTimes,
			UpdateTime:     value.UpdateTime,
		})
	}
	return &types.ShopPageListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "查询店铺列表成功",
		},
		Data: &types.ShopListResp{
			BasePageResp: types.BasePageResp{
				Total:    result.Total,
				PageNum:  req.PageNum,
				PageSize: req.PageSize,
			},
			List: list,
		},
	}, nil

}
