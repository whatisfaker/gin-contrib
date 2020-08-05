package audit

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//AuditLog 审计日志结构
type AuditLog struct {
	//操作名
	Name string `json:"name"`
	//访问的接口地址
	Path string `json:"path"`
	//操作 (GET,POST,PUT)
	Operation string `json:"operation"`
	//条件 （请求内容，条件）
	Condition string `json:"condition"`
	//操作结果
	Result string `json:"result"`

	UID int64 `json:"uid"`
	//操作时间
	CreateTime time.Time `json:"create_time"`
	//真实IP
	RemoteAddr string `json:"remote_addr"`
	//RealIP
	RealIP string `json:"real_ip"`

	//额外信息
	Ext1    string `json:"ext1"`
	Ext2    string `json:"ext2"`
	ExtInt1 int    `json:"ext_int1"`
	ExtInt2 int    `json:"ext_int2"`
}

func newDefaultAuditLogFromRequest(req *http.Request) *AuditLog {
	aLog := new(AuditLog)
	if req.Method == "GET" {
		aLog.Condition = req.URL.RawQuery
	} else if req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" || req.Method == "PATCH" {
		ct := req.Header.Get("Content-Type")
		//非文件上传
		if !strings.Contains(ct, "multipart/form-data") {
			b, err := ioutil.ReadAll(req.Body)
			//没有异常要写回body
			if err == nil {
				aLog.Condition = ct + " " + string(b)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}
		} else {
			err := req.ParseMultipartForm(0)
			if err == nil {
				b, err := json.Marshal(req.PostForm)
				if err == nil {
					aLog.Condition = "multipart/form-data " + string(b)
				}
			}
		}
	}
	aLog.Operation = req.Method
	aLog.Path = req.URL.Path
	idx := strings.LastIndex(req.RemoteAddr, ":")
	if idx > 0 {
		aLog.RemoteAddr = req.RemoteAddr[:idx]
	} else {
		aLog.RemoteAddr = req.RemoteAddr
	}
	aLog.RealIP = req.Header.Get("X-Real-Ip")
	if aLog.RealIP == "" {
		aLog.RealIP = req.Header.Get("X-Forwarded-For")
	}
	if aLog.RealIP == "" {
		aLog.RealIP = req.RemoteAddr
	}
	return aLog
}

type AuditLogCustomize map[AuditLogParam]interface{}
type AuditLogParam string

const (
	AuditLogCondition AuditLogParam = "al_condition"
	AuditLogResult    AuditLogParam = "al_result"
	AuditLogUser      AuditLogParam = "al_user"
	AuditLogExtInt1   AuditLogParam = "al_ext_int1"
	AuditLogExtInt2   AuditLogParam = "al_ext_int2"
	AuditLogExt1      AuditLogParam = "al_ext_1"
	AuditLogExt2      AuditLogParam = "al_ext_2"
)

//@Deprecated
// func (c AuditLogCustomize) Set(k AuditLogParam, v string) AuditLogCustomize {
// 	c[k] = v
// 	return c
// }

func (c AuditLogCustomize) SetCondition(v string) AuditLogCustomize {
	c[AuditLogCondition] = v
	return c
}

func (c AuditLogCustomize) SetResult(v string) AuditLogCustomize {
	c[AuditLogResult] = v
	return c
}

func (c AuditLogCustomize) SetUID(uid int64) AuditLogCustomize {
	c[AuditLogUser] = uid
	return c
}

func (c AuditLogCustomize) SetExt1(ext string) AuditLogCustomize {
	c[AuditLogExt1] = ext
	return c
}

func (c AuditLogCustomize) SetExt2(ext string) AuditLogCustomize {
	c[AuditLogExt2] = ext
	return c
}

func (c AuditLogCustomize) SetExtID1(id int) AuditLogCustomize {
	c[AuditLogExtInt1] = id
	return c
}

func (c AuditLogCustomize) SetExtID2(id int) AuditLogCustomize {
	c[AuditLogExtInt2] = id
	return c
}

func (c AuditLogCustomize) Customize(al *AuditLog) (*AuditLog, error) {
	for k, v := range c {
		switch k {
		case AuditLogCondition:
			al.Condition = v.(string)
		case AuditLogResult:
			al.Result = v.(string)
		case AuditLogExt1:
			al.Ext1 = v.(string)
		case AuditLogExt2:
			al.Ext2 = v.(string)
		case AuditLogExtInt1:
			al.ExtInt1 = v.(int)
		case AuditLogExtInt2:
			al.ExtInt2 = v.(int)
		case AuditLogUser:
			al.UID = v.(int64)
		}
	}
	return al, nil
}
