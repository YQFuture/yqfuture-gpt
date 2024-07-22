package user

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 切换组织
func NewChangeOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeOrgLogic {
	return &ChangeOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeOrgLogic) ChangeOrg(req *types.ChangeOrgReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("切换组织失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "切换失败",
		}, nil
	}

	orgId, err := strconv.ParseInt(req.OrgId, 10, 64)
	if err != nil {
		l.Logger.Error("切换组织失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "切换失败",
		}, nil
	}

	// 调用RPC接口 切换用户当前组织
	changeOrgResp, err := l.svcCtx.UserClient.ChangeOrg(l.ctx, &user.ChangeOrgReq{
		UserId: userId,
		OrgId:  orgId,
	})
	if err != nil {
		l.Logger.Error("切换组织失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "切换失败",
		}, nil
	}

	if changeOrgResp.Code == consts.UserNotInOrg {
		l.Logger.Error("用户不在组织中", err)
		return &types.BaseResp{
			Code: consts.UserNotInOrg,
			Msg:  "用户不在组织中",
		}, nil
	}

	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "切换成功",
	}, nil
}
