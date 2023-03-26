package middleware

import (
	"chatgpt-bot/constant/config"
	botError "chatgpt-bot/constant/error"
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	mutex      sync.Mutex
	rate       float64
	capacity   int64
	duration   int64
	tokens     map[string]float64
	lastUpdate map[string]time.Time
}

func NewLimiter(capacity int64, duration int64) *Limiter {
	return &Limiter{
		rate:       float64(capacity) / float64(duration),
		capacity:   capacity,
		tokens:     make(map[string]float64),
		lastUpdate: make(map[string]time.Time),
		duration:   duration,
	}
}

func (l *Limiter) Allow(user string) (bool, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.refill(user)

	if l.tokens[user] >= 1 {
		l.tokens[user]--
		return true, nil
	} else {
		remainSecond := int64((1 - l.tokens[user]) / l.rate)
		text := fmt.Sprintf(botError.RateLimitMessageTemplate,
			l.capacity, l.duration/60, remainSecond,
			l.duration/60, l.capacity, remainSecond, config.AllowByInviteCount)
		return false, fmt.Errorf(text)
	}

}

func (l *Limiter) refill(user string) {
	if _, ok := l.tokens[user]; !ok {
		l.tokens[user] = float64(l.capacity)
		l.lastUpdate[user] = time.Now()
		return
	} else {
		now := time.Now()
		elapsed := now.Sub(l.lastUpdate[user])
		l.tokens[user] = l.tokens[user] + elapsed.Seconds()*l.rate
		if l.tokens[user] > float64(l.capacity) {
			l.tokens[user] = float64(l.capacity)
		}
		l.lastUpdate[user] = now
	}

}

func (l *Limiter) GetDuration() int64 {
	return l.duration
}

func (l *Limiter) GetCapacity() int64 {
	return l.capacity
}
