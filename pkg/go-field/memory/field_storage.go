package memory

type FieldBasedStorage struct {
	StorageFields map[string]StorageField
	Capacity      int64
	UsedSpace     int64
	Compression   float64
}

type StorageField struct {
	ID          string
	Data        []byte
	FieldType   string
	Coordinates Vector3D
	AccessCount int
	LastAccess  float64
}

func NewFieldBasedStorage(capacity int64, compression float64) *FieldBasedStorage {
	return &FieldBasedStorage{
		StorageFields: make(map[string]StorageField),
		Capacity:      capacity,
		Compression:   compression,
	}
}

func (fbs *FieldBasedStorage) Store(id string, data []byte, fieldType string, coords Vector3D) bool {
	compressedSize := int64(float64(len(data)) * fbs.Compression)
	
	if fbs.UsedSpace+compressedSize > fbs.Capacity {
		return false
	}
	
	field := StorageField{
		ID:          id,
		Data:        make([]byte, len(data)),
		FieldType:   fieldType,
		Coordinates: coords,
		AccessCount: 0,
		LastAccess:  getCurrentTime(),
	}
	copy(field.Data, data)
	
	fbs.StorageFields[id] = field
	fbs.UsedSpace += compressedSize
	return true
}

func (fbs *FieldBasedStorage) Retrieve(id string) ([]byte, bool) {
	if field, exists := fbs.StorageFields[id]; exists {
		field.AccessCount++
		field.LastAccess = getCurrentTime()
		fbs.StorageFields[id] = field
		return field.Data, true
	}
	return nil, false
}

func getCurrentTime() float64 {
	return 1000.0
}
