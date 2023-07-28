package stringx

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJoinStringsInASCII(t *testing.T) {
	data := map[string]string{}
	data["id"] = "12323"
	data["version"] = "2ss2"
	data["timestamp"] = "1661503953"
	data["nonce"] = "23sd"
	data["appid"] = "ssss"
	val := JoinStringsInASCII(data, "&", false, false)
	require.Equal(t, "appid=ssss&id=12323&nonce=23sd&timestamp=1661503953&version=2ss2", val)
}
