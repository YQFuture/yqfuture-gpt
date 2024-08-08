package org

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"yufuture-gpt/app/user/cmd/rpc/client/org"
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
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.ChatDetailResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	shopId, err := strconv.ParseInt(req.ShopId, 10, 64)
	if err != nil {
		l.Logger.Error("获取店铺ID失败", err)
		return &types.ChatDetailResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	shopUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		l.Logger.Error("获取客服用户ID失败", err)
		return &types.ChatDetailResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	buyerId, err := strconv.ParseInt(req.BuyerId, 10, 64)
	if err != nil {
		l.Logger.Error("获取客服用户ID失败", err)
		return &types.ChatDetailResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 获取聊天记录
	chatDetailResp, err := l.svcCtx.OrgClient.GetChatDetail(l.ctx, &org.ChatDetailReq{
		UserId:     userId,
		ShopId:     shopId,
		ShopUserId: shopUserId,
		BuyerId:    buyerId,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Query:      req.Query,
	})
	if err != nil {
		l.Logger.Error("获取聊天记录失败", err)
		return &types.ChatDetailResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	var data []types.ChatDetail
	for _, v := range chatDetailResp.List {
		data = append(data, types.ChatDetail{
			OwnerId:     strconv.FormatInt(v.OwnerId, 10),
			SenderType:  v.SenderType,
			ContentType: v.ContentType,
			Content:     v.Content,
			CreateTime:  v.CreateTime,
		})
	}
	/*	return &types.ChatDetailResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: data,
	}, nil*/

	return returnMockChatDetailResp()
}

func returnMockChatDetailResp() (resp *types.ChatDetailResp, err error) {
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
