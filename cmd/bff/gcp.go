package main

import (
	"context"
	"fmt"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
	grpccredentials "google.golang.org/grpc/credentials"
)

type perRPCCredentials struct {
	token string
}

var _ grpccredentials.PerRPCCredentials = (*perRPCCredentials)(nil)

func (p *perRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + p.token,
	}, nil
}

func (p *perRPCCredentials) RequireTransportSecurity() bool {
	return true
}

func makeCreds(ctx context.Context, endpoint, saEmail string) (grpccredentials.PerRPCCredentials, error) {
	c, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("create IAM credentials client: %w", err)
	}
	defer c.Close()

	// https://cloud.google.com/iam/docs/creating-short-lived-service-account-credentials#sa-credentials-oidc
	resp, err := c.GenerateIdToken(ctx, &credentialspb.GenerateIdTokenRequest{
		Name:     fmt.Sprintf("projects/-/serviceAccounts/%s", saEmail),
		Audience: endpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &perRPCCredentials{
		token: resp.GetToken(),
	}, nil
}
