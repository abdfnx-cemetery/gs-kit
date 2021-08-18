package retry

import (
	"context"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"

	"github.com/abdfnx/gs-kit/data"
)

type PolicyType int

const (
	PolicyConstant PolicyType = iota
	PolicyExponential
)

type Data struct {
	Policy PolicyType `mapstructure:"policy"`
	Duration time.Duration `mapstructure:"duration"`
	InitialInterval     time.Duration `mapstructure:"initialInterval"`
	RandomizationFactor float32       `mapstructure:"randomizationFactor"`
	Multiplier          float32       `mapstructure:"multiplier"`
	MaxInterval         time.Duration `mapstructure:"maxInterval"`
	MaxElapsedTime      time.Duration `mapstructure:"maxElapsedTime"`
	MaxRetries int64 `mapstructure:"maxRetries"`
}

func DefaultData() Data {
	return Data{
		Policy:              PolicyConstant,
		Duration:            5 * time.Second,
		InitialInterval:     backoff.DefaultInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      backoff.DefaultMaxElapsedTime,
		MaxRetries:          -1,
	}
}

func DefaultDataWithNoRetry() Data {
	c := DefaultData()
	c.MaxRetries = 0

	return c
}

func DecodeData(c *Data, input interface{}) error {
	var emptyData Data
	if *c == emptyData {
		*c = DefaultData()
	}

	return Data.Decode(input, c)
}

func DecodeDataWithPrefix(c *Data, input interface{}, prefix string) error {
	input, err := Data.PrefixedBy(input, prefix)
	if err != nil {
		return err
	}

	return DecodeData(c, input)
}

func (c *Data) NewBackOff() backoff.BackOff {
	var b backoff.BackOff
	switch c.Policy {
	case PolicyConstant:
		b = backoff.NewConstantBackOff(c.Duration)
	case PolicyExponential:
		eb := backoff.NewExponentialBackOff()
		eb.InitialInterval = c.InitialInterval
		eb.RandomizationFactor = float64(c.RandomizationFactor)
		eb.Multiplier = float64(c.Multiplier)
		eb.MaxInterval = c.MaxInterval
		eb.MaxElapsedTime = c.MaxElapsedTime
		b = eb
	}

	if c.MaxRetries >= 0 {
		b = backoff.WithMaxRetries(b, uint64(c.MaxRetries))
	}

	return b
}

func (c *Data) NewBackOffWithContext(ctx context.Context) backoff.BackOff {
	b := c.NewBackOff()

	return backoff.WithContext(b, ctx)
}

func NotifyRecover(operation backoff.Operation, b backoff.BackOff, notify backoff.Notify, recovered func()) error {
	var notified bool

	return backoff.RetryNotify(func() error {
		err := operation()

		if err == nil && notified {
			notified = false
			recovered()
		}

		return err
	}, b, func(err error, d time.Duration) {
		if !notified {
			notify(err, d)
			notified = true
		}
	})
}

func (p *PolicyType) DecodeString(value string) error {
	switch strings.ToLower(value) {
		case "constant":
			*p = PolicyConstant
		case "exponential":
			*p = PolicyExponential
		default:
			return errors.Errorf("unexpected back off policy type: %s", value)
	}

	return nil
}
