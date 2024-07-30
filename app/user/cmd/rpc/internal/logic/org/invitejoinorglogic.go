package orglogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"time"
	"yufuture-gpt/app/user/model/orm"
	"yufuture-gpt/common/consts"
	"yufuture-gpt/common/utils"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type InviteJoinOrgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type InviteJoinOrgMessageContent struct {
	OrgId   string `json:"orgId"`
	OrgName string `json:"orgName"`
}

func NewInviteJoinOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteJoinOrgLogic {
	return &InviteJoinOrgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// InviteJoinOrg 邀请用户加入团队
func (l *InviteJoinOrgLogic) InviteJoinOrg(in *user.InviteJoinOrgReq) (*user.InviteJoinOrgResp, error) {
	// 查找用户加入的团队数量
	count, err := l.svcCtx.BsUserOrgModel.FindUserOrgCount(l.ctx, in.InviteUserId)
	if err != nil {
		l.Logger.Error("查找用户加入的团队数量失败: ", err)
		return nil, err
	}
	if count >= 2 {
		return &user.InviteJoinOrgResp{
			Code: consts.OrgNumLimit,
		}, nil
	}
	// 查找当前用户信息
	bsUser, err := l.svcCtx.BsUserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Error("查找用户信息失败: ", err)
		return nil, err
	}
	// 查找当前团队信息
	org, err := l.svcCtx.BsOrganizationModel.FindOne(l.ctx, bsUser.NowOrgId)
	if err != nil {
		l.Logger.Error("查找用户申请加入的团队信息失败: ", err)
		return nil, err
	}

	// 插入邀请用户加入团队消息
	messageId := l.svcCtx.SnowFlakeNode.Generate().Int64()
	messageContentId := l.svcCtx.SnowFlakeNode.Generate().Int64()
	now := time.Now()
	bsMessage := &orm.BsMessage{
		Id: messageId,
		// 这里的UserId是消息接收者 也就是被邀请加入团队的用户ID
		UserId:      in.InviteUserId,
		OrgId:       0,
		MessageType: 0,
		ContentId:   messageContentId,
		ReadFlag:    0,
		DealFlag:    0,
		IgnoreFlag:  0,
		CreateTime:  now,
		UpdateTime:  now,
		CreateBy:    in.UserId,
		UpdateBy:    in.UserId,
	}
	// 构建消息内容
	messageContent := InviteJoinOrgMessageContent{
		OrgId:   strconv.FormatInt(org.Id, 10),
		OrgName: org.OrgName.String,
	}
	messageContentString, err := utils.AnyToString(messageContent)
	if err != nil {
		l.Logger.Error("消息内容转字符串失败: ", err)
		return nil, err
	}
	bsMessageContent := &orm.BsMessageContent{
		Id:                 messageContentId,
		MessageType:        0,
		MessageContentType: 3,
		MessageContent:     messageContentString,
		CreateTime:         now,
		UpdateTime:         now,
		CreateBy:           in.UserId,
		UpdateBy:           in.UserId,
	}

	err = l.svcCtx.BsMessageModel.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.BsMessageModel.SessionInsert(l.ctx, bsMessage, session)
		if err != nil {
			return err
		}
		_, err = l.svcCtx.BsMessageContentModel.SessionInsert(l.ctx, bsMessageContent, session)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		l.Logger.Error("插入邀请用户加入团队消息失败: ", err)
		return nil, err
	}

	return &user.InviteJoinOrgResp{}, nil
}
