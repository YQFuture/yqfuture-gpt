package orglogic

import (
	"context"
	"strconv"
	"time"
	model "yufuture-gpt/app/user/model/mongo"
	"yufuture-gpt/app/user/model/orm"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"
	"yufuture-gpt/common/utils"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AgreeApplyJoinOrgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAgreeApplyJoinOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AgreeApplyJoinOrgLogic {
	return &AgreeApplyJoinOrgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AgreeApplyJoinOrg 同意用户申请加入团队
func (l *AgreeApplyJoinOrgLogic) AgreeApplyJoinOrg(in *user.AgreeApplyJoinOrgReq) (*user.AgreeApplyJoinOrgResp, error) {
	// 找出对应的消息
	bsMessage, err := l.svcCtx.BsMessageModel.FindOne(l.ctx, in.MessageId)
	if err != nil {
		l.Logger.Error("找出对应的消息失败", err)
		return nil, err
	}
	if bsMessage.DealFlag == 1 {
		l.Logger.Error("已经处理过的消息内容", err)
		return nil, err
	}
	// 找出对应的消息内容·
	bsMessageContent, err := l.svcCtx.BsMessageContentModel.FindOne(l.ctx, bsMessage.ContentId)
	if err != nil {
		l.Logger.Error("找出对应的消息内容失败", err)
		return nil, err
	}
	// 反序列化消息内容
	var applyJoinOrgMessageContent ApplyJoinOrgMessageContent
	messageContent := bsMessageContent.MessageContent
	err = utils.StringToAny(messageContent, &applyJoinOrgMessageContent)
	if err != nil {
		l.Logger.Error("反序列化消息内容换失败", err)
		return nil, err
	}
	orgId := bsMessage.OrgId

	// 使用分布式锁 保证用户加入的团队数和团队加入的用户数得到控制
	key := applyJoinOrgMessageContent.UserId + ":" + strconv.FormatInt(orgId, 10)
	lock, err := redis.AcquireDistributedLock(l.ctx, l.svcCtx.Redis, key, 20)
	if err != nil {
		l.Logger.Error("获取分布式锁失败", err)
		return nil, err
	}
	if !lock {
		l.Logger.Error("获取分布式锁失败")
		return nil, err
	}
	defer func() {
		err = redis.ReleaseDistributedLock(l.ctx, l.svcCtx.Redis, key)
		if err != nil {
			l.Logger.Error("释放分布式锁失败", err)
		}
	}()

	// 判断用户加入的团队数量
	applyUserId, err := strconv.ParseInt(applyJoinOrgMessageContent.UserId, 10, 64)
	if err != nil {
		l.Logger.Error("转换用户id失败", err)
		return nil, err
	}
	count, err := l.svcCtx.BsUserOrgModel.FindUserOrgCount(l.ctx, applyUserId)
	if err != nil {
		l.Logger.Error("查找用户加入的团队数量失败: ", err)
		return nil, err
	}
	if count >= 2 {
		return &user.AgreeApplyJoinOrgResp{
			Code: consts.OrgNumLimit,
		}, nil
	}
	// 判断团队加入的用户数量
	org, err := l.svcCtx.BsOrganizationModel.FindOne(l.ctx, orgId)
	if err != nil {
		l.Logger.Error("查找用户申请加入的团队信息失败: ", err)
		return nil, err
	}
	// 判断团队的用户数量
	count, err = l.svcCtx.BsUserOrgModel.FindOrgUserCount(l.ctx, orgId)
	if err != nil {
		l.Logger.Error("查找团队的用户数量失败: ", err)
		return nil, err
	}
	if count >= org.MaxSeat {
		return &user.AgreeApplyJoinOrgResp{
			Code: consts.UserNumLimit,
		}, nil
	}

	// 获取MongoDB中的团队权限文档
	dborgpermission, err := l.svcCtx.DborgpermissionModel.FindOne(l.ctx, org.MongoPermId)
	if err != nil {
		l.Logger.Error("获取团队权限文档失败: ", err)
		return nil, err
	}
	// 将用户保存到MongoDB文档
	dborgpermission.UserList = append(dborgpermission.UserList, &model.User{
		Id: applyUserId,
	})
	// 更新MongoDB中的团队权限文档
	_, err = l.svcCtx.DborgpermissionModel.Update(l.ctx, dborgpermission)
	if err != nil {
		l.Logger.Error("更新团队权限文档失败: ", err)
		return nil, err
	}

	// 更新用户组织关系
	now := time.Now()
	bsUserOrg := &orm.BsUserOrg{
		UserId:     applyUserId,
		OrgId:      orgId,
		Status:     1,
		CreateTime: now,
		UpdateTime: now,
		CreateBy:   in.UserId,
		UpdateBy:   in.UserId,
	}
	_, err = l.svcCtx.BsUserOrgModel.Insert(l.ctx, bsUserOrg)
	if err != nil {
		l.Logger.Error("更新用户组织关系失败", err)
		return nil, err
	}

	// 更新消息状态
	bsMessage.DealFlag = 1
	err = l.svcCtx.BsMessageModel.Update(l.ctx, bsMessage)
	if err != nil {
		l.Logger.Error("更新消息状态失败", err)
	}

	return &user.AgreeApplyJoinOrgResp{}, nil
}
