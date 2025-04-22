package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func SaveRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		contentType := c.ContentType()

		if method == "POST" || method == "PUT" || method == "PATCH" {
			if strings.HasPrefix(contentType, "application/json") {
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err == nil {
					c.Set("RequestBody", bodyBytes)
					c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // reset
				}
			} else if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
				_ = c.Request.ParseForm()
				formData := make(map[string]string)
				for key, val := range c.Request.PostForm {
					if len(val) > 0 {
						formData[key] = val[0]
					}
				}
				c.Set("RequestForm", formData)
			} else if strings.HasPrefix(contentType, "multipart/form-data") {
				_ = c.Request.ParseMultipartForm(10 << 20) // up to 10MB
				formData := make(map[string]string)
				for key, val := range c.Request.PostForm {
					if len(val) > 0 {
						formData[key] = val[0]
					}
				}
				c.Set("RequestForm", formData)
			}
		}
		c.Next()
	}
}
