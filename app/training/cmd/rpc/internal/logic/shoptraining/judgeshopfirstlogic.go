package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type JudgeShopFirstLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJudgeShopFirstLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JudgeShopFirstLogic {
	return &JudgeShopFirstLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// JudgeShopFirst 判断店铺是否首次登录 即从未进行过训练
func (l *JudgeShopFirstLogic) JudgeShopFirst(in *training.JudgeShopFirstReq) (*training.JudgeShopFirstResp, error) {
	first, err := l.svcCtx.TsShopModel.JudgeShopFirst(l.ctx, in.UserId, in.Uuid)
	if err != nil {
		l.Logger.Error("判断店铺是否初次登录失败", err)
		return nil, err
	}
	return &training.JudgeShopFirstResp{
		First: int64(first),
	}, nil
}
