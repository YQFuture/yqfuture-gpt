package userlogic

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	loginlogic "yufuture-gpt/app/user/cmd/rpc/internal/logic/login"
	"yufuture-gpt/app/user/model/orm"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindPhoneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindPhoneLogic {
	return &BindPhoneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BindPhone 绑定手机号
func (l *BindPhoneLogic) BindPhone(in *user.BindPhoneReq) (*user.BindPhoneResp, error) {
	bsUser, err := l.svcCtx.BsUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		l.Logger.Error("根据手机号获取用户失败", err)
		return nil, err
	}
	// 如果手机号未注册 创建新用户 并返回当前用户ID
	if bsUser == nil {
		// 先保存权限相关数据到MongoDB 失败直接返回错误
		// 从MySQL中获取权限模板
		bsPermTemplateList, err := l.svcCtx.BsPermTemplateModel.FindListByBundleType(l.ctx, 0)
		if err != nil {
			l.Logger.Error("获取权限模板失败", err)
			return nil, err
		}
		// 根据权限模板构建MongoDB文档
		dborgpermission := loginlogic.BuildDefaultMongoPermDoc(*bsPermTemplateList)
		// 保存MongoDB文档 获取返回的ID
		result, err := l.svcCtx.DborgpermissionModel.InsertOne(l.ctx, dborgpermission)
		if err != nil {
			l.Logger.Error("保存MongoDB文档失败", err)
			return nil, err
		}
		oid, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			l.Logger.Error("获取MongoDB文档ID失败", err)
			return nil, err
		}
		mongoPermId := oid.Hex()

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
			Openid: sql.NullString{
				String: in.Openid,
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
			BundleType:      0,
			MaxSeat:         1,
			MonthPowerLimit: 0,
			MonthUsedPower:  0,
			MongoPermId:     mongoPermId,
			CreateTime:      now,
			UpdateTime:      now,
			CreateBy:        userId,
			UpdateBy:        userId,
		}
		// 构建用户组织中间表
		bsUserOrg := &orm.BsUserOrg{
			UserId:     userId,
			OrgId:      orgId,
			Status:     1,
			CreateTime: now,
			UpdateTime: now,
			CreateBy:   userId,
			UpdateBy:   userId,
		}

		// 在同一个事务中保存三张表的数据
		err = l.svcCtx.BsUserModel.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			_, err := l.svcCtx.BsUserModel.SessionInsert(l.ctx, newBsUser, session)
			if err != nil {
				l.Logger.Error("保存用户信息失败: ", err)
				return err
			}
			_, err = l.svcCtx.BsOrganizationModel.SessionInsert(l.ctx, bsOrganization, session)
			if err != nil {
				l.Logger.Error("保存组织信息失败: ", err)
				return err
			}
			_, err = l.svcCtx.BsUserOrgModel.SessionInsert(l.ctx, bsUserOrg, session)
			if err != nil {
				l.Logger.Error("保存用户组织中间表失败: ", err)
				return err
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		return &user.BindPhoneResp{
			UserId: userId,
		}, nil
	}

	// 如果手机号已注册 将该微信OpenId绑定到手机用户 并返回手机用户ID
	err = l.svcCtx.BsUserModel.BindOpenId(l.ctx, in.Openid, bsUser.Id)
	if err != nil {
		return nil, err
	}
	return &user.BindPhoneResp{
		UserId: bsUser.Id,
	}, nil
}
