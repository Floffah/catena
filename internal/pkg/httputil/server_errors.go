package httputil

import (
	"bytes"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxLoggedResponseBodyBytes = 4 * 1024

type errorLoggingWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *errorLoggingWriter) Write(data []byte) (int, error) {
	remaining := maxLoggedResponseBodyBytes - w.body.Len()
	if remaining > 0 {
		if len(data) > remaining {
			w.body.Write(data[:remaining])
		} else {
			w.body.Write(data)
		}
	}

	return w.ResponseWriter.Write(data)
}

// ServerErrorLogger records enough detail to diagnose typed 5xx responses
// without changing the safe error payload sent to clients.
func ServerErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		writer := &errorLoggingWriter{ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		if c.Writer.Status() < 500 {
			return
		}

		body := strings.TrimSpace(writer.body.String())
		if body == "" {
			body = "<empty response body>"
		}

		if len(c.Errors) > 0 {
			log.Printf(
				"server error: method=%s path=%s status=%d errors=%q response=%q",
				c.Request.Method,
				c.Request.URL.RequestURI(),
				c.Writer.Status(),
				c.Errors.String(),
				body,
			)
			return
		}

		log.Printf(
			"server error: method=%s path=%s status=%d response=%q",
			c.Request.Method,
			c.Request.URL.RequestURI(),
			c.Writer.Status(),
			body,
		)
	}
}
