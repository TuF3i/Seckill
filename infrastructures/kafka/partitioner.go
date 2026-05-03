package kafka

import (
	"hash/fnv"

	"github.com/segmentio/kafka-go"
)

type KeyPartitioner struct{}

func (p *KeyPartitioner) Balance(msg kafka.Message, partitions ...int) int {
	// 从message的Key中获取roomId
	roomId := string(msg.Key)

	// 使用FNV哈希计算分区索引
	h := fnv.New32a()
	_, _ = h.Write([]byte(roomId))
	partitionIdx := int(h.Sum32()) % len(partitions)

	return partitions[partitionIdx]
}

func (p *KeyPartitioner) RequiresConsistency() bool {
	return true // 保证相同roomId始终进入同一分区
}

func NewKeyPartitioner() *KeyPartitioner {
	return &KeyPartitioner{}
}
