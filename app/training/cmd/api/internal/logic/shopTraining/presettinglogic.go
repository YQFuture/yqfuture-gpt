package shopTraining

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

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
	// 后台会进行较长时间的轮询等待 所以会返回超时错误 无需处理 仅打印日志 前端关注店铺和商品的训练状态即可
	if err != nil {
		l.Logger.Error("开启预训练店铺异常", err)
	}
	return &types.ShopTrainingResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "success",
		},
	}, nil
}
