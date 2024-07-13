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

type JudgeShopFirstLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewJudgeShopFirstLogic 判断店铺是否首次登录 即从未进行过训练
func NewJudgeShopFirstLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JudgeShopFirstLogic {
	return &JudgeShopFirstLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JudgeShopFirstLogic) JudgeShopFirst(req *types.JudgeShopFirstReq) (resp *types.JudgeShopFirstResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	result, err := l.svcCtx.ShopTrainingClient.JudgeShopFirst(l.ctx, &training.JudgeShopFirstReq{
		Uuid:   req.Uuid,
		UserId: userId,
	})
	if err != nil {
		l.Logger.Error("判断店铺是否首次登录失败", err)
		return nil, err
	}
	return &types.JudgeShopFirstResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "判断店铺是否首次登录成功",
		},
		Data: result.First,
	}, nil
}
