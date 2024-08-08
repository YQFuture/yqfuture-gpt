package orglogic

import (
	"context"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShopUserBuyerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetShopUserBuyerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopUserBuyerListLogic {
	return &GetShopUserBuyerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取店铺客服买家列表
func (l *GetShopUserBuyerListLogic) GetShopUserBuyerList(in *user.ShopUserBuyerListReq) (*user.ShopUserBuyerListResp, error) {
	// todo: add your logic here and delete this line

	return &user.ShopUserBuyerListResp{}, nil
}
