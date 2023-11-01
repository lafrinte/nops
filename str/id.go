package str

import "github.com/rs/xid"

func ID() xid.ID {
	return xid.New()
}
