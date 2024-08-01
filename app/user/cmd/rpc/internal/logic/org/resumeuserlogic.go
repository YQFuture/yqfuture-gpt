package orglogic

import (
	"context"
	"errors"
	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResumeUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResumeUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResumeUserLogic {
	return &ResumeUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ResumeUser 恢复用户
func (l *ResumeUserLogic) ResumeUser(in *user.ResumeUserReq) (*user.ResumeUserResp, error) {
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
		return nil, errors.New("只有团队管理员才能恢复用户")
	}
	// 恢复用户
	err = l.svcCtx.BsUserOrgModel.ChangeStatusByUserIdAndOrgId(l.ctx, in.ResumeUserId, bsUser.NowOrgId, 1)
	if err != nil {
		l.Logger.Error("恢复用户失败: ", err)
		return nil, err
	}
	return &user.ResumeUserResp{}, nil
}
