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
	"strings"
	"time"

	"github.com/guysports/go-betfair-api/pkg/betting"
	"github.com/guysports/go-betfair-api/pkg/types"
)

type (
	Login struct {
		RootCAPath string `json:"rootcapath"`
		CertPath   string `json:"certpath,required"`
		KeyPath    string `json:"keypath,required"`
		User       string `json:"user,required"`
		Password   string `json:"password,required"`
	}

	SessionData struct {
		Key       string    `json:"sessionkey,string"`
		ExpiresAt time.Time `json:"expiresat,string"`
	}
)

const (
	sessionExpiry = 4 * time.Hour // Set to whatever is configured in the Betfair account
	sessionFile   = "session.json"
)

// UnmarshalJSON implements custom unmarshaler for time object in session data
func (s *SessionData) UnmarshalJSON(b []byte) error {
	var sessionMap map[string]string

	if err := json.Unmarshal(b, &sessionMap); err != nil {
		return err
	}

	for k, v := range sessionMap {
		if strings.ToLower(k) == "sessionkey" {
			s.Key = v
		}
		if strings.ToLower(k) == "expiresat" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}
			s.ExpiresAt = t
		}
	}
	return nil
}

// MarshalJSON implements a custom marshaler for time object in session data
func (s *SessionData) MarshalJSON() ([]byte, error) {
	marshaler := struct {
		Key       string `json:"sessionkey"`
		ExpiresAt string `json:"expiresat"`
	}{
		Key:       s.Key,
		ExpiresAt: s.ExpiresAt.Format(time.RFC3339),
	}
	return json.Marshal(marshaler)
}

// NewLogin reads the specified path and unmarshals into a Login struct
func NewLogin(path string) (*Login, error) {
	login := Login{}

	loginData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(loginData, &login); err != nil {
		return nil, err
	}

	return &login, nil
}

// BetfairAuthenticate either
func (l *Login) BetfairAuthenticate(ctx context.Context, appKey string, client *betting.API) (*betting.API, error) {
	cfg := types.Config{
		CertPath: l.CertPath,
		KeyPath:  l.KeyPath,
		User:     l.User,
		Password: l.Password,
		AppKey:   appKey,
	}
	if l.RootCAPath != "" {
		cfg.RootCAPath = l.RootCAPath
	}
	var err error
	if client == nil {
		client, err = betting.NewAPI(ctx, &cfg)
		if err != nil {
			return nil, err
		}
	}

	// Check for existing session key and use if it hasn't expired
	sessionAuth := SessionData{}
	sessionDirExists := false
	sessionHome := os.Getenv("UNIT_TEST_HOME")
	if sessionHome == "" {
		sessionHome = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".betfair")
	}
	if _, err := os.Stat(sessionHome); !os.IsNotExist(err) {
		sessionDirExists = true
		// Check sessiondata
		sessionBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", sessionHome, sessionFile))
		if err == nil {
			// File exists then unmarshal
			err = json.Unmarshal(sessionBytes, &sessionAuth)
			if err == nil {
				// Check expiry and use session key if still valid
				if time.Now().Sub(sessionAuth.ExpiresAt) < sessionExpiry {
					client.Client.SetSessionKey(sessionAuth.Key)
					return client, nil
				}
			}
		}
	}

	return l.BetfairAuthenticateImpl(client, sessionHome, sessionDirExists)
}

func (l *Login) BetfairAuthenticateImpl(client *betting.API, sessionHome string, sessionDirExists bool) (*betting.API, error) {
	authData, err := client.Client.Authenticate()
	if err != nil {
		return nil, err
	}

	// Store the session key in the $HOME/.betfair directory with an expiry time for reuse
	if !sessionDirExists {
		err = os.Mkdir(sessionHome, 0666)
		if err != nil {
			return nil, err
		}
		os.Chmod(sessionHome, 0755)
	}
	sessionToStore := SessionData{
		Key:       authData.SessionToken,
		ExpiresAt: time.Now().Add(sessionExpiry),
	}
	// Not too fussed about an error, as it just means another login next time around
	bytesToWrite, _ := json.Marshal(&sessionToStore)
	_ = os.WriteFile(fmt.Sprintf("%s/%s", sessionHome, sessionFile), bytesToWrite, 0666)

	return client, nil
}
