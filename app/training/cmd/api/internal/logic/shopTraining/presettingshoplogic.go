package shopTraining

import (
	"context"
	"encoding/json"
	"strings"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSettingShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPreSettingShopLogic 预设店铺
func NewPreSettingShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingShopLogic {
	return &PreSettingShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreSettingShopLogic) PreSettingShop(req *types.PresettingShopReq) (resp *types.PresettingShopResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	_, err = l.svcCtx.ShopTrainingClient.PreSettingShop(l.ctx, &training.PreSettingShopReq{
		UserId:        userId,
		Uuid:          req.Uuid,
		ShopName:      req.ShopName,
		PlatformType:  req.PlatFormType,
		Authorization: strings.TrimPrefix(req.Authorization, "Bearer "),
		Cookies:       req.Cookies,
	})
	// 后台会进行较长时间的轮询等待 所以会返回超时错误 无需处理 仅打印日志 前端关注店铺和商品的训练状态即可
	if err != nil {
		l.Logger.Error("开启预设店铺异常", err)
	}
	return &types.PresettingShopResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "success",
		},
	}, nil
}
