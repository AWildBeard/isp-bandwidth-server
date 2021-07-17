package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
	"strings"
)

var (
	response = []byte(strings.Repeat("0123456789", (100 * 1000) / 8))
)

type sizeReader struct {
	size uint64
}

// This will almost never fill up the buffer all the way, unless a small buffer is chosen.
// In practice, there is almost no difference in speed between this and something that takes all the
// variables and counting into count to fill the buffer with each call. Because the math and looping
// required ends up taking more compute than just a partial buf fill. To make partial buf fills even
// faster, a larger 'response` buffer above is set to ~12,500 B (12KB). There was no difference in testing
// with targeting a 'response' buffer that is the same size or larger as the observed `buf`. This means
// that the majority of latency in response time is well outside our control here. So the 'simpleton'
// implementation was chosen for readability
func (sr *sizeReader) Read(buf []byte) (n int, err error) {
	if sr.size <= 0 { // Exit condition
		return 0, io.EOF
	}

	copySize := len(response)
	if len(response) > len(buf) {
		copySize = len(buf)
	}

	// make sure neither len(response), or len(buf), are greater than sr.size. Else we write too much
	if uint64(copySize) > sr.size {
		copySize = int(sr.size)
	}

	numWritten := copy(buf, response[:copySize])
	sr.size -= uint64(numWritten)
	if sr.size == 0 {
		return numWritten, io.EOF
	} else {
		return numWritten, nil
	}
}

func MbDownloadHandler(c *gin.Context) {
	downloadSize, err := strconv.ParseUint(c.Param("size"), 10, 64)
	if err != nil || downloadSize <=0 || downloadSize > 5000 {
		c.JSON(400, &Error{"Failed to parse download downloadSize from path"})
		fmt.Printf("Error: %v, downloadSize: %v\n", err, downloadSize)
		return
	}

	// Convert downloadSize into it's Mb form
	downloadSize = (downloadSize * 1000 * 1000) / 8

	sizedReader := sizeReader{size: downloadSize}
	c.DataFromReader(200, int64(downloadSize), "text/data", &sizedReader, nil)
}

func MbUploadHandler(c *gin.Context) {
	uploadSize, err := strconv.ParseUint(c.Param("size"), 10, 64)
	if err != nil || uploadSize <=0 || uploadSize > 5000 {
		c.JSON(400, &Error{"Failed to parse download size from path"})
		fmt.Printf("Error: %v, uploadSize: %v\n", err, uploadSize)
		return
	}

	// Convert uploadSize into it's Mb form
	uploadSize = (uploadSize * 1000 * 1000) / 8

	if written, err := io.CopyN(io.Discard, c.Request.Body, int64(uploadSize)); err != nil || written != int64(uploadSize) {
		c.JSON(400, &Error{
			fmt.Sprintf("Failed to read request payload. Failed at %v", written),
		})
		fmt.Printf("Error: %v, uploadSize: %v\n", err, uploadSize)
		return
	}

	c.JSON(200, &UploadStatus{
		Ok: true,
	})
}

func PingHandler(c *gin.Context) {
	c.Data(200, gin.MIMEPlain, []byte("Pong"))
}
