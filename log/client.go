package log

import (
	"bytes"
	"fmt"
	"net/http"

	"huangliqun.github.com/registry"

	stdlog "log"
)

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stdlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stdlog.SetFlags(0)
	stdlog.SetOutput(&clientLogger{url: serviceURL})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to send long message. Service responded wrong %v", res.StatusCode)
	}
	return len(data), nil
}
