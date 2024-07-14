package es

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

const GoodsTrainingResult = "goods_training_result"

// SaveGoodsTrainingResult 商品训练结果写入ES 店铺训练结果作为多个商品训练结果的集合同样写入该索引
func SaveGoodsTrainingResult(ctx context.Context, es *elastic.Client, document interface{}) (*elastic.IndexResponse, error) {
	res, err := es.Index().Index(GoodsTrainingResult).BodyJson(document).Refresh("true").Do(ctx)
	if err != nil {
		logx.Errorf("商品训练结果写入ES失败, err :%v", err)
		return nil, err
	}
	logx.Infof("商品训练结果写入ES成功, res :%v", res)
	return res, nil
}

// SearchGoodsTrainingResults 在ES中搜索商品训练结果
func SearchGoodsTrainingResults(ctx context.Context, es *elastic.Client, query elastic.Query) (*elastic.SearchResult, error) {
	searchResult, err := es.Search().
		Index(GoodsTrainingResult).
		Query(query).
		Do(ctx)
	if err != nil {
		logx.Errorf("在ES中搜索商品训练结果失败, query :%v, err :%v", query, err)
		return nil, err
	}
	logx.Infof("在ES中搜索商品训练结果成功, query :%v, found %d hits", query, searchResult.Hits.TotalHits.Value)
	return searchResult, nil
}
