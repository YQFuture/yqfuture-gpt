package mqs

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
)

type TrainingGoodsSuccess struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrainingGoodsSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *TrainingGoodsSuccess {
	return &TrainingGoodsSuccess{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TrainingGoodsSuccess) Consume(key, val string) error {
	logx.Infof("店铺训练商品消息消费成功 key :%s , val :%s", key, val)
	//TODO 解析消息为结构体 便于后续操作

	//TODO 消费消息，即发送商品消息给GPT进行训练，并解析返回结果

	//TODO 训练结果写入ES
	//http://192.168.3.118:9200/training_goods/_search?q=id:1
	document := struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Price int    `json:"price"`
	}{
		Id:    5,
		Name:  "Foo",
		Price: 10,
	}

	es := l.svcCtx.Elasticsearch

	res, err := es.Index().Index("training_goods").BodyJson(document).Refresh("true").Do(context.Background())
	if err != nil {
		logx.Errorf("写入ES失败, err :%s", err.Error())
		return err
	}
	logx.Infof("写入ES成功, res :%s", res)

	//默认8个消费者，所以每次消费后延迟10秒，即每个消费者每分钟消费6条数据
	time.Sleep(time.Millisecond * time.Duration(l.svcCtx.Config.TrainingGoodsConf.ConsumeDelay))
	return nil
}

type Document struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
