package listener

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	defaultMaxNumberOfMessages = int32(10)
	defaultVisibiltyTimeout    = int32(5)
	defaultWaitTimeSeconds     = int32(10)
)

type (
	client struct {
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

	// Message is a wrapper for the AWS message that is returned from SQS.
	Message struct {
		Type      string    `json:"Type"`
		MessageID string    `json:"MessageId"`
		TopicArn  string    `json:"TopicArn"`
		Message   string    `json:"Message"`
		Timestamp time.Time `json:"Timestamp"`
	}

	// This function is where the logic for the data goes.
	ProcessorFunc func(Message) error
)

// New returns a new client that listens to the queue specified
// Options include WithEndpoint, WithMaxNumerOfMessages, WithVisibilityTimeout and WithWaitTimeSeconds.
func New(cfg aws.Config, queueURL string, optFuncs ...option) (*client, error) {
	var options configOptions
	for _, optFunc := range optFuncs {
		err := optFunc(&options)
		if err != nil {
			return nil, err
		}
	}

	var cli client
	cli.input = &sqs.ReceiveMessageInput{
		QueueUrl: &queueURL,
	}
	cli.input.AttributeNames = []types.QueueAttributeName{
		types.QueueAttributeNameAll,
	}

	if options.endpoint != nil {
		cli.sqsClient = sqs.NewFromConfig(cfg,
			sqs.WithEndpointResolver(
				sqs.EndpointResolverFromURL(*options.endpoint),
			))
	} else {
		cli.sqsClient = sqs.NewFromConfig(cfg)
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

	return &cli, nil
}

// Listen triggers a never ending for loop that continually requests the specified queue for messages.
func (c *client) Listen(pf ProcessorFunc, errChan chan<- error) {
	for {
		msgResult, err := c.sqsClient.ReceiveMessage(context.TODO(), c.input)
		if err != nil {
			errChan <- err
		}

		if msgResult.Messages != nil {
			for _, m := range msgResult.Messages {
				var messageToProcess Message
				err := json.Unmarshal([]byte(*m.Body), &messageToProcess)
				if err != nil {
					errChan <- err
				}
				err = pf(messageToProcess)
				if err != nil {
					errChan <- err
				}
				dMInput := &sqs.DeleteMessageInput{
					QueueUrl:      c.input.QueueUrl,
					ReceiptHandle: m.ReceiptHandle,
				}
				_, err = c.sqsClient.DeleteMessage(context.TODO(), dMInput)
				if err != nil {
					errChan <- err
				}
			}
		} else {
			continue
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
