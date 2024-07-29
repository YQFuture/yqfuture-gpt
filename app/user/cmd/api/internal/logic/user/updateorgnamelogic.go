package user

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/user/cmd/rpc/client/user"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrgNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateOrgNameLogic 更新组织名称
func NewUpdateOrgNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrgNameLogic {
	return &UpdateOrgNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrgNameLogic) UpdateOrgName(req *types.UpdateOrgNameReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("更新组织名称失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "更新失败",
		}, nil
	}

	orgId, err := strconv.ParseInt(req.OrgId, 10, 64)
	if err != nil {
		l.Logger.Error("更新组织名称失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "更新失败",
		}, nil
	}

	// 调用RPC接口 更新组织名称
	updateOrgNameResp, err := l.svcCtx.UserClient.UpdateOrgName(l.ctx, &user.UpdateOrgNameReq{
		UserId:  userId,
		OrgId:   orgId,
		OrgName: req.OrgName,
	})
	if err != nil {
		l.Logger.Error("更新组织名称失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "更新失败",
		}, nil
	}
	if updateOrgNameResp.Code == consts.OrgNameIsExist {
		return &types.BaseResp{
			Code: consts.OrgNameIsExist,
			Msg:  "组织名称已存在",
		}, nil
	}

	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "更新成功",
	}, nil
}
