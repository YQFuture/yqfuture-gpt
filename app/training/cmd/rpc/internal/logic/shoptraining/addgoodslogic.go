package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGoodsLogic {
	return &AddGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加商品
func (l *AddGoodsLogic) AddGoods(in *training.AddGoodsReq) (*training.AddGoodsResp, error) {
	// 获取商品url列表
	list := in.List

	//TODO 调用GPT接口获取商品信息
	l.Logger.Info("修改店铺状态失败", list)

	// 保存商品
	for a := range list {
		l.Logger.Info("修改店铺状态失败", a)

	}

	return &training.AddGoodsResp{}, nil
}
