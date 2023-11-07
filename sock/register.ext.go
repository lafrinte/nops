package sock

import (
	P "google.golang.org/protobuf/proto"
)

func (r *RegisterMsg) Marshal() ([]byte, error) {
	return P.Marshal(r)
}

func (r *RegisterMsg) Unmarshal(b []byte) error {
	return P.Unmarshal(b, r)
}
