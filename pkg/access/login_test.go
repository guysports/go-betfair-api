// Copyright 2022 Guy Barden
// login.go - Authenticator to the Betfair Exchange API

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package access

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/guysports/go-betfair-api/pkg/fake"

	"github.com/guysports/go-betfair-api/pkg/betting"
	"github.com/stretchr/testify/assert"
)

func setupTest(tb testing.TB) func(tb testing.TB) {
	os.Setenv("UNIT_TEST_HOME", "../../.betfair")

	return func(tb testing.TB) {
		os.RemoveAll(os.Getenv("UNIT_TEST_HOME"))
		defer os.Unsetenv("UNIT_TEST_HOME")
		return
	}
}

func TestLogin_BetfairAuthenticate(t *testing.T) {
	type fields struct {
		RootCAPath string
		CertPath   string
		KeyPath    string
		User       string
		Password   string
	}
	type args struct {
		ctx    context.Context
		appKey string
		client *betting.API
	}

	var testFields = fields{
		RootCAPath: "rootca.pem",
		CertPath:   "cert.crt",
		KeyPath:    "key.pem",
		User:       "testuser",
		Password:   "testpass",
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       *betting.API
		createFile bool
		wantErr    bool
		wantKey    string
	}{
		{
			name:   "successful authenticate with no session token",
			fields: testFields,
			args: args{
				ctx:    context.TODO(),
				appKey: "appkey",
				client: &betting.API{
					Client: &fake.FakeTransportClient{},
				},
			},
			want: &betting.API{
				Client: &fake.FakeTransportClient{},
			},
			wantKey: "validsessionkey",
		},
		{
			name:   "successful authenticate with session token",
			fields: testFields,
			args: args{
				ctx:    context.TODO(),
				appKey: "appkey",
				client: &betting.API{
					Client: &fake.FakeTransportClient{},
				},
			},
			want: &betting.API{
				Client: &fake.FakeTransportClient{},
			},
			wantKey:    "key",
			createFile: true,
		},
		{
			name:   "error in authenticate",
			fields: testFields,
			args: args{
				ctx:    context.TODO(),
				appKey: "appkey",
				client: &betting.API{
					Client: &fake.FakeTransportClient{
						SessionExpiredError: true,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tearDownTest := setupTest(t)
			defer tearDownTest(t)

			l := &Login{
				RootCAPath: tt.fields.RootCAPath,
				CertPath:   tt.fields.CertPath,
				KeyPath:    tt.fields.KeyPath,
				User:       tt.fields.User,
				Password:   tt.fields.Password,
			}
			if tt.createFile {
				// Create a valid session file
				session := SessionData{
					Key:       "key",
					ExpiresAt: time.Now().Add(sessionExpiry),
				}
				assert.Nil(t, os.Mkdir(os.Getenv("UNIT_TEST_HOME"), 0666))
				assert.Nil(t, os.Chmod(os.Getenv("UNIT_TEST_HOME"), 0755))
				bytesToWrite, _ := session.MarshalJSON()
				_ = os.WriteFile(fmt.Sprintf("%s/%s", os.Getenv("UNIT_TEST_HOME"), sessionFile), bytesToWrite, 0666)
			}

			got, err := l.BetfairAuthenticate(tt.args.ctx, tt.args.appKey, tt.args.client)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
				wantBytes, _ := os.ReadFile(fmt.Sprintf("%s/%s", os.Getenv("UNIT_TEST_HOME"), sessionFile))
				auth := SessionData{}
				_ = json.Unmarshal(wantBytes, &auth)
				assert.Equal(t, tt.wantKey, auth.Key)
			}
		})
	}
}
