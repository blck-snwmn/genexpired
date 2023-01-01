//go:generate ./genexpired --source=$GOFILE
package example

import "time"

type IClaim interface {
}

type (
	XClaim struct {
		expireAt time.Time
	}
)
type Claim struct {
	expireAt time.Time
}

func (c *Claim) Expired(now time.Time) bool {
	return !c.expireAt.Before(now)
}
