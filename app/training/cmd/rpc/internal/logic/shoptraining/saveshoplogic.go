package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveShopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveShopLogic {
	return &SaveShopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存爬取的店铺基本数据
func (l *SaveShopLogic) SaveShop(in *training.SaveShopReq) (*training.SaveShopResp, error) {
	// todo: add your logic here and delete this line

	return &training.SaveShopResp{}, nil
}
