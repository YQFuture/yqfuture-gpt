package shopTraining

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSettingGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 预设商品
func NewPreSettingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingGoodsLogic {
	return &PreSettingGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreSettingGoodsLogic) PreSettingGoods(req *types.PreSettingGoodsReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return nil, err
	}
	goodsIdInt, err := strconv.ParseInt(req.GoodsId, 10, 64)
	if err != nil {
		l.Logger.Error("转换商品ID失败", err)
		return nil, err
	}
	_, err = l.svcCtx.ShopTrainingClient.PreSettingGoods(l.ctx, &training.PreSettingGoodsReq{
		GoodsId:       goodsIdInt,
		UserId:        userId,
		Authorization: req.Authorization,
		Cookies:       req.Cookies,
	})
	// 后台会进行较长时间的轮询等待 所以会返回超时错误 无需处理 仅打印日志 前端关注店铺和商品的训练状态即可
	if err != nil {
		l.Logger.Error("开启预设商品异常", err)
	}
	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "预设商品成功",
	}, nil
}
