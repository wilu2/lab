// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
)

func newAccessLog(db *gorm.DB, opts ...gen.DOOption) accessLog {
	_accessLog := accessLog{}

	_accessLog.accessLogDo.UseDB(db, opts...)
	_accessLog.accessLogDo.UseModel(&model.AccessLog{})

	tableName := _accessLog.accessLogDo.TableName()
	_accessLog.ALL = field.NewAsterisk(tableName)
	_accessLog.ID = field.NewInt32(tableName, "id")
	_accessLog.RouteID = field.NewString(tableName, "route_id")
	_accessLog.RequestID = field.NewString(tableName, "request_id")
	_accessLog.ClientAddr = field.NewString(tableName, "client_addr")
	_accessLog.IsoTime = field.NewString(tableName, "iso_time")
	_accessLog.Timestamp = field.NewInt64(tableName, "timestamp")
	_accessLog.Datestamp = field.NewInt64(tableName, "datestamp")
	_accessLog.Minutestamp = field.NewInt64(tableName, "minutestamp")
	_accessLog.Hourstamp = field.NewInt64(tableName, "hourstamp")
	_accessLog.Weeklystamp = field.NewInt64(tableName, "weeklystamp")
	_accessLog.Monthstamp = field.NewInt64(tableName, "monthstamp")
	_accessLog.Yearstamp = field.NewInt64(tableName, "yearstamp")
	_accessLog.TimeCost = field.NewFloat64(tableName, "time_cost")
	_accessLog.RequestLength = field.NewInt32(tableName, "request_length")
	_accessLog.Connection = field.NewString(tableName, "connection")
	_accessLog.ConnectionRequests = field.NewString(tableName, "connection_requests")
	_accessLog.URI = field.NewString(tableName, "uri")
	_accessLog.OriRequest = field.NewString(tableName, "ori_request")
	_accessLog.QueryString = field.NewString(tableName, "query_string")
	_accessLog.Status = field.NewInt32(tableName, "status")
	_accessLog.BytesSent = field.NewInt32(tableName, "bytes_sent")
	_accessLog.Referer = field.NewString(tableName, "referer")
	_accessLog.UserAgent = field.NewString(tableName, "user_agent")
	_accessLog.ForwardedFor = field.NewString(tableName, "forwarded_for")
	_accessLog.Host = field.NewString(tableName, "host")
	_accessLog.Node = field.NewString(tableName, "node")
	_accessLog.Upstream = field.NewString(tableName, "upstream")

	_accessLog.fillFieldMap()

	return _accessLog
}

type accessLog struct {
	accessLogDo accessLogDo

	ALL                field.Asterisk
	ID                 field.Int32
	RouteID            field.String
	RequestID          field.String
	ClientAddr         field.String
	IsoTime            field.String
	Timestamp          field.Int64
	Datestamp          field.Int64
	Minutestamp        field.Int64
	Hourstamp          field.Int64
	Weeklystamp        field.Int64
	Monthstamp         field.Int64
	Yearstamp          field.Int64
	TimeCost           field.Float64
	RequestLength      field.Int32
	Connection         field.String
	ConnectionRequests field.String
	URI                field.String
	OriRequest         field.String
	QueryString        field.String
	Status             field.Int32
	BytesSent          field.Int32
	Referer            field.String
	UserAgent          field.String
	ForwardedFor       field.String
	Host               field.String
	Node               field.String
	Upstream           field.String

	fieldMap map[string]field.Expr
}

func (a accessLog) Table(newTableName string) *accessLog {
	a.accessLogDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a accessLog) As(alias string) *accessLog {
	a.accessLogDo.DO = *(a.accessLogDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *accessLog) updateTableName(table string) *accessLog {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewInt32(table, "id")
	a.RouteID = field.NewString(table, "route_id")
	a.RequestID = field.NewString(table, "request_id")
	a.ClientAddr = field.NewString(table, "client_addr")
	a.IsoTime = field.NewString(table, "iso_time")
	a.Timestamp = field.NewInt64(table, "timestamp")
	a.Datestamp = field.NewInt64(table, "datestamp")
	a.Minutestamp = field.NewInt64(table, "minutestamp")
	a.Hourstamp = field.NewInt64(table, "hourstamp")
	a.Weeklystamp = field.NewInt64(table, "weeklystamp")
	a.Monthstamp = field.NewInt64(table, "monthstamp")
	a.Yearstamp = field.NewInt64(table, "yearstamp")
	a.TimeCost = field.NewFloat64(table, "time_cost")
	a.RequestLength = field.NewInt32(table, "request_length")
	a.Connection = field.NewString(table, "connection")
	a.ConnectionRequests = field.NewString(table, "connection_requests")
	a.URI = field.NewString(table, "uri")
	a.OriRequest = field.NewString(table, "ori_request")
	a.QueryString = field.NewString(table, "query_string")
	a.Status = field.NewInt32(table, "status")
	a.BytesSent = field.NewInt32(table, "bytes_sent")
	a.Referer = field.NewString(table, "referer")
	a.UserAgent = field.NewString(table, "user_agent")
	a.ForwardedFor = field.NewString(table, "forwarded_for")
	a.Host = field.NewString(table, "host")
	a.Node = field.NewString(table, "node")
	a.Upstream = field.NewString(table, "upstream")

	a.fillFieldMap()

	return a
}

func (a *accessLog) WithContext(ctx context.Context) *accessLogDo {
	return a.accessLogDo.WithContext(ctx)
}

func (a accessLog) TableName() string { return a.accessLogDo.TableName() }

func (a accessLog) Alias() string { return a.accessLogDo.Alias() }

func (a *accessLog) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *accessLog) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 27)
	a.fieldMap["id"] = a.ID
	a.fieldMap["route_id"] = a.RouteID
	a.fieldMap["request_id"] = a.RequestID
	a.fieldMap["client_addr"] = a.ClientAddr
	a.fieldMap["iso_time"] = a.IsoTime
	a.fieldMap["timestamp"] = a.Timestamp
	a.fieldMap["datestamp"] = a.Datestamp
	a.fieldMap["minutestamp"] = a.Minutestamp
	a.fieldMap["hourstamp"] = a.Hourstamp
	a.fieldMap["weeklystamp"] = a.Weeklystamp
	a.fieldMap["monthstamp"] = a.Monthstamp
	a.fieldMap["yearstamp"] = a.Yearstamp
	a.fieldMap["time_cost"] = a.TimeCost
	a.fieldMap["request_length"] = a.RequestLength
	a.fieldMap["connection"] = a.Connection
	a.fieldMap["connection_requests"] = a.ConnectionRequests
	a.fieldMap["uri"] = a.URI
	a.fieldMap["ori_request"] = a.OriRequest
	a.fieldMap["query_string"] = a.QueryString
	a.fieldMap["status"] = a.Status
	a.fieldMap["bytes_sent"] = a.BytesSent
	a.fieldMap["referer"] = a.Referer
	a.fieldMap["user_agent"] = a.UserAgent
	a.fieldMap["forwarded_for"] = a.ForwardedFor
	a.fieldMap["host"] = a.Host
	a.fieldMap["node"] = a.Node
	a.fieldMap["upstream"] = a.Upstream
}

