package Conv

import (
	"encoding/binary"
)

func IntToBytes(i int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}
func BytesToInt(i []byte) int {
	if len(i) != 8 {
		return 0
	}
	return int(binary.BigEndian.Uint64(i))
}
