package application

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
)

type appJoinSerJoinChn struct {
	ID             int32      `gorm:"column:id;type:integer;primaryKey;autoIncrement:true" json:"id"`
	Name           *string    `gorm:"column:name;type:character varying" json:"name"`
	ChannelID      *int32     `gorm:"column:channel_id;type:integer" json:"channel_id"`
	ChannelName    *string    `gorm:"column:channel_name;type:integer" json:"channel_name"`
	ServiceID      *int32     `gorm:"column:service_id;type:integer" json:"service_id"`
	ServiceName    *string    `gorm:"column:service_name;type:integer" json:"service_name"`
	ServiceType    *string    `gorm:"column:service_type;type:character varying" json:"service_type"`
	ApiSet         model.JSON `gorm:"column:api_set;type:jsonb" json:"api_set"`
	DocumentID     *int32     `gorm:"column:document_id;type:integer" json:"document_id"`
	RouteID        *string    `gorm:"column:route_id;type:character varying" json:"route_id"`
	Status         *int32     `gorm:"column:status;type:integer" json:"status"`
	Abandoned      *bool      `gorm:"column:abandoned;type:boolean;not null;default:false" json:"abandoned"`
	DelUniqueKey   *int32     `gorm:"column:del_unique_key;type:integer;not null;default:0" json:"del_unique_key"`
	Creator        *string    `gorm:"column:creator;type:character varying" json:"creator"`
	Ctime          time.Time  `gorm:"column:ctime;type:timestamp without time zone;not null;default:CURRENT_TIMESTAMP" json:"ctime"`
	LastEditor     *string    `gorm:"column:last_editor;type:character varying" json:"last_editor"`
	LastUpdateTime time.Time  `gorm:"column:last_update_time;type:timestamp without time zone;not null;default:CURRENT_TIMESTAMP" json:"last_update_time"`
}

func enableClickHouseLogger(plugins map[string]interface{}) map[string]interface{} {
	var clickHouseLogger = make(map[string]interface{})

	clickHouseLogger[`database`] = viper.GetString("clickhouse.database")
	clickHouseLogger[`logtable`] = "logger_tab"
	clickHouseLogger[`password`] = viper.GetString("clickhouse.password")
	clickHouseLogger[`endpoint_addr`] = fmt.Sprintf(`http://%s:%s`,
		viper.GetString("clickhouse.host"),
		viper.GetString("clickhouse.tcp-port"),
	)
	clickHouseLogger[`user`] = viper.GetString(`clickhouse.username`)

	plugins[`clickhouse-logger`] = clickHouseLogger

	return plugins
}

func enableSqlLogger(plugins map[string]interface{}) map[string]interface{} {
	var sqlLogger = make(map[string]interface{})
	sqlLogger[`host`] = viper.GetString("db.host")
	sqlLogger[`port`] = strconv.Itoa(viper.GetInt("db.port"))
	sqlLogger[`dbname`] = viper.GetString("db.database")
	sqlLogger[`username`] = viper.GetString("db.username")
	sqlLogger[`password`] = viper.GetString("db.password")
	sqlLogger[`tablename`] = "access_log"
	key := viper.GetString("db.type") + "-logger"
	plugins[key] = sqlLogger

	return plugins
}

func enableBodyLogger(plugins map[string]interface{}) map[string]interface{} {
	var bodyLogger = make(map[string]interface{})

	bodyLogger[`path`] = viper.GetString("control.body-logger-dir")

	plugins[`body-logger`] = bodyLogger

	return plugins
}

func configLimitCount(plugins map[string]interface{}) map[string]interface{} {
	limitCount := plugins["limit-count"].(map[string]interface{})
	limitCount["policy"] = "local"
	plugins["limit-count"] = limitCount

	return plugins
}

func configClientControl(plugins map[string]interface{}) map[string]interface{} {
	clientControl := plugins["client-control"].(map[string]interface{})
	if clientControl["max_body_size"] != nil {
		bodySize := clientControl["max_body_size"].(float64)
		clientControl["max_body_size"] = int64(bodySize)
	}
	plugins["client-control"] = clientControl

	return plugins
}

func verifyCreateNameUnique(name string, svcCtx *svc.ServiceContext, ctx context.Context) error {
	var (
		aT = query.Use(svcCtx.Db).Application
	)

	count, err := aT.WithContext(ctx).
		Where(aT.Name.Eq(name),
			aT.Abandoned.Is(false)).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return code.WithCodeMsg(code.DuplicateAppName)
	}

	return nil
}

func verifyUpdateNameUnique(name string, id int32, svcCtx *svc.ServiceContext, ctx context.Context) error {
	var (
		aT = query.Use(svcCtx.Db).Application
	)

	count, err := aT.WithContext(ctx).
		Where(aT.Name.Eq(name),
			aT.ID.Neq(id),
			aT.Abandoned.Is(false)).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return code.WithCodeMsg(code.DuplicateAppName)
	}

	return nil
}

// func getApplicationInfo(appID int32, svcCtx *svc.ServiceContext, ctx context.Context) (*model.Application, error) {
// 	var (
// 		aT = query.Use(svcCtx.Db).Application
// 	)

// 	application, err := aT.WithContext(ctx).
// 		JoinService().
// 		JoinChannel().
// 		Where(aT.ID.Eq(appID)).
// 		FindOne()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return application, nil
// }
