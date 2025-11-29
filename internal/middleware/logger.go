package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestLogger creates a middleware that logs detailed information about requests and responses
func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore the body for downstream handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create custom response writer to capture response
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// Log incoming request
		logger.WithFields(logrus.Fields{
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"client_ip":    c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"referer":      c.Request.Referer(),
			"request_id":   c.GetHeader("X-Request-UserID"),
			"content_type": c.ContentType(),
		}).Info("Incoming request")

		// Log request body if it exists (be careful with sensitive data)
		if len(requestBody) > 0 && len(requestBody) < 10000 { // Limit body size to log
			// Only log body for certain content types to avoid logging binary data
			contentType := c.ContentType()
			if contentType == "application/json" || contentType == "application/x-www-form-urlencoded" {
				logger.WithFields(logrus.Fields{
					"path":         c.Request.URL.Path,
					"request_body": string(requestBody),
				}).Debug("Request body")
			}
		}

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		// Determine log level based on status code
		logLevel := logrus.InfoLevel
		if statusCode >= 500 {
			logLevel = logrus.ErrorLevel
		} else if statusCode >= 400 {
			logLevel = logrus.WarnLevel
		}

		// Log response
		responseFields := logrus.Fields{
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"status":        statusCode,
			"duration_ms":   duration.Milliseconds(),
			"client_ip":     c.ClientIP(),
			"response_size": blw.body.Len(),
		}

		// Add error information if present
		if len(c.Errors) > 0 {
			responseFields["errors"] = c.Errors.String()
		}

		// Log response body for debugging (only for non-2xx responses and if not too large)
		if (statusCode < 200 || statusCode >= 300) && blw.body.Len() > 0 && blw.body.Len() < 10000 {
			responseFields["response_body"] = blw.body.String()
		}

		logger.WithFields(responseFields).Log(logLevel, "Request completed")
	}
}

// SimpleRequestLogger creates a simpler version of request logger with less detail
func SimpleRequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		logger.WithFields(logrus.Fields{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"client_ip": c.ClientIP(),
		}).Info("→ Incoming request")

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		logLevel := logrus.InfoLevel
		if statusCode >= 500 {
			logLevel = logrus.ErrorLevel
		} else if statusCode >= 400 {
			logLevel = logrus.WarnLevel
		}

		logger.WithFields(logrus.Fields{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status":      statusCode,
			"duration_ms": duration.Milliseconds(),
		}).Log(logLevel, "← Request completed")
	}
}
