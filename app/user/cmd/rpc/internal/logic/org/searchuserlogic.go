package orglogic

import (
	"context"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SearchUser 查找用户
func (l *SearchUserLogic) SearchUser(in *user.SearchUserReq) (*user.SearchUserResp, error) {
	userList, err := l.svcCtx.BsUserModel.FindListByPhone(l.ctx, in.Query)
	if err != nil {
		l.Logger.Error("根据用手机号查找用户失败", err)
		return nil, err
	}

	var userInfoList []*user.SearchUserInfo
	for _, userInfo := range *userList {
		userInfoList = append(userInfoList, &user.SearchUserInfo{
			UserId: userInfo.Id,
			Phone:  userInfo.Phone.String,
		})
	}
	return &user.SearchUserResp{
		Result: userInfoList,
	}, nil
}
