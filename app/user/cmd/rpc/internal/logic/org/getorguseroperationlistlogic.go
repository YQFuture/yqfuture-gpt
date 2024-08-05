package orglogic

import (
	"context"
	"errors"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgUserOperationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrgUserOperationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgUserOperationListLogic {
	return &GetOrgUserOperationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrgUserOperationList 获取组织用户操作记录列表
func (l *GetOrgUserOperationListLogic) GetOrgUserOperationList(in *user.OrgUserOperationListReq) (*user.OrgUserOperationListResp, error) {
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
		return nil, errors.New("只有团队管理员才能获取组织用户操作记录列表")
	}

	listResult, err := l.svcCtx.DbuseroperationModel.FindListByUserIdAndOrgId(l.ctx, in.OperationUserId, bsUser.NowOrgId, in.StartTime, in.EndTime, in.Query)
	if err != nil {
		l.Logger.Error("获取组织用户操作记录列表失败: ", err)
		return nil, err
	}

	orgUserOperationListResp := &user.OrgUserOperationListResp{}
	for _, operation := range *listResult {
		orgUserOperationListResp.List = append(orgUserOperationListResp.List, &user.OrgUserOperation{
			CreateTime:    operation.CreateAt.Unix(),
			OperationDesc: operation.OperationDesc,
		})
	}

	return orgUserOperationListResp, nil
}
