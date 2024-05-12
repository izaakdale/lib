package listener

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

var (
	client                  *Client
	ErrClientNotInitialised = errors.New("uninitialised client")
)

const (
	defaultMaxNumberOfMessages = int32(10)
	defaultVisibiltyTimeout    = int32(5)
	defaultWaitTimeSeconds     = int32(10)
)

type (
	Client struct {
		sqsClient sqsConsumeAPI
		input     *sqs.ReceiveMessageInput
	}
	sqsConsumeAPI interface {
		ReceiveMessage(ctx context.Context,
			params *sqs.ReceiveMessageInput,
			optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
		DeleteMessage(ctx context.Context,
			params *sqs.DeleteMessageInput,
			optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
	}
	configOptions struct {
		endpoint            *string
		maxNumberOfMessages *int32
		visibilityTimeout   *int32
		waitTimeSeconds     *int32
	}
	option func(opt *configOptions) error
	// This function allows the user to define how their messages should be processed.
	ProcessorFunc func(context.Context, []byte) error
)

// Initialise creates a new client that listens to the queue specified and assigns it to the package client.
// Options include WithEndpoint, WithMaxNumerOfMessages, WithVisibilityTimeout and WithWaitTimeSeconds.
func Initialise(cfg aws.Config, queueURL string, optFuncs ...option) error {
	var options configOptions
	for _, optFunc := range optFuncs {
		err := optFunc(&options)
		if err != nil {
			return err
		}
	}

	var cli Client
	cli.input = &sqs.ReceiveMessageInput{
		QueueUrl: &queueURL,
	}
	cli.input.AttributeNames = []types.QueueAttributeName{
		types.QueueAttributeNameAll,
	}

	if options.endpoint == nil || *options.endpoint == "" {
		cli.sqsClient = sqs.NewFromConfig(cfg)
	} else {
		cli.sqsClient = sqs.NewFromConfig(cfg,
			sqs.WithEndpointResolver(
				sqs.EndpointResolverFromURL(*options.endpoint),
			))
	}

	if options.maxNumberOfMessages != nil {
		cli.input.MaxNumberOfMessages = *options.maxNumberOfMessages
	} else {
		cli.input.MaxNumberOfMessages = defaultMaxNumberOfMessages
	}
	if options.visibilityTimeout != nil {
		cli.input.VisibilityTimeout = *options.visibilityTimeout
	} else {
		cli.input.VisibilityTimeout = defaultVisibiltyTimeout
	}
	if options.waitTimeSeconds != nil {
		cli.input.WaitTimeSeconds = *options.waitTimeSeconds
	} else {
		cli.input.WaitTimeSeconds = defaultWaitTimeSeconds
	}

	client = &cli
	return nil
}

// Listen triggers a never ending for loop that continually requests the specified queue for messages.
func Listen(ctx context.Context, pf ProcessorFunc, errChan chan<- error) {
	if client == nil {
		errChan <- ErrClientNotInitialised
		return
	}
	for {
		msgResult, err := client.sqsClient.ReceiveMessage(ctx, client.input)
		if err != nil {
			errChan <- fmt.Errorf("failed to receive message: %w", err)
		}

		if msgResult != nil {
			if msgResult.Messages != nil {
				for _, m := range msgResult.Messages {
					err = pf(ctx, []byte(*m.Body))
					if err != nil {
						errChan <- err
					}
					dMInput := &sqs.DeleteMessageInput{
						QueueUrl:      client.input.QueueUrl,
						ReceiptHandle: m.ReceiptHandle,
					}
					_, err = client.sqsClient.DeleteMessage(ctx, dMInput)
					if err != nil {
						errChan <- err
					}
				}
			} else {
				continue
			}
		}
	}
}

// WithEndpoint adds an specific endpoint to be used by the AWS API.
// Useful for local development e.g. localstack URL.
func WithEndpoint(e string) option {
	return func(opt *configOptions) error {
		opt.endpoint = &e
		return nil
	}
}

// WithMaxNumerOfMessages dictaces how many messages can be returned from a queue in one go.
// Defaults to 10.
func WithMaxNumerOfMessages(n int32) option {
	return func(opt *configOptions) error {
		opt.maxNumberOfMessages = &n
		return nil
	}
}

// WithVisibilityTimeout dictates how long a queue message will be hidden from other clients.
// Defaults to 5 seconds.
func WithVisibilityTimeout(v int32) option {
	return func(opt *configOptions) error {
		opt.visibilityTimeout = &v
		return nil
	}
}

// WithWaitTimeSeconds dictates how long the request to the queue will wait for a message.
// Defaults to 10 seconds
func WithWaitTimeSeconds(s int32) option {
	return func(opt *configOptions) error {
		opt.waitTimeSeconds = &s
		return nil
	}
}
