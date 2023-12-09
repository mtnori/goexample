package proxyrolecreds

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"goexample/internal/proxyrolecreds/client"
	"net/http"
)

const ProviderName = `ProxyRoleProvider`

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func WithEndpoint(endpoint string) func(*Options) {
	return func(options *Options) {
		options.Endpoint = endpoint
	}
}

func WithHTTPClient(client HTTPClient) func(*Options) {
	return func(options *Options) {
		options.HTTPClient = client
	}
}

type Options struct {
	Endpoint   string
	HTTPClient HTTPClient
}

type getCredentialsAPIClient interface {
	GetCredentials(context.Context) (*client.GetCredentialsOutput, error)
}

type Provider struct {
	client getCredentialsAPIClient
}

func New(endpoint string, optFns ...func(*Options)) *Provider {
	o := Options{
		Endpoint: endpoint,
	}

	for _, fn := range optFns {
		fn(&o)
	}

	p := &Provider{
		client: client.New(client.Options{
			HTTPClient: o.HTTPClient,
			Endpoint:   o.Endpoint,
		}),
	}

	return p
}

func (p *Provider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	resp, err := p.getCredentials(ctx)
	if err != nil {
		return aws.Credentials{}, fmt.Errorf("failed to load credentials, %w", err)
	}

	creds := aws.Credentials{
		AccessKeyID:     resp.AccessKeyID,
		SecretAccessKey: resp.SecretAccessKey,
		SessionToken:    resp.Token,
		Source:          ProviderName,
		CanExpire:       true,
		Expires:         *resp.Expiration,
	}

	return creds, nil
}

func (p *Provider) getCredentials(ctx context.Context) (*client.GetCredentialsOutput, error) {
	return p.client.GetCredentials(ctx)
}
