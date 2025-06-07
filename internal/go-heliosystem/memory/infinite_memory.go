package memory

import "time"

type InfiniteMemory struct {
	Segments    map[string]MemorySegment
	Timeline    []TimeStamp
	Compression float64
	GrowthRate  float64
}

type MemorySegment struct {
	ID        string
	Data      []byte
	Timestamp time.Time
	Priority  float64
	Compressed bool
}

type TimeStamp struct {
	Time      time.Time
	SegmentID string
	Event     string
}

func NewInfiniteMemory(compression, growth float64) *InfiniteMemory {
	return &InfiniteMemory{
		Segments:    make(map[string]MemorySegment),
		Timeline:    make([]TimeStamp, 0),
		Compression: compression,
		GrowthRate:  growth,
	}
}

func (im *InfiniteMemory) Store(id string, data []byte, priority float64) {
	segment := MemorySegment{
		ID:        id,
		Data:      data,
		Timestamp: time.Now(),
		Priority:  priority,
		Compressed: false,
	}
	
	im.Segments[id] = segment
	
	timestamp := TimeStamp{
		Time:      segment.Timestamp,
		SegmentID: id,
		Event:     "store",
	}
	im.Timeline = append(im.Timeline, timestamp)
}

func (im *InfiniteMemory) CompressOldSegments() {
	threshold := time.Now().Add(-24 * time.Hour)
	
	for id, segment := range im.Segments {
		if segment.Timestamp.Before(threshold) && !segment.Compressed {
			compressedData := im.compress(segment.Data)
			segment.Data = compressedData
			segment.Compressed = true
			im.Segments[id] = segment
		}
	}
}

func (im *InfiniteMemory) compress(data []byte) []byte {
	compressionRatio := int(float64(len(data)) * im.Compression)
	if compressionRatio < 1 {
		compressionRatio = 1
	}
	return data[:compressionRatio]
}

func (im *InfiniteMemory) Retrieve(id string) *MemorySegment {
	if segment, exists := im.Segments[id]; exists {
		return &segment
	}
	return nil
}
