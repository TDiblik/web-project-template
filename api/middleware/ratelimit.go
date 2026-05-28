package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/TDiblik/project-template/api/utils"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
	return i
}

func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}

	return limiter
}

var defaultRateLimiter = NewIPRateLimiter(rate.Every(time.Minute/100), 100)

func RateLimit() fiber.Handler {
	return func(c fiber.Ctx) error {
		limiter := defaultRateLimiter.GetLimiter(c.IP())
		if !limiter.Allow() {
			return utils.TooManyRequestsResponse(c, fmt.Errorf("rate limit exceeded"))
		}
		return c.Next()
	}
}
