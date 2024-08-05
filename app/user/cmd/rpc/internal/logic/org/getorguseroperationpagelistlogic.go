package orglogic

import (
	"context"
	"errors"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgUserOperationPageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrgUserOperationPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgUserOperationPageListLogic {
	return &GetOrgUserOperationPageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrgUserOperationPageList 获取组织用户操作记录分页列表
func (l *GetOrgUserOperationPageListLogic) GetOrgUserOperationPageList(in *user.OrgUserOperationPageListReq) (*user.OrgUserOperationPageListResp, error) {
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
		return nil, errors.New("只有团队管理员才能获取组织用户操作记录分页列表")
	}

	total, err := l.svcCtx.DbuseroperationModel.FindPageTotalByUserIdAndOrgId(l.ctx, in.OperationUserId, bsUser.NowOrgId, in.StartTime, in.EndTime, in.Query)
	if err != nil {
		l.Logger.Error("获取组织用户操作记录分页列表总数失败: ", err)
		return nil, err
	}
	pageListResult, err := l.svcCtx.DbuseroperationModel.FindPageListByUserIdAndOrgId(l.ctx, in.OperationUserId, bsUser.NowOrgId, in.PageNum, in.PageSize, in.StartTime, in.EndTime, in.Query)
	if err != nil {
		l.Logger.Error("获取组织用户操作记录分页列表失败: ", err)
		return nil, err
	}

	orgUserOperationPageListResp := &user.OrgUserOperationPageListResp{
		PageNum:  in.PageNum,
		PageSize: in.PageSize,
		Total:    total,
	}
	for _, operation := range *pageListResult {
		orgUserOperationPageListResp.List = append(orgUserOperationPageListResp.List, &user.OrgUserOperation{
			OperationDesc: operation.OperationDesc,
			CreateTime:    operation.CreateAt.Unix(),
		})
	}

	return orgUserOperationPageListResp, nil
}
