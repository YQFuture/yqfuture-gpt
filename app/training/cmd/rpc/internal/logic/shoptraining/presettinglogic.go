package shoptraininglogic

import (
	"context"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/utills"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingLogic {
	return &PreSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreSettingLogic) PreSetting(in *training.ShopTrainingReq) (*training.ShopTrainingResp, error) {
	list, err := l.svcCtx.TsShopModel.FindList(l.ctx)
	if err != nil {
		l.Logger.Error("查询店铺列表失败", err)
		return nil, err
	}
	//将返回体转字符串
	result, err := utills.AnyToString(list)
	if err != nil {
		l.Logger.Error("序列化数据失败", err)
		return nil, err
	}
	return &training.ShopTrainingResp{
		Result: result,
	}, nil
}
