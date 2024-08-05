package loginlogic

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"
	model "yufuture-gpt/app/user/model/mongo"
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

	// 先保存权限相关数据到MongoDB 失败直接返回错误
	// 从MySQL中获取权限模板
	bsPermTemplateList, err := l.svcCtx.BsPermTemplateModel.FindListByBundleType(l.ctx, 0)
	if err != nil {
		l.Logger.Error("获取权限模板失败", err)
		return nil, err
	}
	// 根据权限模板构建MongoDB文档
	dborgpermission := BuildDefaultMongoPermDoc(*bsPermTemplateList)
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
		l.Logger.Error("保存用户信息失败", err)
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

// BuildDefaultMongoPermDoc 构建默认的MongoDB组织权限文档 仅在创建新用户时调用 所以套餐类型一定为免费版
func BuildDefaultMongoPermDoc(bsPermTemplateList []*orm.BsPermTemplate) *model.Dborgpermission {
	var permissionList []*model.Permission
	for _, bsPermTemplate := range bsPermTemplateList {
		permissionList = append(permissionList, BuildGeneralPermission(bsPermTemplate, bsPermTemplate.Id, 0))
	}
	return &model.Dborgpermission{
		PermissionList: permissionList,
		RoleList:       []*model.Role{},
		UserList:       []*model.User{},
	}
}

func BuildGeneralPermission(bsPermTemplate *orm.BsPermTemplate, id, resourceId int64) *model.Permission {
	return &model.Permission{
		Id:         id,
		Name:       bsPermTemplate.Name,
		ParentId:   bsPermTemplate.ParentId,
		Perm:       bsPermTemplate.Perm,
		Url:        bsPermTemplate.Url.String,
		Sort:       bsPermTemplate.Sort,
		ResourceId: resourceId,
		TemplateId: bsPermTemplate.Id,
	}
}