func (a accessLog) clone(db *gorm.DB) accessLog {
	a.accessLogDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a accessLog) replaceDB(db *gorm.DB) accessLog {
	a.accessLogDo.ReplaceDB(db)
	return a
}

type accessLogDo struct{ gen.DO }

func (a accessLogDo) Debug() *accessLogDo {
	return a.withDO(a.DO.Debug())
}

func (a accessLogDo) WithContext(ctx context.Context) *accessLogDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a accessLogDo) ReadDB() *accessLogDo {
	return a.Clauses(dbresolver.Read)
}

func (a accessLogDo) WriteDB() *accessLogDo {
	return a.Clauses(dbresolver.Write)
}

func (a accessLogDo) Session(config *gorm.Session) *accessLogDo {
	return a.withDO(a.DO.Session(config))
}

func (a accessLogDo) Clauses(conds ...clause.Expression) *accessLogDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a accessLogDo) Returning(value interface{}, columns ...string) *accessLogDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a accessLogDo) Not(conds ...gen.Condition) *accessLogDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a accessLogDo) Or(conds ...gen.Condition) *accessLogDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a accessLogDo) Select(conds ...field.Expr) *accessLogDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a accessLogDo) Where(conds ...gen.Condition) *accessLogDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a accessLogDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *accessLogDo {
	return a.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (a accessLogDo) Order(conds ...field.Expr) *accessLogDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a accessLogDo) Distinct(cols ...field.Expr) *accessLogDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a accessLogDo) Omit(cols ...field.Expr) *accessLogDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a accessLogDo) Join(table schema.Tabler, on ...field.Expr) *accessLogDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a accessLogDo) LeftJoin(table schema.Tabler, on ...field.Expr) *accessLogDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a accessLogDo) RightJoin(table schema.Tabler, on ...field.Expr) *accessLogDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a accessLogDo) Group(cols ...field.Expr) *accessLogDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a accessLogDo) Having(conds ...gen.Condition) *accessLogDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a accessLogDo) Limit(limit int) *accessLogDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a accessLogDo) Offset(offset int) *accessLogDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a accessLogDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *accessLogDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a accessLogDo) Unscoped() *accessLogDo {
	return a.withDO(a.DO.Unscoped())
}

func (a accessLogDo) Create(values ...*model.AccessLog) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a accessLogDo) CreateInBatches(values []*model.AccessLog, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a accessLogDo) Save(values ...*model.AccessLog) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a accessLogDo) First() (*model.AccessLog, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.AccessLog), nil
	}
}

func (a accessLogDo) Take() (*model.AccessLog, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.AccessLog), nil
	}
}

func (a accessLogDo) Last() (*model.AccessLog, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.AccessLog), nil
	}
}

func (a accessLogDo) Find() ([]*model.AccessLog, error) {
	result, err := a.DO.Find()
	return result.([]*model.AccessLog), err
}

func (a accessLogDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.AccessLog, err error) {
	buf := make([]*model.AccessLog, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a accessLogDo) FindInBatches(result *[]*model.AccessLog, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a accessLogDo) Attrs(attrs ...field.AssignExpr) *accessLogDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a accessLogDo) Assign(attrs ...field.AssignExpr) *accessLogDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a accessLogDo) Joins(fields ...field.RelationField) *accessLogDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a accessLogDo) Preload(fields ...field.RelationField) *accessLogDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a accessLogDo) FirstOrInit() (*model.AccessLog, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.AccessLog), nil
	}
}

func (a accessLogDo) FirstOrCreate() (*model.AccessLog, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.AccessLog), nil
	}
}

func (a accessLogDo) FindByPage(offset int, limit int) (result []*model.AccessLog, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a accessLogDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a accessLogDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a accessLogDo) Delete(models ...*model.AccessLog) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *accessLogDo) withDO(do gen.Dao) *accessLogDo {
	a.DO = *do.(*gen.DO)
	return a
}
