package shoptraininglogic

import (
	"context"
	"time"
	"yufuture-gpt/common/utils"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type TrainingGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTrainingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrainingGoodsLogic {
	return &TrainingGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 训练商品
func (l *TrainingGoodsLogic) TrainingGoods(in *training.TrainingGoodsReq) (*training.TrainingGoodsResp, error) {
	//根据店铺shopId查找出商品列表，需要筛选出enabled字段为1的商品
	goodsId := in.GoodsId
	goods, err := l.svcCtx.TsGoodsModel.FindOne(l.ctx, goodsId)
	if err != nil {
		return nil, err
	}
	// TODO 调用url 解析获取商品图片数组

	//TODO 将商品列表推到消息队列
	var goodsString string
	goodsString, err = utils.AnyToString(goods)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.KqPusherClient.Push(goodsString)
	if err != nil {
		l.Logger.Error("推送商品到kafka失败", goods)
		return nil, err
	}

	goods.TrainingStatus = 1
	goods.TrainingTimes += 1
	goods.UpdateTime = time.Now()
	goods.UpdateBy = in.UserId
	err = l.svcCtx.TsGoodsModel.Update(l.ctx, goods)
	if err != nil {
		l.Logger.Error("修改商品状态失败", goods)
		return nil, err
	}
	// 返回正常
	return &training.TrainingGoodsResp{}, nil
}
