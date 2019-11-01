package middlwares

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	. "instatasks/helpers"
	"instatasks/models"
	"io/ioutil"
	// "github.com/fatih/color"
	"net"
	"net/http"
	"strconv"
)

const (
	noWritten     = -1
	defaultStatus = 200
)

// Body crypt middleware
func BodyCrypt() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			privks *models.CachedRSAKeys
			wb     *responseBuffer
			ok     bool
			ua     string
		)

		ua = c.GetHeader("User-Agent")

		privks, ok = models.CachedRSAKeysGlobal[ua]

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Uncnown User-Agent"})
			return
		}
		if c.Request.Body != nil {
			bodyRaw, _ := c.GetRawData()
			if dBB, err := DecryptWithPrivateKey(bodyRaw, &privks.CachedRSAPrivateKey); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Decrypt error"})
				return
			} else {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(dBB))
			}
		}

		if w, ok := c.Writer.(gin.ResponseWriter); ok {
			wb = NewResponseBuffer(w)
			c.Writer = wb
			c.Next()
		} else {
			c.Next()
			return
		}

		data := wb.Body.Bytes()
		wb.Body.Reset()

		if data != nil {

			pubk := models.CachedRSAKeysGlobal[ua].CachedRSAPublicKey

			if eBB, err := EncryptWithPublicKey(data, &pubk); err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "Encrypt error"})
				return
			} else {
				wb.Body.Write(eBB)
			}

		}
		wb.Header().Set("Content-Type", "application/encodedjson")
		wb.Header().Set("Content-Length", strconv.Itoa(wb.Body.Len()))
		wb.Flush()
	}
}

type responseBuffer struct {
	Response gin.ResponseWriter // the actual ResponseWriter to flush to
	status   int                // the HTTP response code from WriteHeader
	Body     *bytes.Buffer      // the response content body
	Flushed  bool
}

func NewResponseBuffer(w gin.ResponseWriter) *responseBuffer {
	return &responseBuffer{
		Response: w, status: defaultStatus, Body: &bytes.Buffer{},
	}
}

func (w *responseBuffer) Pusher() http.Pusher {
	return w.Response.Pusher() // use the actual response pusher
}

func (w *responseBuffer) Header() http.Header {
	return w.Response.Header() // use the actual response header
}

func (w *responseBuffer) Write(buf []byte) (int, error) {
	w.Body.Write(buf)
	return len(buf), nil
}

func (w *responseBuffer) WriteString(s string) (n int, err error) {
	n, err = w.Write([]byte(s))
	return
}

func (w *responseBuffer) Written() bool {
	return w.Body.Len() != noWritten
}

func (w *responseBuffer) WriteHeader(status int) {
	w.status = status
}

func (w *responseBuffer) WriteHeaderNow() {
	//if !w.Written() {
	//	w.size = 0
	//	w.ResponseWriter.WriteHeader(w.status)
	//}
}

func (w *responseBuffer) Status() int {
	return w.status
}

func (w *responseBuffer) Size() int {
	return w.Body.Len()
}

func (w *responseBuffer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	//if w.size < 0 {
	//	w.size = 0
	//}
	return w.Response.(http.Hijacker).Hijack()
}

func (w *responseBuffer) CloseNotify() <-chan bool {
	return w.Response.(http.CloseNotifier).CloseNotify()
}

// Fake Flush
// TBD
func (w *responseBuffer) Flush() {
	w.realFlush()
}

func (w *responseBuffer) realFlush() {
	if w.Flushed {
		return
	}
	w.Response.WriteHeader(w.status)
	if w.Body.Len() > 0 {
		_, err := w.Response.Write(w.Body.Bytes())
		if err != nil {
			panic(err)
		}
		w.Body.Reset()
	}
	w.Flushed = true
}
