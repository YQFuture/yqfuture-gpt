package loginlogic

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
		l.Logger.Error("根据手机号获取用户失败", err)
		return nil, err
	}
	if bsUser != nil {
		return &user.RegisterResp{
			Code: consts.PhoneIsRegistered,
		}, nil
	}

	// 构建新用户
	now := time.Now()
	userId := l.svcCtx.SnowFlakeNode.Generate().Int64()
	orgId := l.svcCtx.SnowFlakeNode.Generate().Int64()
	newBsUser := &orm.BsUser{
		Id:       userId,
		NowOrgId: orgId,
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
			String: "2e61e107-b98d-47a2-92c5-aec081f03978_head_img_default.jpg",
			Valid:  true,
		},
		CreateTime: now,
		UpdateTime: now,
		CreateBy:   userId,
		UpdateBy:   userId,
	}
	// 构建用户对应的组织
	bsOrganization := &orm.BsOrganization{
		Id:      orgId,
		OwnerId: userId,
		OrgName: sql.NullString{
			String: in.Phone + "的组织",
			Valid:  true,
		},
		BundleType: 0,
		CreateTime: now,
		UpdateTime: now,
		CreateBy:   userId,
		UpdateBy:   userId,
	}
	// 构建用户组织中间表
	bsUserOrg := &orm.BsUserOrg{
		UserId: userId,
		OrgId:  orgId,
	}

	// 在同一个事务中保存三张表的数据
	err = l.svcCtx.BsUserModel.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err = l.svcCtx.BsUserModel.SessionInsert(l.ctx, newBsUser, session)
		if err != nil {
			return err
		}
		_, err := l.svcCtx.BsOrganizationModel.SessionInsert(l.ctx, bsOrganization, session)
		if err != nil {
			return err
		}
		_, err = l.svcCtx.BsUserOrgModel.SessionInsert(l.ctx, bsUserOrg, session)
		if err != nil {
			return err
		}
		return nil
	})

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
