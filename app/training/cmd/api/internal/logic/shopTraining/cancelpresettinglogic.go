package shopTraining

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"
	"yufuture-gpt/common/consts"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelPreSettingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消预设
func NewCancelPreSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPreSettingLogic {
	return &CancelPreSettingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelPreSettingLogic) CancelPreSetting(req *types.CancelPreSettingReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	l.Logger.Info("userId", userId)
	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "取消预设成功",
	}, nil
}
