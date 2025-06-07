package serialization

import (
	"encoding/json"
	"fmt"
)

type ElderSerializer struct {
	Format string
	Version string
}

type SerializedElder struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Data     map[string]interface{} `json:"data"`
	Metadata map[string]string      `json:"metadata"`
}

func (es *ElderSerializer) Serialize(entity interface{}) ([]byte, error) {
	serialized := SerializedElder{
		Type:     "elder_entity",
		Data:     make(map[string]interface{}),
		Metadata: map[string]string{"format": es.Format, "version": es.Version},
	}
	
	return json.Marshal(serialized)
}

func (es *ElderSerializer) Deserialize(data []byte) (interface{}, error) {
	var serialized SerializedElder
	err := json.Unmarshal(data, &serialized)
	if err != nil {
		return nil, fmt.Errorf("deserialization failed: %v", err)
	}
	
	return serialized, nil
}
