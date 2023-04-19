package publisher

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var (
	client                  *Client
	ErrClientNotInitialised = errors.New("uninitialised client")
)

type (
	Client struct {
		sns      snsPublishAPI
		TopicArn string
	}
	snsPublishAPI interface {
		Publish(ctx context.Context,
			params *sns.PublishInput,
			optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
	}
	configOptions struct {
		endpoint  *string
		publisher snsPublishAPI
	}
	option func(opt *configOptions) error
)

// Initialise creates a new publisher client and assigns it to the package level client.
// Otional parameters include WithEndpoint.
func Initialise(cfg aws.Config, topicArn string, optFuncs ...option) error {
	var options configOptions
	for _, optFunc := range optFuncs {
		err := optFunc(&options)
		if err != nil {
			return err
		}
	}

	if options.publisher != nil {
		client = &Client{
			TopicArn: topicArn,
			sns:      options.publisher,
		}
		return nil
	}

	var cli = Client{
		TopicArn: topicArn,
	}

	if options.endpoint == nil || *options.endpoint == "" {
		cli.sns = sns.NewFromConfig(cfg)
	} else {
		cli.sns = sns.NewFromConfig(cfg,
			sns.WithEndpointResolver(
				sns.EndpointResolverFromURL(*options.endpoint),
			))
	}

	client = &cli

	return nil
}

// WithEndpoint adds an specific endpoint to be used by the AWS API.
// Useful for local development e.g. localstack URL.
func WithEndpoint(e string) option {
	return func(opt *configOptions) error {
		opt.endpoint = &e
		return nil
	}
}

// WithPublisher allows the client to use their own publisher with the package.
// Useful for stubbing.
func WithPublisher(p snsPublishAPI) option {
	return func(opt *configOptions) error {
		opt.publisher = p
		return nil
	}
}

// Publish sends a message to the Topic initialised in the client.
// Returns the message id and an error
func Publish(ctx context.Context, msg string) (*string, error) {
	if client == nil {
		return nil, ErrClientNotInitialised
	}

	input := &sns.PublishInput{
		Message:  &msg,
		TopicArn: &client.TopicArn,
	}
	result, err := client.sns.Publish(ctx, input)
	if err != nil {
		return nil, err
	}

	return result.MessageId, nil
}
