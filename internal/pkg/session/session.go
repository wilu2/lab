package session

import (
	"errors"
	"time"

	"github.com/spf13/viper"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	authType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
	"gitlab.intsig.net/textin-gateway/pkg/md5"
)

type Session struct {
	ID        int32     `gorm:"column:id;type:integer;primaryKey;autoIncrement:true" json:"id"`
	UserID    *int32    `gorm:"column:user_id;type:integer" json:"user_id"`
	Session   *string   `gorm:"column:session;type:character varying" json:"session"`
	Abandoned *bool     `gorm:"column:abandoned;type:boolean;not null;default:false" json:"abandoned"`
	Expiry    time.Time `gorm:"column:expiry;type:timestamp without time zone;not null;default:CURRENT_TIMESTAMP" json:"expiry"`
}

func GeneSession(userInfo authType.UserInfo) (string, error) {
	var (
		timeOut = viper.GetDuration("jwt.timeout")
		key     = viper.GetString("jwt.key")
	)
	dbIns, _ := dal.GetPostgresFactoryOr(nil)
	session := md5.GetMD5Encode(apisix.GetFlakeUidStr() + string(userInfo.ID) + key)
	// session := md5.GetMD5Encode(apisix.GetFlakeUidStr() + userInfo.Role + string(userInfo.ID) + key)
	expiry := time.Now().Add(timeOut)
	sessionItem := Session{
		UserID:  &userInfo.ID,
		Session: &session,
		Expiry:  expiry,
	}
	dbIns.GetDb().Create(&sessionItem)
	if sessionItem.ID == 0 { // 存在用户删除，jwt 还存在的情况
		return "", errors.New("")
	}
	return session, nil
}

func ParseSession(session string) (userInfo authType.UserInfo, err error) {
	var (
		dbIns, _    = dal.GetPostgresFactoryOr(nil)
		abandoned   = false
		sessionItem Session
		userItem    model.User
	)
	dbIns.GetDb().Where(&Session{Session: &session, Abandoned: &abandoned}).First(
		sessionItem)
	if sessionItem.ID == 0 { // 存在用户删除，jwt 还存在的情况
		return authType.UserInfo{}, errors.New("")
	}
	dbIns.GetDb().Where(&model.User{ID: *sessionItem.UserID}).First(
		userItem)
	channels := utils.ConvertChannels(userItem.Channels)
	return authType.UserInfo{
		ID: userInfo.ID,
		// Role:     *userItem.Role,
		Channels: channels,
	}, nil
}
