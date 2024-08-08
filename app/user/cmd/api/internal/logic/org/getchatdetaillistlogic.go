package org

import (
	"context"
	"time"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatDetailListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetChatDetailListLogic 获取聊天记录列表
func NewGetChatDetailListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatDetailListLogic {
	return &GetChatDetailListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatDetailListLogic) GetChatDetailList(req *types.ChatDetailReq) (resp *types.ChatDetailResp, err error) {

	return &types.ChatDetailResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: []types.ChatDetail{
			{
				OwnerId:     "1",
				SenderType:  2,
				ContentType: 1,
				Content:     "你好",
				CreateTime:  time.Now().Unix(),
			},
			{
				OwnerId:     "1818638829322506240",
				SenderType:  0,
				ContentType: 1,
				Content:     "您好 客服小云为您服务",
				CreateTime:  time.Now().Unix(),
			},
			{
				OwnerId:     "1818638829322506240",
				SenderType:  1,
				ContentType: 1,
				Content:     "您好 人工客服已上线",
				CreateTime:  time.Now().Unix(),
			},
			{
				OwnerId:     "1",
				SenderType:  2,
				ContentType: 1,
				Content:     "这瓜保熟吗",
				CreateTime:  time.Now().Unix(),
			},
			{
				OwnerId:     "1818638829322506240",
				SenderType:  0,
				ContentType: 1,
				Content:     "我们开水果店的",
				CreateTime:  time.Now().Unix(),
			},
			{
				OwnerId:     "1818638829322506240",
				SenderType:  0,
				ContentType: 1,
				Content:     "能卖你生瓜蛋子吗",
				CreateTime:  time.Now().Unix(),
			},
		},
	}, nil
}
