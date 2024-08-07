package org

import (
	"context"
	"encoding/json"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
	"time"
	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"
	"yufuture-gpt/app/user/cmd/rpc/client/org"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgUserOperationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetOrgUserOperationListLogic 导出团队用户操作记录列表
func NewGetOrgUserOperationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgUserOperationListLogic {
	return &GetOrgUserOperationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrgUserOperationListLogic) GetOrgUserOperationList(req *types.OrgUserOperationListReq, w http.ResponseWriter) error {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return err

	}
	operationUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		l.Logger.Error("获取要查询角色的用户ID失败", err)
		return err
	}

	// 调用RPC接口 获取团队用户操作记录列表
	operationList, err := l.svcCtx.OrgClient.GetOrgUserOperationList(l.ctx, &org.OrgUserOperationListReq{
		UserId:          userId,
		OperationUserId: operationUserId,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Query:           req.Query,
	})
	if err != nil {
		l.Logger.Error("获取团队用户操作记录列表失败", err)
		return err
	}

	// 创建一个新的Excel文件
	f := excelize.NewFile()

	// 添加工作表
	_, err = f.NewSheet("Sheet1")
	if err != nil {
		l.Logger.Error("创建工作表失败", err)
		return err
	}

	// 设置列名
	err = f.SetCellValue("Sheet1", "A1", "时间")
	if err != nil {
		l.Logger.Error("创建工作表失败", err)
		return err
	}
	err = f.SetCellValue("Sheet1", "B1", "操作")
	if err != nil {
		l.Logger.Error("创建工作表失败", err)
		return err
	}

	// 填充数据
	for i, record := range operationList.List {
		// 假设 record.CreateTime 是 int64 类型
		err = f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), time.Unix(record.CreateTime, 0).Format("2006-01-02 15:04:05"))
		if err != nil {
			l.Logger.Error("填充数据失败", err)
			return err
		}
		err = f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), record.OperationDesc)
		if err != nil {
			l.Logger.Error("填充数据失败", err)
			return err
		}
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=export.xlsx")

	// 将Excel文件写入HTTP响应
	err = f.Write(w)
	if err != nil {
		l.Logger.Error("返回工作表失败", err)
		return err
	}

	return nil
}
