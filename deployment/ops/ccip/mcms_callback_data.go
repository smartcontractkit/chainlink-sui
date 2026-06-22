package ccipops

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
)

// SerializeMcmsObjectAddrs BCS-encodes fixed 32-byte object addresses for MCMS callback
// validate_obj_addrs payloads. Order must match the Move callback's deserialize order.
func SerializeMcmsObjectAddrs(ids ...string) ([]byte, error) {
	s := &bcs.Serializer{}
	for _, id := range ids {
		addr, err := objectIDTo32Bytes(id)
		if err != nil {
			return nil, err
		}
		s.FixedBytes(addr[:])
	}
	return s.ToBytes(), nil
}

func objectIDTo32Bytes(id string) ([32]byte, error) {
	var out [32]byte
	if len(id) < 2 || id[:2] != "0x" {
		id = "0x" + id
	}
	for len(id) < 66 {
		id = "0x0" + id[2:]
	}
	if len(id) != 66 {
		return out, fmt.Errorf("invalid object id length %q", id)
	}
	for i := 0; i < 32; i++ {
		var b byte
		for j := 0; j < 2; j++ {
			c := id[2+i*2+j]
			switch {
			case c >= '0' && c <= '9':
				b = b*16 + (c - '0')
			case c >= 'a' && c <= 'f':
				b = b*16 + (c - 'a' + 10)
			case c >= 'A' && c <= 'F':
				b = b*16 + (c - 'A' + 10)
			default:
				return out, fmt.Errorf("invalid hex in object id %q", id)
			}
		}
		out[i] = b
	}
	return out, nil
}
