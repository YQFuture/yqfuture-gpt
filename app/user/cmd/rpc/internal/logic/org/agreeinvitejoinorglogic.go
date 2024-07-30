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

type AgreeInviteJoinOrgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAgreeInviteJoinOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AgreeInviteJoinOrgLogic {
	return &AgreeInviteJoinOrgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AgreeInviteJoinOrg 同意用户邀请加入团队
func (l *AgreeInviteJoinOrgLogic) AgreeInviteJoinOrg(in *user.AgreeInviteJoinOrgReq) (*user.AgreeInviteJoinOrgResp, error) {
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
	var inviteJoinOrgMessageContent InviteJoinOrgMessageContent
	messageContent := bsMessageContent.MessageContent
	err = utils.StringToAny(messageContent, &inviteJoinOrgMessageContent)
	if err != nil {
		l.Logger.Error("反序列化消息内容换失败", err)
		return nil, err
	}

	// 判断用户加入的团队数量
	count, err := l.svcCtx.BsUserOrgModel.FindUserOrgCount(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Error("查找用户加入的团队数量失败: ", err)
		return nil, err
	}
	if count >= 2 {
		return &user.AgreeInviteJoinOrgResp{
			Code: consts.OrgNumLimit,
		}, nil
	}

	// 更新用户组织关系
	inviteOrgId, err := strconv.ParseInt(inviteJoinOrgMessageContent.OrgId, 10, 64)
	if err != nil {
		l.Logger.Error("转换组织id失败", err)
		return nil, err
	}
	bsUserOrg := &orm.BsUserOrg{
		UserId: in.UserId,
		OrgId:  inviteOrgId,
	}
	_, err = l.svcCtx.BsUserOrgModel.Insert(l.ctx, bsUserOrg)
	if err != nil {
		l.Logger.Error("更新用户组织关系失败", err)
		return nil, err
	}

	return &user.AgreeInviteJoinOrgResp{}, nil
}
