package stringToNodeID

import "hash/fnv"

func StringToNodeID(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}
