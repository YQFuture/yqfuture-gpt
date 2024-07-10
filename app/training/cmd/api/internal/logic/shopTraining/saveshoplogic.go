package shopTraining

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 保存爬取的店铺基本数据
func NewSaveShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveShopLogic {
	return &SaveShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveShopLogic) SaveShop(req *types.SaveShopReq) (resp *types.SaveShopResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	// 构建保存消息体
	var saveGoodsList []*training.SaveGoods
	for _, v := range req.GoodsList {
		saveGoodsList = append(saveGoodsList, &training.SaveGoods{
			GoodsName:  v.GoodsName,
			GoodsUrl:   v.GoodsUrl,
			PlatformId: v.PlatFormId,
		})
	}
	_, err = l.svcCtx.ShopTrainingClient.SaveShop(l.ctx, &training.SaveShopReq{
		UserId:       userId,
		Uuid:         req.Uuid,
		ShopName:     req.ShopName,
		PlatformType: 2,
		List:         saveGoodsList,
	})
	if err != nil {
		l.Logger.Error("保存店铺基本信息失败", err)
		return nil, err
	}
	return &types.SaveShopResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "保存店铺基本信息成功",
		},
	}, nil
}
