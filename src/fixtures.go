package gc

import "time"

type Fixtures interface {
	SetTime(t time.Time)
	GetTime() time.Time
	SetValue(s string)
	GetValue() string
	SetToken(s string)
	GetToken() string
}
