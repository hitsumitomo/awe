package test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"math/rand"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hitsumitomo/awe"
)

var p = awe.Packet{
	Op:    "PUB",
	Key:   "topic.person.say",
	Value: []byte("Why, sometimes I've believed as many as six impossible things before breakfast."),
}
func TestEncodeDecodePacket(t *testing.T) {
	data := p.Encode()

	var result awe.Packet
	err := result.Decode(data)
	if err != nil {
		t.Errorf("Decode() error = %v", err)
		return
	}

	if !reflect.DeepEqual(p, result) {
		t.Errorf("Encopde/Decode mismatch: expected %v, got %v", p, result)
	}
}

type Example struct {
    Int8Field    int8
    Uint8Field   uint8
    Int16Field   int16
    Uint16Field  uint16
    Int32Field   int32
    Uint32Field  uint32
    Int64Field   int64
    Uint64Field  uint64
    ByteField    byte
    Float32Field float32
    Float64Field float64
    StringField  string
    BytesField   []byte
    ArrayField   [64]byte
    Bool         bool
}

func TestMarshalUnmarshal(t *testing.T) {
	tests := []struct {
			name string
			example Example
	}{
	{
		name: "zero_values",
		example: Example{},
	},
	{
		name: "positive_values",
		example: Example{
			Int8Field:    2,
			Uint8Field:   3,
			Int16Field:   4,
			Uint16Field:  5,
			Int32Field:   6,
			Uint32Field:  7,
			Int64Field:   8,
			Uint64Field:  9,
			ByteField:    10,
			Float32Field: 11.1,
			Float64Field: 12.2,
			StringField:  "test",
			BytesField:   []byte{13, 14, 15},
			ArrayField:   [64]byte{16, 17, 18},
			Bool:         true,
		},
	},
	{
		name: "negative_values",
		example: Example{
			Int8Field:    -2,
			Int16Field:   -4,
			Int32Field:   -6,
			Int64Field:   -8,
			Float32Field: -11.1,
			Float64Field: -12.2,
			StringField:  "negative",
			BytesField:   []byte{13, 14, 15},
			ArrayField:   [64]byte{16, 17, 18},
			Bool:         false,
		},
	},
	{
		name: "mixed_values",
		example: Example{
			Int8Field:    -2,
			Uint8Field:   3,
			Int16Field:   -4,
			Uint16Field:  5,
			Int32Field:   -6,
			Uint32Field:  7,
			Int64Field:   -8,
			Uint64Field:  9,
			ByteField:    10,
			Float32Field: -11.1,
			Float64Field: 12.2,
			StringField:  "mixed",
			BytesField:   []byte{13, 14, 15},
			ArrayField:   [64]byte{16, 17, 18},
			Bool:         true,
		},
	},}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.example.Marshal()
			if err != nil {
				t.Errorf("Marshal() error = %v", err)
				return
			}

			var result Example
			err = result.Unmarshal(data)
			if err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}

			if !reflect.DeepEqual(tt.example, result) {
				t.Errorf("Marshal/Unmarshal mismatch: expected %v, got %v", tt.example, result)
			}
		})
	}
}

var example = Example{
	Int8Field:    -2,
	Uint8Field:   3,
	Int16Field:   -4,
	Uint16Field:  5,
	Int32Field:   -6,
	Uint32Field:  7,
	Int64Field:   -8,
	Uint64Field:  9,
	ByteField:    10,
	Float32Field: -11.1,
	Float64Field: 12.2,
	StringField:  "mixed",
	BytesField:   []byte{13, 14, 15},
	ArrayField:   [64]byte{16, 17, 18},
	Bool:         true,
}

func TestMarshalUnmarshalCrc32(t *testing.T) {
	data, err := example.Marshal(awe.Crc32)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}

	var result Example
	err = result.Unmarshal(data)
	if err != nil {
		t.Errorf("Unmarshal() error = %v", err)
		return
	}

	if !reflect.DeepEqual(example, result) {
		t.Errorf("Marshal/Unmarshal mismatch: expected %v, got %v", example, result)
	}
}

func TestMarshalUnmarshalCompress(t *testing.T) {
	data, err := example.Marshal(awe.Compress)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}

	var result Example
	err = result.Unmarshal(data)
	if err != nil {
		t.Errorf("Unmarshal() error = %v", err)
		return
	}

	if !reflect.DeepEqual(example, result) {
		t.Errorf("Marshal/Unmarshal mismatch: expected %v, got %v", example, result)
	}
}

func TestMarshalUnmarshalCompressCrc32(t *testing.T) {
	data, err := example.Marshal(awe.Compress | awe.Crc32)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}

	var result Example
	err = result.Unmarshal(data)
	if err != nil {
		t.Errorf("Unmarshal() error = %v", err)
		return
	}

	if !reflect.DeepEqual(example, result) {
		t.Errorf("Marshal/Unmarshal mismatch: expected %v, got %v", example, result)
	}
}

func TestMarshalUnmarshalCompressCrc32FastUnmarshal(t *testing.T) {
	data, err := example.Marshal(awe.Compress | awe.Crc32 | awe.FastUnmarshal)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}

	var result Example
	err = result.Unmarshal(data)
	if err != nil {
		t.Errorf("Unmarshal() error = %v", err)
		return
	}

	if !reflect.DeepEqual(example, result) {
		t.Errorf("Marshal/Unmarshal mismatch: expected %v, got %v", example, result)
	}
}

