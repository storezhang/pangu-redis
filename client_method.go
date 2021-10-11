package redis

import (
	`context`

	`github.com/go-redis/redis/v8`
)

func (c *Client) addValues(ctx context.Context, key string, pushType addValuesType, opts ...valuesOption) (affected int64, err error) {
	_options := defaultValuesOptions()
	for _, opt := range opts {
		opt.applyValues(_options)
	}

	values := make([]interface{}, 0, len(_options.values))
	for _, value := range _options.values {
		var marshaled interface{}
		if marshaled, err = c.marshal(value, _options.label, _options.serializer); nil != err {
			return
		}

		values = append(values, marshaled)
	}

	client := c.getClient(_options.options)
	switch pushType {
	case addValuesTypeLPush:
		affected, err = client.LPush(ctx, key, values...).Result()
	case addValuesTypeRPush:
		affected, err = client.RPush(ctx, key, values...).Result()
	}

	return
}

func (c *Client) len(ctx context.Context, key string, lenType lenType, opts ...option) (int64, error) {
	_options := defaultOptions()
	for _, opt := range opts {
		opt.apply(_options)
	}

	var redisCmd *redis.IntCmd
	client := c.getClient(_options)
	switch lenType {
	case lenTypeLLen:
		redisCmd = client.LLen(ctx, key)
	case lenTypeZCard:
		redisCmd = client.ZCard(ctx, key)
	}

	return redisCmd.Result()
}