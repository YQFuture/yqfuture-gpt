package userlogic

import (
	"context"
	"errors"
	"yufuture-gpt/app/user/model/orm"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCurrentUserDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserDataLogic {
	return &GetCurrentUserDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCurrentUserData 获取当前登录用户数据
func (l *GetCurrentUserDataLogic) GetCurrentUserData(in *user.CurrentUserDataReq) (*user.CurrentUserDataResp, error) {
	bsUser, err := l.svcCtx.BsUserModel.FindOne(l.ctx, in.UserId)
	if err != nil && !errors.Is(err, orm.ErrNotFound) {
		l.Logger.Error("获取用户失败", err)
		return nil, err
	}
	// 判断是否未绑定手机号
	if errors.Is(err, orm.ErrNotFound) || bsUser.Phone.Valid == false || bsUser.Phone.String == "" {
		return &user.CurrentUserDataResp{
			Code: consts.PhoneIsNotBound,
		}, nil
	}
	// 查询用户组织关系
	bsUserOrg, err := l.svcCtx.BsUserOrgModel.FindUserOrgByUserIdAndOrgId(l.ctx, bsUser.Id, bsUser.NowOrgId)
	if err != nil {
		l.Logger.Error("查询用户组织关系失败", err)
		return nil, err
	}
	var orgUserStatus int64
	if bsUserOrg == nil {
		// 中间表关系不存在 说明被删除
		orgUserStatus = 2
	} else if bsUserOrg.Status == 0 {
		// 状态为0 说明被暂停权限
		orgUserStatus = 1
	}

	// 查询当前组织信息
	organization, err := l.svcCtx.BsOrganizationModel.FindOne(l.ctx, bsUser.NowOrgId)
	if err != nil {
		l.Logger.Error("查询当前组织信息失败", err)
		return nil, err
	}

	// 更新全体公告到消息表
	err = l.svcCtx.BsMessageModel.SyncNotice(l.ctx, bsUser.Id)
	if err != nil {
		l.Logger.Error("更新全体公告到消息表失败", err)
	}

	// 查询当前未读消息数
	var unreadMsgCount int64
	count, err := l.svcCtx.BsMessageModel.FindUnreadCount(l.ctx, bsUser.Id, bsUser.NowOrgId)
	if err != nil {
		l.Logger.Error("查询当前未读消息数失败", err)
	} else {
		unreadMsgCount = count
	}

	return &user.CurrentUserDataResp{
		Code: consts.Success,
		Result: &user.CurrentUserData{
			Id:       bsUser.Id,
			HeadImg:  bsUser.HeadImg.String,
			NickName: bsUser.NickName.String,
			Phone:    bsUser.Phone.String,
			NowOrg: &user.OrgInfo{
				OrgId:      organization.Id,
				OrgName:    organization.OrgName.String,
				BundleType: organization.BundleType,
				IsAdmin:    organization.OwnerId == bsUser.Id,
			},
			UnreadMsgCount: unreadMsgCount,
			OrgUserStatus:  orgUserStatus,
		},
	}, nil
}
