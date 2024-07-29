package user

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/user/cmd/rpc/client/user"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateHeadImgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateHeadImgLogic 更新头像
func NewUpdateHeadImgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHeadImgLogic {
	return &UpdateHeadImgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateHeadImgLogic) UpdateHeadImg(req *types.UpdateHeadImgReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("更新头像失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "更新失败",
		}, nil
	}

	// 调用RPC接口 更新头像
	_, err = l.svcCtx.UserClient.UpdateHeadImg(l.ctx, &user.UpdateHeadImgReq{
		UserId:  userId,
		HeadImg: req.HeadImg,
	})
	if err != nil {
		l.Logger.Error("更新头像失败", err)
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
