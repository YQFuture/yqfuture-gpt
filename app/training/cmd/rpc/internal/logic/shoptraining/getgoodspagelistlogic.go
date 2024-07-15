package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGoodsPageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGoodsPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsPageListLogic {
	return &GetGoodsPageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGoodsPageList 查询商品列表
func (l *GetGoodsPageListLogic) GetGoodsPageList(in *training.GoodsPageListReq) (*training.GoodsPageListResp, error) {
	total, err := l.svcCtx.TsGoodsModel.GetGoodsPageTotal(l.ctx, in)
	if err != nil {
		l.Logger.Error("查询商品列表总数失败", err)
		return nil, err
	}
	list, err := l.svcCtx.TsGoodsModel.GetGoodsPageList(l.ctx, in)
	if err != nil {
		l.Logger.Error("查询商品列表失败", err)
		return nil, err
	}

	var goodRespList []*training.GoodsResp
	for _, goods := range *list {
		goodRespList = append(goodRespList, &training.GoodsResp{
			Id:              goods.Id,
			ShopId:          goods.ShopId,
			GoodsName:       goods.GoodsName,
			GoodsUrl:        goods.GoodsUrl,
			TrainingSummary: goods.TrainingSummary,
			Enabled:         goods.Enabled,
			TrainingTimes:   goods.TrainingTimes,
			TrainingStatus:  goods.TrainingStatus,
			UpdateTime:      goods.UpdateTime.Unix(),
		})
	}

	return &training.GoodsPageListResp{
		Total: int64(total),
		List:  goodRespList,
	}, nil
}
