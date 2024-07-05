package shoptraininglogic

import (
	"context"
	"github.com/olivere/elastic/v7"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGoodsTrainingResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGoodsTrainingResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsTrainingResultLogic {
	return &GetGoodsTrainingResultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGoodsTrainingResult 获取商品训练结果
func (l *GetGoodsTrainingResultLogic) GetGoodsTrainingResult(in *training.GetGoodsTrainingResultReq) (*training.GetGoodsTrainingResultResp, error) {
	es := l.svcCtx.Elasticsearch
	query := elastic.NewTermQuery("id", in.GoodsId)
	result, err := es.Search().Index("training_goods").Query(query).Size(1).Do(l.ctx)
	if err != nil {
		l.Logger.Error("从ES获取商品训练结果失败", err)
		return nil, err
	}
	//序列化结果并返回
	var trainingResult string
	trainingResult, err = utils.AnyToString(result.Hits.Hits[0].Source)
	if err != nil {
		l.Logger.Error("序列化商品训练结果失败", err)
		return nil, err
	}
	return &training.GetGoodsTrainingResultResp{
		Result: trainingResult,
	}, nil
}
