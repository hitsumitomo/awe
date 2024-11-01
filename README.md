# awe

`awe` - Ultrafast Go Marshalling Library

## Features

- Significantly outperforms standard Go serialization methods.
- Offers customizable encoding options via bitmap flags.
- Includes `awec` compiler for automatic generation of `Marshal` and `Unmarshal` methods for your structures.

## Supported Data Types

- Integer types: `int8`, `uint8`, `int16`, `uint16`, `int32`, `uint32`, `int64`, `uint64`
- Basic types: `byte`, `bool`, `string`
- Floating-point types: `float32`, `float64`
- Byte slices `[]byte` and fixed-size byte arrays `[NNN]byte`
- **Note**: Nested structs are not supported. Only top-level fields with supported types can be used.

## Building the awec Compiler Tool
To build the awec compiler tool, run:
```sh
go build awec.go
```

Place awec in your path, for example:
```sh
sudo mv awec /usr/local/bin
```

Run awec on the file containing structs to generate `Marshal` and `Unmarshal` methods:
```sh
awec example.go
```
awec generates `example_marshalling.go` where `Marshal`/`Unmarshal` methods are stored.

## Packet Structure

The package also provides a `Packet` struct designed for efficient encoding and decoding operations, useful for network communication or data storage.

```go
type Packet struct {
    Op    string
    Key   string
    Value []byte
}
```

Example usage:
```go
package main
import (
    "log"
    "github.com/hitsumitomo/awe"
)

func main() {
    packet := awe.Packet{
        Op:    "SET",
        Key:   "example_key",
        Value: []byte("example_value"),
    }

    // encode Packet to []byte
    encoded := packet.Encode()

    // decode Packet from []byte
    var decodedPacket awe.Packet
    if err := decodedPacket.Decode(encoded); err != nil {
        log.Fatal("Decoding failed:", err)
    }
}
```

## Example Struct

```go
package main

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
// another structures or code...
```
Executing `awec file.go` creates `file_marshalling.go`:
<pre>
Parsing file: file.go
Marshal, Unmarshal methods for struct "Example" successfully generated.
</pre>
- `func (e *Example) Marshal(flags ...byte) ([]byte, error)`: Serializes the `Example` struct into a binary format.
- `func (e *Example) Unmarshal(data []byte) error`: Deserializes a binary format into an `Example` struct.

## Basic Example

```go
package main

import (
    "log"
    "github.com/hitsumitomo/awe"
)

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

func main() {
    example := Example{
        Int8Field:    1,
        Uint8Field:   2,
        Int16Field:   -3,
        Uint16Field:  4,
        Int32Field:   5,
        Uint32Field:  6,
        Int64Field:   -7,
        Uint64Field:  8,
        ByteField:    'a',
        Float32Field: 3.14,
        Float64Field: -2.718,
        StringField:  "example",
        BytesField:   []byte("data"),
        ArrayField:   [64]byte{'t', 'e', 's', 't'},
        Bool:         true,
    }

    // Ensure that awec has already generated the marshalling methods.
    encoded, err := example.Marshal(awe.FastUnmarshal)
    if err != nil {
        log.Fatal("Marshalling failed:", err)
    }

    var decodedExample Example
    if err := decodedExample.Unmarshal(encoded); err != nil {
        log.Fatal("Unmarshalling failed:", err)
    }
}
```
By default, the binary format is structured as follows: `[marker][length][payload]`.

Bitmap flags can be used to include additional features:
- `[marker][length][crc32][payload]`  - Adds a CRC32 checksum for error detection.
- `[marker][length][gzipped payload]` - Compresses the payload using gzip for efficient storage and transmission.

## Additional Bitmap Flags

Example: `example.Marshal(flag)` (stored in marker)
- `awe.FastUnmarshal` - Use packet data directly, reducing allocations.
- `awe.Crc32` - Add CRC32 checksum.
- `awe.Pub`,`awe.Sub`,`awe.Unsub`,`awe.Queue` - Message broker operations.
- `awe.Compress` - Apply Gzip compression.

```go
example.Marshal(awe.Crc32 | awe.Pub)
example.Marshal(awe.Compress)
example.Marshal(awe.FastUnmarshal)
```

## Benchmarks
<pre>
BenchmarkPacketEncode-4              12939886        91.92 ns/op          112 B/op          1 allocs/op
BenchmarkPacketDecode-4              13830751        84.96 ns/op           19 B/op          2 allocs/op
BenchmarkMarshal-4                    9811832       118.7 ns/op           144 B/op          1 allocs/op
BenchmarkUnmarshal-4                 12128236        98.03 ns/op            8 B/op          2 allocs/op
BenchmarkUnmarshalFastUnmarshal-4    16830313        69.58 ns/op            5 B/op          1 allocs/op
BenchmarkMarshalCRC32-4               5802518       201.8 ns/op           144 B/op          1 allocs/op
BenchmarkUnmarshalCRC32-4             6515298       182.2 ns/op             8 B/op          2 allocs/op
BenchmarkMarshalCompressed-4             7299    149795 ns/op          814240 B/op         21 allocs/op
BenchmarkUnmarshalCompressed-4          87684     12635 ns/op           41800 B/op         10 allocs/op
BenchmarkProtoMarshal-4               1818622       662.4 ns/op           256 B/op          1 allocs/op
BenchmarkProtoUnmarshal-4             1413729       817.4 ns/op           360 B/op          4 allocs/op
BenchmarkJSONMarshal-4                 300334      3820 ns/op             544 B/op          2 allocs/op
BenchmarkJSONUnmarshal-4                66055     17996 ns/op             240 B/op          7 allocs/op
BenchmarkGobMarshal-4                  110130     10704 ns/op            1904 B/op         49 allocs/op
BenchmarkGobUnmarshal-4               2808040       417.8 ns/op           304 B/op          5 allocs/op
</pre>