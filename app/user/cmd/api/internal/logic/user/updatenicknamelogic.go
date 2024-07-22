package user

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNickNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateNickNameLogic 更新昵称
func NewUpdateNickNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNickNameLogic {
	return &UpdateNickNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNickNameLogic) UpdateNickName(req *types.UpdateNickNameReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("更新昵称失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "更新失败",
		}, nil
	}

	// 调用RPC接口 更新昵称
	_, err = l.svcCtx.UserClient.UpdateNickName(l.ctx, &user.UpdateNickNameReq{
		UserId:   userId,
		NickName: req.NickName,
	})
	if err != nil {
		l.Logger.Error("更新昵称失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "更新失败",
		}, nil
	}

	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "更新成功",
	}, nil
}
