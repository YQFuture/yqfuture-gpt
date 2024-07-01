package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type JudgeFirstShopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJudgeFirstShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JudgeFirstShopLogic {
	return &JudgeFirstShopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 判断店铺是否初次登录(从未进行过训练)
func (l *JudgeFirstShopLogic) JudgeFirstShop(in *training.JudgeFirstShopReq) (*training.JudgeFirstShopResp, error) {
	first, err := l.svcCtx.TsShopModel.JudgeFirstShop(l.ctx, in)
	if err != nil {
		l.Logger.Error("判断店铺是否初次登录失败", err)
		return nil, err
	}
	return &training.JudgeFirstShopResp{
		First: int64(first),
	}, nil
}
