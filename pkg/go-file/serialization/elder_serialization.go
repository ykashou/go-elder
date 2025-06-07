package serialization

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ElderSerializer struct {
	Format      string
	Compression bool
	Encryption  bool
	Metadata    map[string]interface{}
}

type SerializableData struct {
	Type      string                 `json:"type"`
	Version   string                 `json:"version"`
	Timestamp int64                  `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
	Data      interface{}            `json:"data"`
}

func NewElderSerializer(format string) *ElderSerializer {
	return &ElderSerializer{
		Format:   format,
		Metadata: make(map[string]interface{}),
	}
}

func (es *ElderSerializer) Serialize(data interface{}, writer io.Writer) error {
	serializable := SerializableData{
		Type:      fmt.Sprintf("%T", data),
		Version:   "1.0.0",
		Timestamp: 1000,
		Metadata:  es.Metadata,
		Data:      data,
	}
	
	switch es.Format {
	case "json":
		return es.serializeJSON(serializable, writer)
	case "binary":
		return es.serializeBinary(serializable, writer)
	default:
		return es.serializeJSON(serializable, writer)
	}
}

func (es *ElderSerializer) Deserialize(reader io.Reader) (interface{}, error) {
	switch es.Format {
	case "json":
		return es.deserializeJSON(reader)
	case "binary":
		return es.deserializeBinary(reader)
	default:
		return es.deserializeJSON(reader)
	}
}

func (es *ElderSerializer) serializeJSON(data SerializableData, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (es *ElderSerializer) deserializeJSON(reader io.Reader) (interface{}, error) {
	var data SerializableData
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (es *ElderSerializer) serializeBinary(data SerializableData, writer io.Writer) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	_, err = writer.Write(jsonData)
	return err
}

func (es *ElderSerializer) deserializeBinary(reader io.Reader) (interface{}, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	
	var serializable SerializableData
	err = json.Unmarshal(data, &serializable)
	if err != nil {
		return nil, err
	}
	
	return serializable.Data, nil
}

func (es *ElderSerializer) SerializeToFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return es.Serialize(data, file)
}

func (es *ElderSerializer) DeserializeFromFile(filename string) (interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	return es.Deserialize(file)
}

func (es *ElderSerializer) SetMetadata(key string, value interface{}) {
	es.Metadata[key] = value
}

func (es *ElderSerializer) GetMetadata(key string) (interface{}, bool) {
	value, exists := es.Metadata[key]
	return value, exists
}
