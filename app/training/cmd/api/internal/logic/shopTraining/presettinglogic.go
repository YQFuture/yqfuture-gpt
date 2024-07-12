package shopTraining

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSettingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingLogic {
	return &PreSettingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreSettingLogic) PreSetting(req *types.ShopTrainingReq) (resp *types.ShopTrainingResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}

	_, err = l.svcCtx.ShopTrainingClient.PreSetting(l.ctx, &training.ShopTrainingReq{
		UserId:        userId,
		Uuid:          req.Uuid,
		ShopName:      req.ShopName,
		PlatformType:  req.PlatFormType,
		Authorization: req.Authorization,
		Cookies:       req.Cookies,
	})
	if err != nil {
		l.Logger.Error("开启预训练失败", err)
		return nil, err
	}
	return &types.ShopTrainingResp{
		BaseResp: types.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, nil
}
