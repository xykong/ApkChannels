package sign

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func readNextBytes(file *os.File, number int) []byte {
	buf := make([]byte, number)
	_, err := file.Read(buf)
	if err != nil {
		fmt.Println("解码失败", err)
	}
	return buf
}

func ChannelV1(src, dst string) {

	file, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file failed: %v", src)
	}
	defer file.Close()

	offset, err := file.Seek(-22, 2)

	log.Printf("offset: %d, err: %v", offset, err)

	fmt.Println("Success Open File")

	//var mark int64 = 0
	//err = binary.Read(file, binary.BigEndian, &mark)
	//fmt.Printf("%02X ", mark)

	var eocd []byte = []byte{0x50, 0x4b, 0x05, 0x06}

	var buffer bytes.Buffer
	_, _ = io.CopyN(&buffer, file, int64(len(eocd)))
	_bytes := buffer.Bytes()

	for _, b := range _bytes {
		fmt.Printf("%02X ", b)
	}

	fmt.Println()
	if bytes.Compare(eocd, _bytes) == 0 {
		fmt.Println("Equal")
	}
}
