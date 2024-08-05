package org

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/user/cmd/rpc/client/org"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgUserOperationPageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetOrgUserOperationPageListLogic 获取团队用户操作记录分页列表
func NewGetOrgUserOperationPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgUserOperationPageListLogic {
	return &GetOrgUserOperationPageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrgUserOperationPageListLogic) GetOrgUserOperationPageList(req *types.OrgUserOperationPageListReq) (resp *types.OrgUserOperationPageListResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.OrgUserOperationPageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	operationUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		l.Logger.Error("获取要查询角色的用户ID失败", err)
		return &types.OrgUserOperationPageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 获取团队用户操作记录分页列表
	userOperationResp, err := l.svcCtx.OrgClient.GetOrgUserOperationPageList(l.ctx, &org.OrgUserOperationPageListReq{
		UserId:          userId,
		OperationUserId: operationUserId,
		Query:           req.Query,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		PageSize:        req.PageSize,
		PageNum:         req.PageNum,
	})
	if err != nil {
		l.Logger.Error("获取团队用户操作记录分页列表失败", err)
		return &types.OrgUserOperationPageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	orgUserOperationPageListResp := &types.OrgUserOperationPageListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: types.OrgUserOperationPage{
			BasePageResp: types.BasePageResp{
				PageSize: userOperationResp.PageSize,
				PageNum:  userOperationResp.PageNum,
				Total:    userOperationResp.Total,
			},
		},
	}

	for _, userOperation := range userOperationResp.List {
		orgUserOperationPageListResp.Data.List = append(orgUserOperationPageListResp.Data.List, types.OrgUserOperation{
			CreateTime:    userOperation.CreateTime,
			OperationDesc: userOperation.OperationDesc,
		})
	}

	return orgUserOperationPageListResp, nil
}
