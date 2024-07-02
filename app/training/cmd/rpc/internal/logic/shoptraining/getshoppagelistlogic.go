package shoptraininglogic

import (
	"context"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShopPageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetShopPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopPageListLogic {
	return &GetShopPageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetShopPageList 查询店铺列表
func (l *GetShopPageListLogic) GetShopPageList(in *training.ShopPageListReq) (*training.ShopPageListResp, error) {
	total, err := l.svcCtx.TsShopModel.GetShopPageTotal(l.ctx, in)
	if err != nil {
		l.Logger.Error("查询店铺列表总数失败", err)
		return nil, err
	}
	list, err := l.svcCtx.TsShopModel.GetShopPageList(l.ctx, in)
	if err != nil {
		l.Logger.Error("查询店铺列表失败", err)
		return nil, err
	}

	var shopRespList []*training.ShopResp
	for _, shop := range *list {
		shopRespList = append(shopRespList, &training.ShopResp{
			Id:             shop.Id,
			Uuid:           shop.Uuid,
			ShopName:       shop.ShopName,
			PlatformType:   shop.PlatformType,
			TrainingTimes:  shop.TrainingTimes,
			TrainingStatus: shop.TrainingStatus,
			UpdateTime:     shop.UpdateTime.Unix(),
		})
	}

	return &training.ShopPageListResp{
		Total: int64(total),
		List:  shopRespList,
	}, nil
}
