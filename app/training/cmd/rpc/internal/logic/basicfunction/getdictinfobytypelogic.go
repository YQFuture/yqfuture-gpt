package basicfunctionlogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictInfoByTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictInfoByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictInfoByTypeLogic {
	return &GetDictInfoByTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDictInfoByTypeLogic) GetDictInfoByType(in *training.DictInfoByTypeReq) (*training.DictInfoByTypeResp, error) {
	dictType, err := l.svcCtx.BsDictTypeModel.FindOneByKey(l.ctx, in.Key)
	if err != nil {
		l.Logger.Error("根据key查询字典类型失败", err)
		return nil, err
	}
	dictInfoList, err := l.svcCtx.BsDictInfoModel.FindListByTypeId(l.ctx, dictType.Id)
	if err != nil {
		l.Logger.Error("根据typeId查询字典信息失败", err)
		return nil, err
	}
	var infoList []*training.DictInfo
	for _, dictInfo := range *dictInfoList {
		info := &training.DictInfo{
			Name:     dictInfo.Name,
			Value:    dictInfo.Value.String,
			OrderNum: dictInfo.OrderNum,
		}
		infoList = append(infoList, info)
	}
	return &training.DictInfoByTypeResp{
		List: infoList,
	}, nil
}
