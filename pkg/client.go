//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Joerg Schad, Gergely Brautigam
//

package pkg

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/arangodb-managed/apis/common/auth"
	iam "github.com/arangodb-managed/apis/iam/v1"
)

// Client is responsible for connecting to the Oasis API
type Client struct {
	ApiKeyID      string
	ApiKeySecret  string
	ApiEndpoint   string
	ApiPortSuffix string
	ctxWithToken  context.Context
	conn          *grpc.ClientConn
}

// Connect connects to oasis api
func (c *Client) Connect() error {
	ctx := context.Background()
	log := zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})

	var err error
	c.conn, err = c.mustDialAPI(log)
	if err != nil {
		return err
	}

	token, err := c.getToken(ctx, log, c.ApiKeyID, c.ApiKeySecret)
	if err != nil {
		log.Error().Err(err).Msg("Could not get Auth Token")
		return err
	}

	c.ctxWithToken = auth.WithAccessToken(ctx, token)
	return nil
}

// mustDialAPI dials the ArangoDB Oasis API
func (c *Client) mustDialAPI(log zerolog.Logger) (*grpc.ClientConn, error) {
	// Set up a connection to the server.
	tc := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.Dial(c.ApiEndpoint+c.ApiPortSuffix, grpc.WithTransportCredentials(tc))
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to ArangoDB Oasis API")
		return nil, err
	}
	return conn, nil
}

func (c *Client) getToken(ctx context.Context, log zerolog.Logger, apiKeyID, apiKeySecret string) (string, error) {

	iamc := iam.NewIAMServiceClient(c.conn)

	resp, err := iamc.AuthenticateAPIKey(ctx, &iam.AuthenticateAPIKeyRequest{
		Id:     apiKeyID,
		Secret: apiKeySecret,
	})
	if err != nil {
		log.Error().Err(err).Msg("Authentication failed")
		return "", err
	}
	log.Print("Retrieved Auth token successfully.")
	return resp.GetToken(), nil
}
