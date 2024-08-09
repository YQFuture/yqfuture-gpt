package orglogic

import (
	"context"
	"errors"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GivePowerAvgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGivePowerAvgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GivePowerAvgLogic {
	return &GivePowerAvgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GivePowerAvg 平均分配算力
func (l *GivePowerAvgLogic) GivePowerAvg(in *user.GivePowerAvgReq) (*user.GivePowerAvgResp, error) {
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
		return nil, errors.New("只有团队管理员才能平均分配算力")
	}

	// 获取当前团队用户列表
	userOrgList, err := l.svcCtx.BsUserOrgModel.FindListByOrgId(l.ctx, bsUser.NowOrgId)
	if err != nil {
		l.Logger.Error("获取团队用户列表失败: ", err)
		return nil, err
	}

	// 计算平均算力
	avgPower := bsOrg.MonthPowerLimit / int64(len(*userOrgList)) / 10000 * 10000
	err = l.svcCtx.BsUserOrgModel.UpdateUserPowerAvg(l.ctx, bsUser.NowOrgId, avgPower)
	if err != nil {
		l.Logger.Error("更新用户算力失败: ", err)
		return nil, err

	}
	return &user.GivePowerAvgResp{}, nil
}
