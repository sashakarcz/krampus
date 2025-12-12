package middleware

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// Decompress middleware handles Content-Encoding: deflate and gzip
// Also handles raw deflate data without Content-Encoding header (Santa does this)
func Decompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		encoding := c.GetHeader("Content-Encoding")
		contentType := c.GetHeader("Content-Type")

		// If no encoding header but Content-Type is JSON, check if body is compressed
		if encoding == "" && strings.Contains(contentType, "json") {
			// Read the body to check if it's compressed
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(400, gin.H{"error": "Failed to read body"})
				return
			}
			c.Request.Body.Close()

			// Try to decompress as raw deflate
			reader := flate.NewReader(bytes.NewReader(bodyBytes))
			decompressed, err := io.ReadAll(reader)
			reader.Close()

			if err == nil && len(decompressed) > 0 {
				// Successfully decompressed, use decompressed data
				c.Request.Body = io.NopCloser(bytes.NewReader(decompressed))
			} else {
				// Not compressed or failed to decompress, use original data
				c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}

			c.Next()
			return
		}

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
