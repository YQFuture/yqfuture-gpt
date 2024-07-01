package basicFunction

import (
	"context"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictInfoByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDictInfoByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictInfoByTypeLogic {
	return &GetDictInfoByTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictInfoByTypeLogic) GetDictInfoByType(req *types.DictInfoByTypeReq) (resp *types.DictInfoByTypeResp, err error) {
	result, err := l.svcCtx.BasicFunctionClient.GetDictInfoByType(l.ctx, &training.DictInfoByTypeReq{
		Key: req.Key,
	})
	if err != nil {
		l.Logger.Error("根据key查询字典信息失败", err)
		return nil, err
	}

	infoList := result.List
	var dictInfoList []*types.DictInfo
	for _, dictInfo := range infoList {
		info := &types.DictInfo{
			Name:     dictInfo.Name,
			Value:    dictInfo.Value,
			OrderNum: dictInfo.OrderNum,
		}
		dictInfoList = append(dictInfoList, info)
	}
	return &types.DictInfoByTypeResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "success",
		},
		Data: dictInfoList,
	}, nil
}
