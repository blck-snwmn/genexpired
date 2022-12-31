//go:generate ./genexpired --source=$GOFILE
package example

import "time"

type Claim struct {
	expireAt time.Time
}

func (c *Claim) Expired(now time.Time) bool {
	return !c.expireAt.Before(now)
}
