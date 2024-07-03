package mqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/model/orm"
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
	var tsGoods orm.TsGoods
	err := json.Unmarshal([]byte(val), &tsGoods)
	if err != nil {
		logx.Errorf("解析训练商品消息失败, err :%s", err.Error())
		return err
	}
	//TODO 消费消息，即发送商品消息给GPT进行训练，并解析返回结果

	//TODO 训练结果写入ES
	//http://192.168.3.118:9200/training_goods/_search?q=id:1
	//http://192.168.3.118:9200/training_goods/_doc/iURcdpABA6PejL7zs03J
	document := struct {
		Id             int64     `json:"id"`
		TrainingResult string    `json:"training_result"`
		CreateTime     time.Time `json:"create_time"`
	}{
		Id:             tsGoods.Id,
		TrainingResult: "作为一个AI模型，我的训练结果通常是以数字形式存储在大量的参数中，这些参数编码了我从数据中学到的知识。由于这些参数通常是非常大的数值数组，它们并不容易被直接理解。不过，我可以提供一个简化的例子来说明AI模型训练结果的含义。\n\n假设我们有一个简单的线性回归模型，它的目的是学习一个函数来预测基于单个输入特征的输出值。模型的训练结果就是找到最佳的权重（w）和偏差（b），使得模型的预测尽可能接近实际的数据点。\n\n例如，如果我们有一系列的数据点 {(x1, y1), (x2, y2), ..., (xn, yn)}，其中 x 是输入特征，y 是对应的目标值，那么训练过程可能会得到如下的权重和偏差：\n\n权重 w = 2.5\n偏差 b = 0.7\n\n这意味着我们的模型学到的函数是 f(x) = 2.5x + 0.7。对于一个新的输入值 x_new，模型会使用这个函数来预测输出值：\n\n预测的 y_new = f(x_new) = 2.5 * x_new + 0.7\n\n这个简化的例子展示了即使是非常基础的AI模型，其训练结果也是通过一组数值（在这个例子中是权重和偏差）来表达的，这些数值捕获了输入数据和目标输出之间的关系。\n\n对于更复杂的模型，比如深度学习模型，训练结果会包含大量的权重和偏差参数，这些参数构成了模型学到的复杂函数的表示。尽管这些参数本身可能不容易解释，但它们共同定义了模型如何处理和预测新的数据。",
		CreateTime:     time.Now(),
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
