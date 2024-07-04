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
	logx.Infof("开始消费店铺训练商品消息 key :%s , val :%s", key, val)
	// 解析消息为结构体 便于后续操作
	var tsGoods orm.TsGoods
	err := json.Unmarshal([]byte(val), &tsGoods)
	if err != nil {
		logx.Errorf("解析训练商品消息失败, err :%s", err.Error())
		return err
	}
	//TODO 从消息中直接解析训练请求对象

	request := &TrainingRequest{}

	// 消费消息，即发送商品消息给GPT进行训练，并解析返回结果
	result := trainingGoods(l.svcCtx.Config.GptImageURL, request)
	var response string
	var token int
	if result.Status == false {
		logx.Errorf("训练失败, err :%s, 商品信息：%v", result.Msg, tsGoods)
		response = "训练失败"
		token = 0
	} else {
		response = result.Data.Response
		token = result.Data.Token
	}

	// 训练结果写入ES 这里的训练结果只从接口返回值中提取出的文字描述
	document := &Document{
		Id:             tsGoods.Id,
		TrainingResult: response,
		TOKEN:          token,
		CreateTime:     time.Now(),
	}
	es := l.svcCtx.Elasticsearch
	res, err := es.Index().Index("training_goods").Id(strconv.FormatInt(tsGoods.Id, 10)).BodyJson(document).Refresh("true").Do(context.Background())
	if err != nil {
		logx.Errorf("训练结果写入ES失败, err :%s", err.Error())
		return err
	}
	logx.Infof("训练结果写入ES成功, res :%v, result :%v", res, result)

	// 训练日志写入ES 这里的训练结果会保存整个的接口返回值
	go func() {
		goodsTrainingLog := &GoodsTrainingLog{
			Id:             tsGoods.Id,
			TrainingResult: result,
			TOKEN:          token,
			CreateTime:     time.Now(),
		}
		res, err := es.Index().Index("goods_training_log").BodyJson(goodsTrainingLog).Refresh("true").Do(context.Background())
		if err != nil {
			logx.Errorf("训练日志写入ES失败, err :%s", err.Error())
		}
		logx.Infof("训练日志写入ES成功, res :%v result :%v", res, result)
	}()

	// 更新tsGoods在MySQL中的字段
	if len(response) > 20 {
		tsGoods.TrainingSummary = response[:20]
	} else {
		tsGoods.TrainingSummary = response
	}
	tsGoods.TrainingStatus = 2
	err = l.svcCtx.TsGoodsModel.Update(l.ctx, &tsGoods)
	if err != nil {
		logx.Errorf("更新tsGoods失败, res :%s", err)
	}

	// 默认8个消费者, 所以每次消费后延迟指定时间, 以此控制消费频率, 可通过配置文件配置
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
	Id             int64          `json:"id"`
	TrainingResult TrainingResult `json:"training_result"`
	TOKEN          int            `json:"token"`
	CreateTime     time.Time      `json:"create_time"`
}

type TrainingRequest struct {
	//SystemPrompt string   `json:"system_prompt"` // 可选，传入会替代默认的prompt
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

func trainingGoods(url string, request *TrainingRequest) TrainingResult {
	var result TrainingResult
	err := utills.HTTPPostAndParseJSON(url, request, &result)
	if err != nil {
		logx.Errorf("训练商品失败, err :%s", err.Error())
		return TrainingResult{
			Status: false,
			Msg:    err.Error(),
		}
	}
	return result
}
