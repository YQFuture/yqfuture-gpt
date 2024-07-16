package loginlogic

import (
	"context"
	"database/sql"
	"time"
	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"
	"yufuture-gpt/app/user/model/orm"
	"yufuture-gpt/common/consts"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Register 注册
func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 判断手机号是否注册
	bsUser, err := l.svcCtx.BsUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		return nil, err
	}
	if bsUser != nil {
		return &user.RegisterResp{
			Code: consts.PhoneIsRegistered,
		}, nil
	}

	// 构建新用户并保存
	now := time.Now()
	id := l.svcCtx.SnowFlakeNode.Generate().Int64()
	newBsUser := &orm.BsUser{
		Id: id,
		Phone: sql.NullString{
			String: in.Phone,
			Valid:  true,
		},
		UserName: sql.NullString{
			String: in.Phone,
			Valid:  true,
		},
		NickName: sql.NullString{
			String: in.Phone,
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
	return &user.RegisterResp{
		Code: consts.Success,
		Result: &user.UserInfo{
			Id:       newBsUser.Id,
			Phone:    newBsUser.Phone.String,
			NickName: newBsUser.NickName.String,
			HeadImg:  newBsUser.HeadImg.String,
		},
	}, nil
}
