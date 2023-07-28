package user

import (
	"context"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/base"
)

func getChannelList(channelIDList []int32, svcCtx *svc.ServiceContext, ctx context.Context) []base.BaseInfo {

	var (
		cT      = query.Use(svcCtx.Db).Channel
		chnList = make([]base.BaseInfo, 0)
	)

	query := cT.WithContext(ctx).Select(cT.ID, cT.Name).Where(cT.ID.In(channelIDList...))
	items, _ := query.Debug().Find()

	if len(items) > 0 {
		for _, item := range items {
			chnList = append(chnList, base.BaseInfo{
				ID:   int(item.ID),
				Name: item.Name,
			})
		}
	}

	return chnList
}

func verifyCreateAccountUniq(account string, svcCtx *svc.ServiceContext, ctx context.Context) error {
	var (
		uT = query.Use(svcCtx.Db).User
	)

	count, err := uT.WithContext(ctx).
		Where(uT.Account.Eq(account),
			uT.Abandoned.Is(false)).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return code.WithCodeMsg(code.DuplicateAccount)
	}

	return nil
}

func verifyUpdateAccountUniq(account string, ID int32, svcCtx *svc.ServiceContext, ctx context.Context) error {
	var (
		uT = query.Use(svcCtx.Db).User
	)

	count, err := uT.WithContext(ctx).
		Where(uT.Account.Eq(account),
			uT.ID.Neq(ID),
			uT.Abandoned.Is(false)).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return code.WithCodeMsg(code.DuplicateAccount)
	}

	return nil
}

func verifyChannel(channelIDList []int32, userID int32, svcCtx *svc.ServiceContext, ctx context.Context) error {
	var (
		cT = query.Use(svcCtx.Db).Channel
	)

	count, err := cT.WithContext(ctx).
		Where(cT.ID.In(channelIDList...)).Count()

	if err != nil {
		return err
	}

	if int(count) != len(channelIDList) {
		return code.WithCodeMsg(code.ChannelNotExist)
	}

	return nil

}
