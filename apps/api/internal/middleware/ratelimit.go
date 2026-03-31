package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
	"github.com/terraforge-gg/terraforge/internal/dto"
)

type RateLimiterConfig struct {
	Requests  int
	Window    time.Duration
	Namespace string
}

var (
	RateLimitGeneral = RateLimiterConfig{Requests: 100, Window: time.Minute, Namespace: "general"}
	RateLimitWrite   = RateLimiterConfig{Requests: 30, Window: time.Minute, Namespace: "write"}
	RateLimitSearch  = RateLimiterConfig{Requests: 60, Window: time.Minute, Namespace: "search"}
)

const slidingWindowScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local clearBefore = now - window

redis.call('ZREMRANGEBYSCORE', key, 0, clearBefore)
local count = redis.call('ZCARD', key)

if count < limit then
    redis.call('ZADD', key, now, now .. ':' .. math.random(1000000))
    redis.call('EXPIRE', key, math.ceil(window / 1000000000) + 1)
    return {1, limit - count - 1}
else
    local oldest = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
    local resetAt = now + window
    if #oldest > 0 then
        resetAt = tonumber(oldest[2]) + window
    end
    return {0, resetAt}
end
`

var slidingWindowRateLimiter = redis.NewScript(slidingWindowScript)

func RateLimiter(rdb *redis.Client, cfg RateLimiterConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			key := buildKey(c, cfg.Namespace)

			now := time.Now().UnixNano()
			window := cfg.Window.Nanoseconds()

			result, err := slidingWindowRateLimiter.Run(
				c.Request().Context(), rdb,
				[]string{key}, now, window, cfg.Requests,
			).Result()

			if err != nil {
				return next(c)
			}

			vals := result.([]interface{})

			allowed := vals[0].(int64) == 1

			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(cfg.Requests))
			c.Response().Header().Set("X-RateLimit-Window", cfg.Window.String())

			if allowed {
				remaining := vals[1].(int64)
				c.Response().Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
				return next(c)
			}

			resetAt := vals[1].(int64)
			retryAfter := (resetAt - now) / int64(time.Second)
			if retryAfter < 1 {
				retryAfter = 1
			}

			c.Response().Header().Set("X-RateLimit-Remaining", "0")
			c.Response().Header().Set("Retry-After", strconv.FormatInt(retryAfter, 10))

			return c.JSON(http.StatusTooManyRequests, dto.ProblemDetails{
				Title:  "Too Many Requests",
				Status: http.StatusTooManyRequests,
				Detail: fmt.Sprintf("Rate limit exceeded. Try again in %d seconds.", retryAfter),
			})
		}
	}
}

func buildKey(c *echo.Context, namespace string) string {
	identifier := c.RealIP()
	path := c.Request().URL.Path
	return fmt.Sprintf("ratelimit:%s:%s:%s", namespace, identifier, path)
}
