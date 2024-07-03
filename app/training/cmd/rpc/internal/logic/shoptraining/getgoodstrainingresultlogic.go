package shoptraininglogic

import (
	"context"
	"github.com/olivere/elastic/v7"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/utills"

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
	// TODO 根据商品ID从ES中获取训练结果

	es := l.svcCtx.Elasticsearch
	query := elastic.NewTermQuery("id", in.GoodsId) // 替换"目标标题值"为你要搜索的实际值
	result, err := es.Search().Index("training_goods").Query(query).Size(1).Do(l.ctx)
	if err != nil {
		return nil, err
	}
	var trainingResult string
	// 处理查询结果
	trainingResult, err = utills.AnyToString(result.Hits.Hits[0].Source)

	// TODO 序列化结果并返回
	return &training.GetGoodsTrainingResultResp{
		Result: trainingResult,
	}, nil
}
