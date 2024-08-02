package orglogic

import (
	"context"
	"errors"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GivePowerShopAvgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGivePowerShopAvgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GivePowerShopAvgLogic {
	return &GivePowerShopAvgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GivePowerShopAvg 平均分配店铺算力
func (l *GivePowerShopAvgLogic) GivePowerShopAvg(in *user.GivePowerShopAvgReq) (*user.GivePowerShopAvgResp, error) {
	// 获取当前用户数据和团队数据
	bsUser, err := l.svcCtx.BsUserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Error("获取用户数据失败: ", err)
		return nil, err
	}
	bsOrg, err := l.svcCtx.BsOrganizationModel.FindOne(l.ctx, bsUser.NowOrgId)
	if err != nil {
		l.Logger.Error("获取团队数据失败: ", err)
		return nil, err
	}
	if bsOrg.OwnerId != bsUser.Id {
		l.Logger.Error("当前用户不是当前团队管理员")
		return nil, errors.New("只有团队管理员才能平均分配店铺算力")
	}

	// 获取当前团队店铺列表
	bsShopList, err := l.svcCtx.BsShopModel.FindListByOrgId(l.ctx, bsOrg.Id)
	if err != nil {
		l.Logger.Error("获取团队店铺列表失败: ", err)
		return nil, err
	}

	// 计算平均算力
	avgPower := bsOrg.MonthPowerLimit / int64(len(*bsShopList))
	err = l.svcCtx.BsShopModel.UpdateShopPowerAvg(l.ctx, bsOrg.Id, avgPower)
	if err != nil {
		l.Logger.Error("平均分配店铺算力失败: ", err)
		return nil, err
	}

	return &user.GivePowerShopAvgResp{}, nil
}
