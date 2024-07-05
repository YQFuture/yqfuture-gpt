package shopTraining

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/utils"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

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
	l.Logger.Info("token中的userId", userId)
	result, err := l.svcCtx.ShopTrainingClient.PreSetting(l.ctx, &training.ShopTrainingReq{})
	if err != nil {
		l.Logger.Error("获取店铺列表失败", err)
		return nil, err
	}
	var shopList []*orm.TsShop
	err = utils.StringToAny(result.Result, &shopList)
	if err != nil {
		l.Logger.Error("反序列化数据失败", err)
		return nil, err
	}
	return &types.ShopTrainingResp{
		BaseResp: types.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		Data: shopList,
	}, nil
}
