//
// DISCLAIMER
//
// Copyright 2020-2023 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//

package provider

import (
	"context"
	"crypto/tls"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/arangodb-managed/apis/common/auth"
	commonGrpc "github.com/arangodb-managed/apis/common/v1/grpc"
	iam "github.com/arangodb-managed/apis/iam/v1"
	lh "github.com/arangodb-managed/log-helper"
)

// Client is responsible for connecting to the Arango Graph API
type Client struct {
	ApiKeyID       string
	ApiKeySecret   string
	ApiEndpoint    string
	ApiPortSuffix  string
	ProjectID      string
	OrganizationID string
	ctxWithToken   context.Context
	conn           *grpc.ClientConn
	log            zerolog.Logger
}

// Connect connects to Arango Graph API
func (c *Client) Connect() error {
	ctx := context.Background()
	c.log = lh.MustNew(lh.DefaultConfig())

	var err error
	c.conn, err = c.mustDialAPI()
	if err != nil {
		return err
	}

	token, err := c.getToken(ctx, c.ApiKeyID, c.ApiKeySecret)
	if err != nil {
		c.log.Error().Err(err).Msg("Could not get Auth Token")
		return err
	}
	// Add Access Token
	ctxWithToken := auth.WithAccessToken(ctx, token)
	// Add the User Agent as well
	ua := commonGrpc.CreateUserAgent("terraform-provider-oasis", currentVersion)
	c.ctxWithToken = commonGrpc.WithUserAgent(ctxWithToken, ua)
	return nil
}

// mustDialAPI dials the Arango Graph API
func (c *Client) mustDialAPI() (*grpc.ClientConn, error) {
	// Set up a connection to the server.
	tc := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.Dial(c.ApiEndpoint+c.ApiPortSuffix, grpc.WithTransportCredentials(tc))
	if err != nil {
		c.log.Error().Err(err).Msg("Failed to connect to Arango Graph API")
		return nil, err
	}
	return conn, nil
}

func (c *Client) getToken(ctx context.Context, apiKeyID, apiKeySecret string) (string, error) {

	iamc := iam.NewIAMServiceClient(c.conn)

	resp, err := iamc.AuthenticateAPIKey(ctx, &iam.AuthenticateAPIKeyRequest{
		Id:     apiKeyID,
		Secret: apiKeySecret,
	})
	if err != nil {
		c.log.Error().Err(err).Msg("Authentication failed")
		return "", err
	}
	c.log.Print("Retrieved Auth token successfully.")
	return resp.GetToken(), nil
}
