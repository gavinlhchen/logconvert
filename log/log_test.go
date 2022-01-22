package log_test

import (
	"context"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/gavinlhchen/logconvert/log"
)

func Test_WithName(t *testing.T) {
	defer log.Flush()

	logger := log.WithName("test")
	logger.Infow("Hello world!", "foo", "bar") //structed logger
	logger.Info("This is a info message", log.Int32("int_key", 10))
	logger.Infof("This is a formatted %s message", "info")
}

func Test_WithValues(t *testing.T) {
	defer log.Flush()

	logger := log.WithValues("key", "value")
	logger.Info("Hello world!")
}

func Test_V(t *testing.T) {
	defer log.Flush()

	log.V(-1).Infow("Hello world!", "key", "value")
	log.V(0).Infow("Hello world!", "key", "value")
	log.V(1).Infow("Hello world!", "key", "value")
	log.V(2).Infow("Hello world!", "key", "value")
}

func Test_Option(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ExitOnError)
	opt := log.NewOptions()
	opt.AddFlags(fs)

	args := []string{"--log.level=debug"}
	err := fs.Parse(args)
	assert.Nil(t, err)

	assert.Equal(t, "debug", opt.Level)
}

func Test_Context(t *testing.T) {
	defer log.Flush()

	ctx := context.Background()
	ctxSub := context.WithValue(ctx, "requestID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")

	log.L(ctxSub).Infow("user create function called.")
}
