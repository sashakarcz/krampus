package middleware

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// Decompress middleware handles Content-Encoding: deflate and gzip
func Decompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		encoding := c.GetHeader("Content-Encoding")

		if encoding == "" {
			c.Next()
			return
		}

		var reader io.ReadCloser
		var err error

		switch strings.ToLower(encoding) {
		case "deflate":
			reader = flate.NewReader(c.Request.Body)
		case "gzip":
			reader, err = gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(400, gin.H{"error": "Failed to decompress gzip"})
				return
			}
		default:
			// Unsupported encoding, pass through
			c.Next()
			return
		}

		// Replace the request body with decompressed version
		defer c.Request.Body.Close()
		c.Request.Body = reader
		c.Request.Header.Del("Content-Encoding")

		c.Next()
	}
}
