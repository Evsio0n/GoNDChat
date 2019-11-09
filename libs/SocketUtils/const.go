package SocketUtils

import (
	"../log"
	"bytes"
	"encoding/binary"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"math"
	"strconv"
)

var uid int
var packageSqense = 4

func commonHeader(command int, length int) []byte {
	var list = UintConverter{numbers: nil}
	list.add(54328, 32)
	list.add(0, 16)
	list.add(uid, 64)
	list.add(0, 32)
	list.add(command, 16)
	list.add(packageSqense, 32)
	list.add(int(math.Min(float64(length), 2048)), 32)
	return list.asUint8List()
}
func commonUnit(value int, radix int) []byte {
	var list = UintConverter{numbers: nil}
	list.add(value, radix)
	return list.asUint8List()
}
func commonString(value string) []byte {
	var list = UintConverter{numbers: nil}
	list.addString(value)
	return list.asUint8List()
}
func commonGroupKey(groupId string, typa int) []byte {
	if typa < 0 || typa > 2 {
		panic("err panic")
	} else {
		var list = UintConverter{numbers: nil}
		list.addString(groupId)
		list.add(typa, 16)
		return list.asUint8List()
	}
}
func packages(command int, data []int, increaseSequence bool) []byte {
	increaseSequence = true
	if len(data) == 0 {
		var header = commonHeader(command, 0)
		if increaseSequence {
			packageSqense++
		}
		var result = BytesCombine(header, IntArraytobytes(data))
		return result
	} else {
		var header = commonHeader(command, len(data))
		if increaseSequence {
			packageSqense++
		}
		result := BytesCombine(header, IntArraytobytes(data))
		return result
	}
}

func getPackageUint(data []int, radix int) int {
	result := 0
	dataInta := IntArraytobytes(data)
	if radix == 8 {
		result = BytesToInt(dataInta, 8)
	} else if radix == 16 {
		result = BytesToInt(dataInta, 16)
	} else if radix == 32 {
		result = BytesToInt(dataInta, 32)
	} else if radix == 64 {
		result = BytesToInt(dataInta, 64)
	}
	return result
}
func getPackageString(data []int) map[string]interface{} {
	cut := data[0:2]
	var result map[string]interface{} = map[string]interface{}{
		"length":  BytesToInt(IntArraytobytes(cut), 16),
		"content": nil,
	}
	var decoder *encoding.Decoder
	var length = BytesToInt(IntArraytobytes(cut), 16)
	decoder = unicode.UTF8.NewDecoder()
	result["content"], _ = decoder.Bytes(IntArraytobytes(data[2 : 2+length]))
	return result
}
func onReceive(event []int) {
	content := event[28]
	status := getPackageUint(event[4:6], 16)
	strings, _ := strconv.ParseInt(string(getPackageUint(event[18:20], 16)), 16, 16)
	command:=
		//TODO command
	squence := getPackageUint(event[20:24], 32)
	length := getPackageUint(event[24:28], 32)
	log.Debug("content\n", content,
		"\nstatus:\n,", status,
		"\ncommond\n", command,
		"\nsquence\n", squence,
		"\nlenth\n", length)
	switch command {
	case 0x75:
		//TODO add package fuction
	}
}

type UintWrapper struct {
	radix int
	value int
}

func (uint *UintWrapper) UintWrapper() {

	if uint.radix%8 == 0 {
		panic("Error at Building UintWrapper")
	}
	log.Debug(uint.radix, uint.value)
}

type UintConverter struct {
	numbers []UintWrapper
}

func (u *UintConverter) addWrapper(wrapper UintWrapper) {
	u.numbers = append(u.numbers, wrapper)
}
func (u *UintConverter) add(value int, radix int) {
	u.addWrapper(UintWrapper{value: value, radix: radix})
}
func (u *UintConverter) addString(value string) {
	bytes2 := []byte(value)
	u.addWrapper(UintWrapper{value: len(bytes2), radix: 16})
	for _, bye := range bytes2 {
		u.add(int(bye), 8)
	}
}
func (u *UintConverter) asUint8List() []byte {
	var size = 0
	for a := 0; a <= len(u.numbers); a++ {
		size += u.numbers[a].radix / 8
	}
	var buffer bytes.Buffer
	var offset = 0
	for a := 0; a <= len(u.numbers); a++ {
		if u.numbers[a].radix == 8 {
			buffer.Write(IntToBytes(u.numbers[a].radix, 8))
		} else if u.numbers[a].radix == 16 {
			buffer.Write(IntToBytes(u.numbers[a].radix, 16))
		} else if u.numbers[a].radix == 32 {
			buffer.Write(IntToBytes(u.numbers[a].radix, 32))
		} else if u.numbers[a].radix == 64 {
			buffer.Write(IntToBytes(u.numbers[a].radix, 64))
		}
		offset += u.numbers[a].radix / 8
	}
	c := make([]byte, offset)
	c = buffer.Bytes()
	return c

}
func IntToBytes(n int, c int) []byte {
	if c == 8 {
		data := uint8(n)
		bytebuf := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytebuf, binary.BigEndian, data)
		return bytebuf.Bytes()
	} else if c == 16 {
		data := uint16(n)
		bytebuf := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytebuf, binary.BigEndian, data)
		return bytebuf.Bytes()
	} else if c == 32 {
		data := uint32(n)
		bytebuf := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytebuf, binary.BigEndian, data)
		return bytebuf.Bytes()
	} else if c == 64 {
		data := uint64(n)
		bytebuf := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytebuf, binary.BigEndian, data)
		return bytebuf.Bytes()
	}
	panic("error")
	//TODO 重写uint包
}

func BytesToInt(bys []byte, c int) int {
	if c == 8 {
		buff := bytes.NewBuffer(bys)
		var data uint8
		_ = binary.Read(buff, binary.BigEndian, &data)
		return int(data)
	} else if c == 16 {
		buff := bytes.NewBuffer(bys)
		var data uint16
		_ = binary.Read(buff, binary.BigEndian, &data)
		return int(data)
	} else if c == 32 {
		buff := bytes.NewBuffer(bys)
		var data uint32
		_ = binary.Read(buff, binary.BigEndian, &data)
		return int(data)
	} else if c == 64 {
		buff := bytes.NewBuffer(bys)
		var data uint64
		_ = binary.Read(buff, binary.BigEndian, &data)
		return int(data)
	}
	panic("error")
}

func BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}
func IntArraytobytes(Slience []int) []byte {
	Slience = make([]int, 0)
	buffer := new(bytes.Buffer)
	_ = binary.Write(buffer, binary.LittleEndian, Slience)
	return buffer.Bytes()
}
