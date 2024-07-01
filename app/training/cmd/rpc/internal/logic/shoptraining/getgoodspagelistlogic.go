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

// 查询商品列表
func (l *GetGoodsPageListLogic) GetGoodsPageList(in *training.GoodsPageListReq) (*training.GoodsPageListResp, error) {
	// todo: add your logic here and delete this line

	return &training.GoodsPageListResp{}, nil
}
