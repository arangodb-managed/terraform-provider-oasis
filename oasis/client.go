package oasis

import (
	"context"
	"crypto/tls"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"

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

func (c *Client) Connect() {
	ctx := context.Background()
	log := zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})

	c.conn = c.mustDialAPI(log)

	token, err := c.getToken(ctx, log, c.ApiKeyID, c.ApiKeySecret)
	if err != nil {
		log.Print("Could not get Auth Token")
	}

	c.ctxWithToken = auth.WithAccessToken(ctx, token)
}

// mustDialAPI dials the ArangoDB Oasis API
func (c *Client) mustDialAPI(log zerolog.Logger) *grpc.ClientConn {
	// Set up a connection to the server.
	tc := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.Dial(c.ApiEndpoint+c.ApiPortSuffix, grpc.WithTransportCredentials(tc))
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to ArangoDB Oasis API")
	}
	return conn
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
