package loginlogic

import (
	"context"
	"database/sql"
	"time"
	"yufuture-gpt/app/user/model/orm"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWechatUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWechatUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWechatUserInfoLogic {
	return &GetWechatUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetWechatUserInfo 获取微信用户信息
func (l *GetWechatUserInfoLogic) GetWechatUserInfo(in *user.WechatUserInfoReq) (*user.WechatUserInfoResp, error) {
	// 从用户表查找用户 并判断用户是否存在
	bsUser, err := l.svcCtx.BsUserModel.FindOneByPhone(l.ctx, in.Openid)
	if err != nil {
		return nil, err
	}
	// 已注册的用户直接返回
	if bsUser != nil {
		// 返回用户信息
		return &user.WechatUserInfoResp{
			Code: consts.Success,
			Result: &user.UserInfo{
				Id:       bsUser.Id,
				Phone:    bsUser.Phone.String,
				NickName: bsUser.NickName.String,
				HeadImg:  bsUser.HeadImg.String,
			},
		}, nil
	}

	// 未注册的用户 先创建用户再返回
	now := time.Now()
	id := l.svcCtx.SnowFlakeNode.Generate().Int64()
	newBsUser := &orm.BsUser{
		Id: id,
		UserName: sql.NullString{
			String: in.Openid,
			Valid:  true,
		},
		NickName: sql.NullString{
			String: in.Openid,
			Valid:  true,
		},
		HeadImg: sql.NullString{
			String: "https://yqfuture.com/_nuxt/favicon.28e9763f.png",
			Valid:  true,
		},
		CreateTime: now,
		UpdateTime: now,
		CreateBy:   id,
		UpdateBy:   id,
	}
	_, err = l.svcCtx.BsUserModel.Insert(l.ctx, newBsUser)
	if err != nil {
		return nil, err
	}

	// 返回用户信息
	return &user.WechatUserInfoResp{
		Code: consts.Success,
		Result: &user.UserInfo{
			Id:       newBsUser.Id,
			Phone:    newBsUser.Phone.String,
			NickName: newBsUser.NickName.String,
			HeadImg:  newBsUser.HeadImg.String,
		},
	}, nil
}
