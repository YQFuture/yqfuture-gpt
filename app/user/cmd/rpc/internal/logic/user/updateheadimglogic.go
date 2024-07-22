package userlogic

import (
	"context"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateHeadImgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateHeadImgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHeadImgLogic {
	return &UpdateHeadImgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateHeadImg 更新头像
func (l *UpdateHeadImgLogic) UpdateHeadImg(in *user.UpdateHeadImgReq) (*user.UpdateHeadImgResp, error) {
	err := l.svcCtx.BsUserModel.UpdateHeadImg(l.ctx, in.HeadImg, in.UserId)
	if err != nil {
		l.Logger.Error("更新头像失败: ", err)
		return nil, err
	}
	return &user.UpdateHeadImgResp{}, nil
}
