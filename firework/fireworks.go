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

// Init :
func (f *CronFirework) Init(rule string) {
	f.DelayFirework = &DelayFirework{
		Time: time.Now(),
	}
	f.Expression = cronexpr.MustParse(rule)
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
