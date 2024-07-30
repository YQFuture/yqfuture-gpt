package orglogic

import (
	"context"
	"strconv"
	"yufuture-gpt/app/user/model/orm"
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

	// 更新用户组织关系
	bsUserOrg := &orm.BsUserOrg{
		UserId: applyUserId,
		OrgId:  bsMessage.OrgId,
	}
	_, err = l.svcCtx.BsUserOrgModel.Insert(l.ctx, bsUserOrg)
	if err != nil {
		l.Logger.Error("更新用户组织关系失败", err)
		return nil, err
	}

	return &user.AgreeApplyJoinOrgResp{}, nil
}
