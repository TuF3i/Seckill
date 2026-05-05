package stringToNodeID

import "hash/fnv"

func StringToNodeID(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))

	val := int64(h.Sum64())
	if val < 0 {
		val = -val
	}
	return val % 1024
}
