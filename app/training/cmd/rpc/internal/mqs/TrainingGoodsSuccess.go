package mqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/utills"
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
	// 解析消息为结构体 便于后续操作
	var tsGoods orm.TsGoods
	err := json.Unmarshal([]byte(val), &tsGoods)
	if err != nil {
		logx.Errorf("解析训练商品消息失败, err :%s", err.Error())
		return err
	}

	//TODO 从mongo中获取商品JSON 获取图片列表

	//TODO 消费消息，即发送商品消息给GPT进行训练，并解析返回结果
	request := &TrainingRequest{}
	result, err := trainingGoods(l.svcCtx.Config.GptImageURL, request)
	if err != nil {
		return err
	}
	response := result.Data.Response
	token := result.Data.Token
	logx.Infof("训练结果：%s", response)
	logx.Infof("消耗的token：%d", token)

	// 训练结果写入ES
	document := &Document{
		Id:             tsGoods.Id,
		TrainingResult: response,
		TOKEN:          token,
		CreateTime:     time.Now(),
	}
	es := l.svcCtx.Elasticsearch
	res, err := es.Index().Index("training_goods").Id(strconv.FormatInt(tsGoods.Id, 10)).BodyJson(document).Refresh("true").Do(context.Background())
	if err != nil {
		logx.Errorf("写入ES失败, err :%s", err.Error())
		return err
	}
	logx.Infof("写入ES成功, res :%s", res)

	//训练日志写入ES
	go func() {
		goodsTrainingLog := &GoodsTrainingLog{
			Id:             tsGoods.Id,
			TrainingResult: response,
			TOKEN:          token,
			CreateTime:     time.Now(),
		}
		res, err := es.Index().Index("goods_training_log").BodyJson(goodsTrainingLog).Refresh("true").Do(context.Background())
		if err != nil {
			logx.Errorf("写入ES失败, err :%s", err.Error())
		}
		logx.Infof("写入ES成功, res :%s", res)
	}()

	//更新tsGoods在MySql中的字段
	tsGoods.TrainingSummary = response[:20]
	err = l.svcCtx.TsGoodsModel.Update(l.ctx, &tsGoods)
	if err != nil {
		logx.Errorf("更新tsGoods失败, res :%s", err)
	}

	//默认8个消费者，所以每次消费后延迟10秒，即每个消费者每分钟消费6条数据
	time.Sleep(time.Millisecond * time.Duration(l.svcCtx.Config.TrainingGoodsConf.ConsumeDelay))
	return nil
}

type Document struct {
	Id             int64     `json:"id"`
	TrainingResult string    `json:"training_result"`
	TOKEN          int       `json:"token"`
	CreateTime     time.Time `json:"create_time"`
}

type GoodsTrainingLog struct {
	Id             int64     `json:"id"`
	TrainingResult string    `json:"training_result"`
	TOKEN          int       `json:"token"`
	CreateTime     time.Time `json:"create_time"`
}

type TrainingRequest struct {
	//SystemPrompt string   `json:"system_prompt"` //可选，传入会替代默认的prompt
	ImageURLs []string `json:"image_urls"`
}

type TrainingResult struct {
	Status bool `json:"status"`
	Data   struct {
		Response string `json:"response"` // 训练结果
		Token    int    `json:"token"`    // 消耗的token
	} `json:"data"`
	Msg string `json:"msg"`
}

func trainingGoods(url string, request *TrainingRequest) (TrainingResult, error) {
	var result TrainingResult
	err := utills.HTTPPostAndParseJSON(url, request, &result)
	if err != nil {
		logx.Errorf("训练商品失败, err :%s", err.Error())
		return result, err
	}
	return result, nil
}
