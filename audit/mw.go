package audit

import (
	"bytes"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

const audiltLogKey = "auditlogcontextkey"

type GinAuditLogCustomize struct {
	ac AuditLogCustomize
}

func (c *GinAuditLogCustomize) Set(alp AuditLogParam, v string) *GinAuditLogCustomize {
	c.ac.Set(alp, v)
	return c
}

func (c *GinAuditLogCustomize) SetUID(v int64) *GinAuditLogCustomize {
	c.ac.SetUID(v)
	return c
}

func (c *GinAuditLogCustomize) Do(ctx *gin.Context) error {
	tmp, ok := ctx.Get(audiltLogKey)
	if !ok {
		return errors.New("No auditlog middleware used")
	}
	alog := tmp.(*AuditLog)
	alog, err := c.ac.Customize(alog)
	if err != nil {
		return err
	}
	ctx.Set(audiltLogKey, alog)
	return nil
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Customize() *GinAuditLogCustomize {
	return &GinAuditLogCustomize{make(map[AuditLogParam]interface{})}
}

var MWAuditlogHandler = func(*AuditLog, *gin.Context) {}

func MWAuditlog(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//如果非四类Method不处理
		if c.Request.Method != "GET" && c.Request.Method != "PUT" && c.Request.Method != "POST" && c.Request.Method != "DELETE" && c.Request.Method != "PATCH" {
			c.Next()
			return
		}
		//Hack writer
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		//Hack request
		aLog := newDefaultAuditLogFromRequest(c.Request)
		aLog.Name = name
		c.Set(audiltLogKey, aLog)
		c.Next()
		aLog = c.MustGet(audiltLogKey).(*AuditLog)
		if aLog.Result == "" {
			aLog.Result = blw.body.String()
		}
		aLog.CreateTime = time.Now()
		if MWAuditlogHandler != nil {
			MWAuditlogHandler(aLog, c)
		}
	}
}
