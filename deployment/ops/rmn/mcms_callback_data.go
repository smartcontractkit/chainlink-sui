package rmn

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

// SerializeSubjects BCS-encodes vector<vector<u8>> for curse_multiple callbacks.
func SerializeSubjects(subjects [][]byte) ([]byte, error) {
	s := &bcs.Serializer{}
	s.Uleb128(uint32(len(subjects)))
	for _, subject := range subjects {
		s.Uleb128(uint32(len(subject)))
		s.FixedBytes(subject)
	}
	return s.ToBytes(), nil
}

// SerializeMcmsObjectAddrsWithAddressVector BCS-encodes pinned object addresses followed by
// vector<address> for MCMS callbacks such as deregister_curser_cap_ids.
func SerializeMcmsObjectAddrsWithAddressVector(pinnedIDs []string, addresses []string) ([]byte, error) {
	data, err := SerializeMcmsObjectAddrs(pinnedIDs...)
	if err != nil {
		return nil, err
	}
	vectorData, err := serializeAddressVector(addresses)
	if err != nil {
		return nil, err
	}
	return append(data, vectorData...), nil
}

// SerializeMcmsObjectAddrsWithBool BCS-encodes pinned object addresses followed by a bool.
func SerializeMcmsObjectAddrsWithBool(pinnedIDs []string, enabled bool) ([]byte, error) {
	data, err := SerializeMcmsObjectAddrs(pinnedIDs...)
	if err != nil {
		return nil, err
	}
	s := &bcs.Serializer{}
	s.Bool(enabled)
	return append(data, s.ToBytes()...), nil
}

func serializeAddressVector(addresses []string) ([]byte, error) {
	s := &bcs.Serializer{}
	s.Uleb128(uint32(len(addresses)))
	for _, addr := range addresses {
		bytes, err := objectIDTo32Bytes(addr)
		if err != nil {
			return nil, err
		}
		s.FixedBytes(bytes[:])
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
	for i := range 32 {
		var b byte
		for j := range 2 {
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
