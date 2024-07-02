package shopTraining

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JudgeFirstShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 判断店铺是否初次登录(从未进行过训练)
func NewJudgeFirstShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JudgeFirstShopLogic {
	return &JudgeFirstShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JudgeFirstShopLogic) JudgeFirstShop(req *types.JudgeFirstShopReq) (resp *types.JudgeFirstShopResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	first, err := l.svcCtx.ShopTrainingClient.JudgeFirstShop(l.ctx, &training.JudgeFirstShopReq{
		Uuid:   req.Uuid,
		UserId: userId,
	})
	if err != nil {
		l.Logger.Error("判断店铺是否初次登录失败", err)
		return nil, err
	}
	return &types.JudgeFirstShopResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "查询店铺列表成功",
		},
		Data: first.First,
	}, nil
}
