package sign

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

var eocd int64 = 0x6054B50

func V1Writer(src, dst string, channel string) error {

	if channel == "" {
		return fmt.Errorf("channel is empty")
	}

	if src == dst || len(dst) == 0 {
		return V2WriteInPlace(src, channel)
	} else {
		return V2WriteStream(src, dst, channel)
	}
}

func seekEOCD(file *os.File) (offset int64, err error) {

	// -22 is the EOCD offset of zip file with empty comment
	offset, err = file.Seek(-22, 2)
	if err != nil {
		return 0, fmt.Errorf("invalid zip file: %v", file.Name())
	}

	var mark int64 = 0
	err = binary.Read(file, binary.LittleEndian, &mark)

	if eocd != mark {
		return 0, fmt.Errorf(".zip file with commnet not support: %v", file.Name())
	}

	return
}

func V2WriteInPlace(src string, channel string) error {

	file, err := os.OpenFile(src, os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file failed: %v", src)
	}
	defer file.Close()

	_, err = seekEOCD(file)
	if err != nil {
		return err
	}

	// .ZIP file comment length
	size := int16(len(channel))

	_, err = file.Seek(-2, 2)
	err = binary.Write(file, binary.LittleEndian, size)
	if err != nil {
		return fmt.Errorf("binary.Write failed: %v", err)
	}

	_, err = file.WriteString(channel)
	if err != nil {
		return fmt.Errorf("file.WriteString failed: %v", err)
	}

	return nil
}

func V2WriteStream(src, dst string, channel string) error {

	file, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file failed: %v", src)
	}
	defer file.Close()

	offset, err := seekEOCD(file)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create file failed: %v", dst)
	}
	defer dstFile.Close()

	stat, _ := file.Stat()
	log.Printf("offset: %v, total: %v", offset, stat)

	_, _ = file.Seek(0, 0)
	_, err = io.CopyN(dstFile, file, offset+20)
	//_, err = io.Copy(dstFile, file)
	if err != nil {
		return fmt.Errorf("io.CopyN file %v to %v failed: %v", src, dst, err)
	}

	// .ZIP file comment length
	size := int16(len(channel))

	err = binary.Write(dstFile, binary.LittleEndian, size)
	if err != nil {
		return fmt.Errorf("binary.Write %v failed: %v", dst, err)
	}

	_, err = dstFile.WriteString(channel)
	if err != nil {
		return fmt.Errorf("file.WriteString %v failed: %v", dst, err)
	}

	return nil
}
