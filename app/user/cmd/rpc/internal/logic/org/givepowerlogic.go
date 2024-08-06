package orglogic

import (
	"context"
	"errors"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GivePowerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGivePowerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GivePowerLogic {
	return &GivePowerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GivePower 分配算力
func (l *GivePowerLogic) GivePower(in *user.GivePowerReq) (*user.GivePowerResp, error) {
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
		return nil, errors.New("只有团队管理员才能分配算力")
	}

	// 获取当前团队已分配的总算力 判断剩余算力是否足够
	totalPower, err := l.svcCtx.BsUserOrgModel.FindOrgTotalGivePower(l.ctx, bsOrg.Id)
	if err != nil {
		l.Logger.Error("获取团队已分配算力失败: ", err)
		return nil, err
	}
	if totalPower+in.Power > bsOrg.MonthPowerLimit {
		l.Logger.Error("当前团队已分配算力不足")
		return &user.GivePowerResp{
			Code: consts.PowerNotEnough,
		}, nil
	}

	// 更新用户算力
	err = l.svcCtx.BsUserOrgModel.UpdateUserPower(l.ctx, in.Power, in.GivePowerUserId, bsOrg.Id)
	if err != nil {
		l.Logger.Error("更新用户算力失败: ", err)
		return nil, err
	}

	return &user.GivePowerResp{}, nil
}
