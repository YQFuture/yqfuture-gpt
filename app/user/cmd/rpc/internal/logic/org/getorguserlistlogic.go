package orglogic

import (
	"context"
	"errors"
	"yufuture-gpt/app/user/cmd/rpc/client/org"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrgUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgUserListLogic {
	return &GetOrgUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrgUserList 获取团队用户列表
func (l *GetOrgUserListLogic) GetOrgUserList(in *user.OrgUserListReq) (*user.OrgUserListResp, error) {
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
		return nil, errors.New("只有团队管理员才能获取团队用户列表")
	}

	userListResult, err := l.svcCtx.BsUserModel.FindListByOrgId(l.ctx, bsOrg.Id)
	if err != nil {
		l.Logger.Error("获取团队用户列表失败: ", err)
		return nil, err
	}

	var orgUserList []*org.OrgUser
	for _, userResult := range *userListResult {
		orgUser := &org.OrgUser{
			UserId:   userResult.Id,
			Phone:    userResult.Phone.String,
			NickName: userResult.NickName.String,
			HeadImg:  userResult.HeadImg.String,
		}
		orgUserList = append(orgUserList, orgUser)
	}

	return &user.OrgUserListResp{
		List: orgUserList,
	}, nil
}
