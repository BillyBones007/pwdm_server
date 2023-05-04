package tokentools

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetUUID(t *testing.T) {
	jwtTools := NewJWTTools()
	tests := []struct {
		name      string
		jwtTools  *JWTTools
		expTime   int64
		sleepTime time.Duration
		uuid      string
		wantUUID  string
		wantValid bool
	}{
		{
			name:      "Valid token",
			jwtTools:  jwtTools,
			expTime:   time.Now().Add(time.Second * 3).Unix(),
			sleepTime: time.Second * 1,
			uuid:      "myID",
			wantUUID:  "myID",
		},
		{
			name:      "Invalid token",
			jwtTools:  jwtTools,
			expTime:   time.Now().Add(time.Second * 3).Unix(),
			sleepTime: time.Second * 1,
			uuid:      "",
			wantUUID:  "",
		},
		{
			name:      "Token is expired",
			jwtTools:  jwtTools,
			expTime:   time.Now().Add(time.Second * 1).Unix(),
			sleepTime: time.Second * 2,
			uuid:      "id",
			wantUUID:  "id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := tt.jwtTools.CreateToken(tt.expTime, tt.uuid)
			if err != nil {
				fmt.Println(err)
				t.Fail()
			}

			time.Sleep(tt.sleepTime)

			// ok, _ := tt.jwtTools.ValidToken(token)
			// assert.Equal(t, ok, tt.wantValid)

			uuid, err := tt.jwtTools.ParseUUID(token)
			if err != nil {
				if strings.Contains(err.Error(), "token is expired") || strings.Contains(err.Error(), "uuid field is empty") {
					fmt.Println("ok")
				}
			} else {
				assert.Equal(t, tt.wantUUID, uuid)
			}
		})
	}
}
