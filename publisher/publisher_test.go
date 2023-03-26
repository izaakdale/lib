package publisher_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/izaakdale/lib/publisher"
	"github.com/stretchr/testify/assert"
)

var (
	inputMessage  = "input"
	outputMessage = "output"
)

type stub struct {
	T *testing.T
	// Message *string
}

func (s stub) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	assert.Equal(s.T, inputMessage, *params.Message)
	return &sns.PublishOutput{
		MessageId: &outputMessage,
	}, nil
}

func TestInitialiseAndPublish(t *testing.T) {
	ctx := context.Background()
	s := stub{t}

	cfg, err := config.LoadDefaultConfig(ctx)
	assert.NoError(t, err)

	_, err = publisher.Publish(ctx, inputMessage)
	assert.Error(t, err)
	assert.EqualError(t, err, publisher.ErrClientNotInitialised.Error())

	err = publisher.Initialise(cfg, "arn:aws:sns:eu-west-2:000000000000:test-test-test", publisher.WithPublisher(s))
	assert.NoError(t, err)

	ret, err := publisher.Publish(ctx, inputMessage)
	assert.NoError(t, err)

	assert.Equal(t, *ret, outputMessage)
}
