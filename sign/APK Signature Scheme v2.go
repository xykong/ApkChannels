package sign

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

var ZipEocdRecMinSize int64 = 22
var ZipEocdRecSig int64 = 0x06054b50
var ZipEocdCommentLengthFieldOffset int64 = 20

func checkError(err error) {
	if err == nil {
		return
	}

	log.Fatalf("Error: %v", err)
}

func seekEocd(src string) error {

	file, err := os.Open(src)
	checkError(err)

	fi, err := file.Stat()
	checkError(err)

	if fi.Size() < ZipEocdRecMinSize {
		return fmt.Errorf("zip file is invalid")
	}

	// Optimization: 99.99% of APKs have a zero-length comment field in the EoCD record and thus
	// the EoCD record offset is known in advance. Try that offset first to avoid unnecessarily
	// reading more data.
	err = seekEocdWithComment(file, 0)

	return nil
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func seekEocdWithComment(file *os.File, maxCommentSize int64) error {

	// ZIP End of Central Directory (EOCD) record is located at the very end of the ZIP archive.
	// The record can be identified by its 4-byte signature/magic which is located at the very
	// beginning of the record. A complication is that the record is variable-length because of
	// the comment field.
	// The algorithm for locating the ZIP EOCD record is as follows. We search backwards from
	// end of the buffer for the EOCD record signature. Whenever we find a signature, we check
	// the candidate record's comment length is such that the remainder of the record takes up
	// exactly the remaining bytes in the buffer. The search is bounded because the maximum
	// size of the comment field is 65535 bytes because the field is an unsigned 16-bit number.

	if (maxCommentSize < 0) || (maxCommentSize > math.MaxInt16) {
		return fmt.Errorf("maxCommentSize: %v", maxCommentSize)
	}

	fi, err := file.Stat()
	checkError(err)

	if fi.Size() < ZipEocdRecMinSize {
		return fmt.Errorf("zip file is invalid")
	}

	// Lower maxCommentSize if the file is too small.
	maxCommentSize = Min(maxCommentSize, fi.Size()-ZipEocdRecMinSize)

	bufCapacity := ZipEocdRecMinSize + maxCommentSize
	bufOffsetInFile := fi.Size() - bufCapacity

	//ByteBuffer buf = ByteBuffer.allocate(ZIP_EOCD_REC_MIN_SIZE + maxCommentSize)
	//buf.order(ByteOrder.LITTLE_ENDIAN)
	//long bufOffsetInFile = fileSize - buf.capacity();
	buf := new(bytes.Buffer)
	_, _ = file.Seek(bufOffsetInFile, 0)

	_, err = io.CopyN(buf, file, bufCapacity)

	eocdOffsetInBuf := seekEocdWithCommentInBuffer(buf)
	if eocdOffsetInBuf == -1 {
		// No EoCD record found in the buffer
		return nil
	}

	//// EoCD found
	//buf.position(eocdOffsetInBuf)
	//ByteBuffer
	//eocd = buf.slice()
	//eocd.order(ByteOrder.LITTLE_ENDIAN)
	//return Pair.create(eocd, bufOffsetInFile+eocdOffsetInBuf)

	return nil
}

/**
 * Returns the position at which ZIP End of Central Directory record starts in the provided
 * buffer or {@code -1} if the record is not present.
 *
 * <p>NOTE: Byte order of {@code zipContents} must be little-endian.
 */
func seekEocdWithCommentInBuffer(buf *bytes.Buffer) int64 {

	//assertByteOrderLittleEndian(zipContents);

	// ZIP End of Central Directory (EOCD) record is located at the very end of the ZIP archive.
	// The record can be identified by its 4-byte signature/magic which is located at the very
	// beginning of the record. A complication is that the record is variable-length because of
	// the comment field.
	// The algorithm for locating the ZIP EOCD record is as follows. We search backwards from
	// end of the buffer for the EOCD record signature. Whenever we find a signature, we check
	// the candidate record's comment length is such that the remainder of the record takes up
	// exactly the remaining bytes in the buffer. The search is bounded because the maximum
	// size of the comment field is 65535 bytes because the field is an unsigned 16-bit number.

	var archiveSize = int64(buf.Cap())
	if archiveSize < ZipEocdRecMinSize {
		return -1
	}

	maxCommentLength := Min(archiveSize-ZipEocdRecMinSize, math.MaxInt16)
	eocdWithEmptyCommentStartPosition := archiveSize - ZipEocdRecMinSize

	for expectedCommentLength := int64(0); expectedCommentLength < maxCommentLength; expectedCommentLength++ {
		eocdStartPos := eocdWithEmptyCommentStartPosition - expectedCommentLength

		var mark int64 = 0
		var err = binary.Read(bytes.NewBuffer(buf.Bytes()[eocdStartPos:]), binary.LittleEndian, &mark)
		checkError(err)

		if mark == ZipEocdRecSig {

			var actualCommentLength int16 = 0
			var err = binary.Read(bytes.NewBuffer(buf.Bytes()[eocdStartPos+ZipEocdCommentLengthFieldOffset:]), binary.LittleEndian, &actualCommentLength)
			checkError(err)

			if int64(actualCommentLength) == expectedCommentLength {
				return int64(actualCommentLength)
			}
		}
	}

	return -1
}
