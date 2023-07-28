// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:             db,
		AccessLog:      newAccessLog(db, opts...),
		ApisixRoute:    newApisixRoute(db, opts...),
		ApisixUpstream: newApisixUpstream(db, opts...),
		Application:    newApplication(db, opts...),
		Channel:        newChannel(db, opts...),
		Config:         newConfig(db, opts...),
		Document:       newDocument(db, opts...),
		Optlog:         newOptlog(db, opts...),
		Service:        newService(db, opts...),
		User:           newUser(db, opts...),
		Version:        newVersion(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	AccessLog      accessLog
	ApisixRoute    apisixRoute
	ApisixUpstream apisixUpstream
	Application    application
	Channel        channel
	Config         config
	Document       document
	Optlog         optlog
	Service        service
	User           user
	Version        version
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:             db,
		AccessLog:      q.AccessLog.clone(db),
		ApisixRoute:    q.ApisixRoute.clone(db),
		ApisixUpstream: q.ApisixUpstream.clone(db),
		Application:    q.Application.clone(db),
		Channel:        q.Channel.clone(db),
		Config:         q.Config.clone(db),
		Document:       q.Document.clone(db),
		Optlog:         q.Optlog.clone(db),
		Service:        q.Service.clone(db),
		User:           q.User.clone(db),
		Version:        q.Version.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:             db,
		AccessLog:      q.AccessLog.replaceDB(db),
		ApisixRoute:    q.ApisixRoute.replaceDB(db),
		ApisixUpstream: q.ApisixUpstream.replaceDB(db),
		Application:    q.Application.replaceDB(db),
		Channel:        q.Channel.replaceDB(db),
		Config:         q.Config.replaceDB(db),
		Document:       q.Document.replaceDB(db),
		Optlog:         q.Optlog.replaceDB(db),
		Service:        q.Service.replaceDB(db),
		User:           q.User.replaceDB(db),
		Version:        q.Version.replaceDB(db),
	}
}

type queryCtx struct {
	AccessLog      *accessLogDo
	ApisixRoute    *apisixRouteDo
	ApisixUpstream *apisixUpstreamDo
	Application    *applicationDo
	Channel        *channelDo
	Config         *configDo
	Document       *documentDo
	Optlog         *optlogDo
	Service        *serviceDo
	User           *userDo
	Version        *versionDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		AccessLog:      q.AccessLog.WithContext(ctx),
		ApisixRoute:    q.ApisixRoute.WithContext(ctx),
		ApisixUpstream: q.ApisixUpstream.WithContext(ctx),
		Application:    q.Application.WithContext(ctx),
		Channel:        q.Channel.WithContext(ctx),
		Config:         q.Config.WithContext(ctx),
		Document:       q.Document.WithContext(ctx),
		Optlog:         q.Optlog.WithContext(ctx),
		Service:        q.Service.WithContext(ctx),
		User:           q.User.WithContext(ctx),
		Version:        q.Version.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
