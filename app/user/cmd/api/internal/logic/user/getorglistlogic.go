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

type GetOrgListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetOrgListLogic 获取用户组织列表
func NewGetOrgListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgListLogic {
	return &GetOrgListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrgListLogic) GetOrgList(req *types.BaseReq) (resp *types.OrgListResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户组织列表失败", err)
		return &types.OrgListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取失败",
			},
		}, nil
	}

	// 调用RPC接口 获取用户组织列表
	orgListResp, err := l.svcCtx.UserClient.GetOrgList(l.ctx, &user.OrgListReq{
		UserId: userId,
	})
	if err != nil {
		l.Logger.Error("获取用户组织列表失败", err)
		return &types.OrgListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取失败",
			},
		}, nil
	}
	var orgList []types.OrgInfo
	for _, org := range orgListResp.Result {
		orgList = append(orgList, types.OrgInfo{
			OrgId:      strconv.FormatInt(org.OrgId, 10),
			OrgName:    org.OrgName,
			BundleType: org.BundleType,
			IsAdmin:    org.IsAdmin,
		})
	}

	return &types.OrgListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "获取成功",
		},
		Data: orgList,
	}, nil
}
