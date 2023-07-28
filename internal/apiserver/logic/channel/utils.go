package channel

import (
	"context"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
)

func verifyCreateNameUnique(name string, svcCtx *svc.ServiceContext, ctx context.Context) error {
	var (
		cT = query.Use(svcCtx.Db).Channel
	)

	count, err := cT.WithContext(ctx).
		Where(cT.Name.Eq(name),
			cT.Abandoned.Is(false)).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return code.WithCodeMsg(code.DuplicateChannelName)
	}

	return nil
}

func verifyUpdateNameUnique(name string, id int32, svcCtx *svc.ServiceContext, ctx context.Context) error {
	var (
		cT = query.Use(svcCtx.Db).Channel
	)

	count, err := cT.WithContext(ctx).
		Where(cT.Name.Eq(name),
			cT.Abandoned.Is(false),
			cT.ID.Neq(id)).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return code.WithCodeMsg(code.DuplicateChannelName)
	}

	return nil
}
