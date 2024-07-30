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

type ApplyJoinOrgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type ApplyJoinOrgMessageContent struct {
	UserId string `json:"userId"`
	Phone  string `json:"phone"`
}

func NewApplyJoinOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyJoinOrgLogic {
	return &ApplyJoinOrgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ApplyJoinOrg 用户申请加入团队
func (l *ApplyJoinOrgLogic) ApplyJoinOrg(in *user.ApplyJoinOrgReq) (*user.ApplyJoinOrgResp, error) {
	// 查找用户加入的团队数量
	count, err := l.svcCtx.BsUserOrgModel.FindUserOrgCount(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Error("查找用户加入的团队数量失败: ", err)
		return nil, err
	}
	if count >= 2 {
		return &user.ApplyJoinOrgResp{
			Code: consts.OrgNumLimit,
		}, nil
	}
	// 查找用户信息
	bsUser, err := l.svcCtx.BsUserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Error("查找用户信息失败: ", err)
		return nil, err
	}

	// 查找用户申请加入的团队信息
	org, err := l.svcCtx.BsOrganizationModel.FindOne(l.ctx, in.OrgId)
	if err != nil {
		l.Logger.Error("查找用户申请加入的团队信息失败: ", err)
		return nil, err
	}

	// 插入用户申请加入团队消息
	messageId := l.svcCtx.SnowFlakeNode.Generate().Int64()
	messageContentId := l.svcCtx.SnowFlakeNode.Generate().Int64()
	now := time.Now()
	bsMessage := &orm.BsMessage{
		Id: messageId,
		// 这里的UserId是消息接收者 也就是用户申请加入的团队管理员ID
		UserId: org.OwnerId,
		// 只有管理员在当前团队 才能看到该条消息
		OrgId:       org.Id,
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
	messageContent := ApplyJoinOrgMessageContent{
		UserId: strconv.FormatInt(bsUser.Id, 10),
		Phone:  bsUser.Phone.String,
	}
	messageContentString, err := utils.AnyToString(messageContent)
	if err != nil {
		l.Logger.Error("消息内容转字符串失败: ", err)
		return nil, err
	}
	bsMessageContent := &orm.BsMessageContent{
		Id:                 messageContentId,
		MessageType:        0,
		MessageContentType: 4,
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
		l.Logger.Error("插入用户申请加入团队消息失败: ", err)
		return nil, err
	}

	return &user.ApplyJoinOrgResp{}, nil
}
