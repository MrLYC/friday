package firework

import (
	"time"

	"github.com/gorhill/cronexpr"
)

// CronFirework :
type CronFirework struct {
	*DelayFirework
	Expression *cronexpr.Expression
}

// NewCronFirework :
func NewCronFirework(rule string, startAt time.Time, firework IFirework) *CronFirework {
	f := &CronFirework{}
	f.DelayFirework = &DelayFirework{}
	f.Time = startAt
	f.IFirework = firework
	f.Expression = cronexpr.MustParse(rule)
	return f
}

// Copy :
func (f *CronFirework) Copy() IFirework {
	return &CronFirework{
		Expression:    f.Expression,
		DelayFirework: f.DelayFirework.Copy().(*DelayFirework),
	}
}

// UpdateTime :
func (f *CronFirework) UpdateTime() bool {
	f.Time = f.Expression.Next(f.Time)
	return true
}
