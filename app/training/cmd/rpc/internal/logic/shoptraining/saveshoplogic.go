package shoptraininglogic

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	yqmongo "yufuture-gpt/app/training/model/mongo"
	"yufuture-gpt/app/training/model/orm"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveShopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveShopLogic {
	return &SaveShopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SaveShop 保存爬取的店铺基本数据
func (l *SaveShopLogic) SaveShop(in *training.SaveShopReq) (*training.SaveShopResp, error) {
	// 根据userid和uuid查找出店铺 新老店铺执行不同流程的操作
	tsShop, err := l.svcCtx.TsShopModel.FindOneByUuidAndUserId(l.ctx, in.UserId, in.Uuid)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找店铺失败", err)
		return nil, err
	}

	// 记录本次爬取的商品数据 用于保存到mongo
	var mongoGoodsList []*orm.TsGoods
	// 对于新店铺 执行新增操作步骤
	if tsShop == nil || tsShop.Id == 0 {
		// 解析数据保存到mysql
		tsShop = buildTsShop(l, in)
		_, err := l.svcCtx.TsShopModel.Insert(l.ctx, tsShop)
		if err != nil {
			l.Logger.Error("保存店铺到mysql失败", err, in)
			return nil, err
		}

		for _, saveGoods := range in.List {
			tsGoods := buildTsGoods(l, in, tsShop, saveGoods)
			mongoGoodsList = append(mongoGoodsList, tsGoods)
			_, err = l.svcCtx.TsGoodsModel.Insert(l.ctx, tsGoods)
			if err != nil {
				l.Logger.Error("保存商品到mysql失败", err, in)
				return nil, err
			}
		}
	} else {
		// 查询老店铺原有商品列表并暂存map
		tsGoodsList, err := l.svcCtx.TsGoodsModel.FindListByShopId(l.ctx, tsShop.Id)
		if err != nil {
			l.Logger.Error("查询老店铺商品列表失败", err, in)
			return nil, err
		}
		tsGoodsMap := make(map[string]*orm.TsGoods)
		if tsGoodsList != nil {
			for _, tsGoods := range *tsGoodsList {
				tsGoodsMap[tsGoods.PlatformId] = tsGoods
			}
		}

		for _, saveGoods := range in.List {
			if tsGoods, ok := tsGoodsMap[saveGoods.PlatformId]; !ok {
				//对于新增的商品，保存到mysql
				tsGoods := buildTsGoods(l, in, tsShop, saveGoods)
				mongoGoodsList = append(mongoGoodsList, tsGoods)
				_, err = l.svcCtx.TsGoodsModel.Insert(l.ctx, tsGoods)
				if err != nil {
					l.Logger.Error("保存商品到mysql失败", err, in)
				}
			} else {
				//老商品无需操作
				mongoGoodsList = append(mongoGoodsList, tsGoods)
			}
		}

	}

	// 不论新老店铺都统一根据本次信息保存一条新记录到mongo
	var saveGoodsList []*yqmongo.GoodsIdList
	for _, tsGoods := range mongoGoodsList {
		saveGoodsList = append(saveGoodsList, &yqmongo.GoodsIdList{
			GoodsId:    tsGoods.Id,
			PlatformId: tsGoods.PlatformId,
			Url:        tsGoods.GoodsUrl,
		})
	}
	dbsavegoodscrawlertitles := buildDbsavegoodscrawlertitles(tsShop, in, saveGoodsList)
	err = l.svcCtx.DbsavegoodscrawlertitlesModel.Insert(l.ctx, dbsavegoodscrawlertitles)
	if err != nil {
		l.Logger.Error("保存店铺到mongo失败", err)
		return nil, err
	}

	return &training.SaveShopResp{}, nil
}

func buildTsShop(l *SaveShopLogic, in *training.SaveShopReq) *orm.TsShop {
	return &orm.TsShop{
		Id:             l.svcCtx.SnowFlakeNode.Generate().Int64(),
		PlatformId:     "", //平台店铺ID 暂时留空
		ShopName:       in.ShopName,
		UserId:         in.UserId,
		Uuid:           in.Uuid,
		GroupId:        0, // 部门ID 为之后的版本预留 暂时填零
		PlatformType:   in.PlatformType,
		TrainingStatus: 0,
		TrainingTimes:  0,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		CreateBy:       in.UserId,
		UpdateBy:       in.UserId,
	}
}

func buildTsGoods(l *SaveShopLogic, in *training.SaveShopReq, tsShop *orm.TsShop, saveGoods *training.SaveGoods) *orm.TsGoods {
	return &orm.TsGoods{
		Id:              l.svcCtx.SnowFlakeNode.Generate().Int64(),
		ShopId:          tsShop.Id,
		PlatformId:      saveGoods.PlatformId,
		GoodsName:       saveGoods.GoodsName,
		GoodsUrl:        saveGoods.GoodsUrl,
		GoodsJsonUrl:    "",
		TrainingSummary: "",
		PlatformType:    in.PlatformType,
		Enabled:         2, // 默认开启
		TrainingStatus:  0,
		TrainingTimes:   0,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		CreateBy:        in.UserId,
		UpdateBy:        in.UserId,
	}
}

func buildDbsavegoodscrawlertitles(tsShop *orm.TsShop, in *training.SaveShopReq, saveGoodsList []*yqmongo.GoodsIdList) *yqmongo.Dbsavegoodscrawlertitles {
	return &yqmongo.Dbsavegoodscrawlertitles{
		ID:       primitive.NewObjectID(),
		CreateAt: time.Now(),
		UpdateAt: time.Now(),

		SerialNumber: in.SerialNumber,
		ShopId:       tsShop.Id,
		ShopName:     in.ShopName,
		UserID:       in.UserId,
		UUID:         in.Uuid,
		Platform:     in.PlatformType,
		GoodsList:    saveGoodsList,
	}
}
