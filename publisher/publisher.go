package publisher

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type (
	client struct {
		snsClient snsPublishAPI
		TopicArn  string
	}
	snsPublishAPI interface {
		Publish(ctx context.Context,
			params *sns.PublishInput,
			optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
	}
	configOptions struct {
		endpoint *string
	}
	option func(opt *configOptions) error
)

// New creates a new publisher client with the AWS config and TopicArn specific.
// Otional parameters include WithEndpoint.
func New(cfg aws.Config, topicArn string, optFuncs ...option) (*client, error) {
	var options configOptions
	for _, optFunc := range optFuncs {
		err := optFunc(&options)
		if err != nil {
			return nil, err
		}
	}

	var cli = client{
		TopicArn: topicArn,
	}
	if options.endpoint != nil {
		cli.snsClient = sns.NewFromConfig(cfg,
			sns.WithEndpointResolver(
				sns.EndpointResolverFromURL(*options.endpoint),
			))
	} else {
		cli.snsClient = sns.NewFromConfig(cfg)
	}

	return &cli, nil
}

// WithEndpoint adds an specific endpoint to be used by the AWS API.
// Useful for local development e.g. localstack URL.
func WithEndpoint(e string) option {
	return func(opt *configOptions) error {
		opt.endpoint = &e
		return nil
	}
}

// Publish sends a message to the Topic initialised in the client.
// Returns the message id and an error
func (c *client) Publish(msg *string) (*string, error) {
	input := &sns.PublishInput{
		Message:  msg,
		TopicArn: &c.TopicArn,
	}

	result, err := c.snsClient.Publish(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return result.MessageId, nil
}