func TestMarshalUnmarshalFastUnmarshal(t *testing.T) {
	data, err := example.Marshal(awe.FastUnmarshal)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}

	var result Example
	err = result.Unmarshal(data)
	if err != nil {
		t.Errorf("Unmarshal() error = %v", err)
		return
	}

	if !reflect.DeepEqual(example, result) {
		t.Errorf("Marshal/Unmarshal mismatch: expected %v, got %v", example, result)
	}
}

func BenchmarkPacketEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p.Encode()
	}
}

func BenchmarkPacketDecode(b *testing.B) {
	data := p.Encode()
	for i := 0; i < b.N; i++ {
		p.Decode(data)
	}
}

func BenchmarkMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := example.Marshal()
		if err != nil {
			b.Fatalf("MarshalExample failed: %v", err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	data, err := example.Marshal()
	if err != nil {
		b.Fatalf("MarshalExample failed: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := example.Unmarshal(data)
		if err != nil {
			b.Fatalf("UnmarshalExample failed: %v", err)
		}
	}
}

func BenchmarkUnmarshalFastUnmarshal(b *testing.B) {
	data, err := example.Marshal(awe.FastUnmarshal)
	if err != nil {
		b.Fatalf("MarshalExample failed: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := example.Unmarshal(data)
		if err != nil {
			b.Fatalf("UnmarshalExample failed: %v", err)
		}
	}
}

func BenchmarkMarshalCRC32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := example.Marshal(awe.Crc32)
		if err != nil {
			b.Fatalf("MarshalExample failed: %v", err)
		}
	}
}

func BenchmarkUnmarshalCRC32(b *testing.B) {
	data, err := example.Marshal(awe.Crc32)
	if err != nil {
		b.Fatalf("MarshalExample failed: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := example.Unmarshal(data)
		if err != nil {
			b.Fatalf("UnmarshalExample failed: %v", err)
		}
	}
}

func BenchmarkMarshalCompressed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := example.Marshal(awe.Compress)
		if err != nil {
			b.Fatalf("MarshalExample failed: %v", err)
		}
	}
}

func BenchmarkUnmarshalCompressed(b *testing.B) {
	data, err := example.Marshal(awe.Compress)
	if err != nil {
		b.Fatalf("MarshalExample failed: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := example.Unmarshal(data)
		if err != nil {
			b.Fatalf("UnmarshalExample failed: %v", err)
		}
	}
}

func randomExampleProto() *ExampleProto {
	return &ExampleProto{
		Int8Field:    int32(rand.Intn(128) - 64),
		Uint8Field:   uint32(rand.Intn(256)),
		Int16Field:   int32(rand.Intn(32768) - 16384),
		Uint16Field:  uint32(rand.Intn(65536)),
		Int32Field:   rand.Int31(),
		Uint32Field:  rand.Uint32(),
		Int64Field:   rand.Int63(),
		Uint64Field:  rand.Uint64(),
		ByteField:    uint32(rand.Intn(256)),
		Float32Field: rand.Float32(),
		Float64Field: rand.Float64(),
		StringField:  "Benchmark String Field",
		BytesField:   []byte("If you don’t know where you want to go, then it doesn’t matter which path you take."),
		ArrayField:   make([]byte, 64),
		BoolField:    rand.Intn(2) == 1,
	}
}

func BenchmarkProtoMarshal(b *testing.B) {
	ex := randomExampleProto()
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(ex)
		if err != nil {
			b.Fatalf("Failed to marshal: %v", err)
		}
	}
}

func BenchmarkProtoUnmarshal(b *testing.B) {
	ex := randomExampleProto()
	data, err := proto.Marshal(ex)
	if err != nil {
		b.Fatalf("Failed to marshal for unmarshal test: %v", err)
	}

	for i := 0; i < b.N; i++ {
		var decoded ExampleProto
		if err := proto.Unmarshal(data, &decoded); err != nil {
			b.Fatalf("Failed to unmarshal: %v", err)
		}
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(example)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	data2, _ := json.Marshal(example)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data2, &example); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGobMarshal(b *testing.B) {
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(example); err != nil {
			b.Fatal(err)
		}
		buf.Reset()
	}
}

func BenchmarkGobUnmarshal(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(example); err != nil && err != io.EOF {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dec := gob.NewDecoder(&buf)
		if err := dec.Decode(&example); err != nil && err != io.EOF {
			b.Fatal(err)
		}
        buf.Reset()
    }
}

// syntax = "proto3";

// package main;
// option go_package = "./";

// message ExampleProto {
//     int32 int8_field = 1;       // int8 in Go is an int32 in protobuf
//     uint32 uint8_field = 2;     // uint8 in Go is uint32 in protobuf
//     int32 int16_field = 3;      // int16 in Go is an int32 in protobuf
//     uint32 uint16_field = 4;    // uint16 in Go is uint32 in protobuf
//     int32 int32_field = 5;      // directly maps to int32
//     uint32 uint32_field = 6;    // directly maps to uint32
//     int64 int64_field = 7;      // directly maps to int64
//     uint64 uint64_field = 8;    // directly maps to uint64
//     uint32 byte_field = 9;      // byte is an alias for uint8 in Go, use uint32 here
//     float float32_field = 10;   // directly maps to float32
//     double float64_field = 11;  // directly maps to float64
//     string string_field = 12;   // directly maps to string
//     bytes bytes_field = 13;     // []byte is represented as bytes
//     bytes array_field = 14;     // array of 64 bytes, also represented as bytes
//     bool bool_field = 15;       // directly maps to bool
// }
