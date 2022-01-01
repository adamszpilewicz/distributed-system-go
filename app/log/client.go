package log

import (
	"bytes"
	"fmt"
	"github.com/adamszpilewicz/distributed-systems/app/registry"
	stlog "log"
	"net/http"
)

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != 200 {
		return 0, fmt.Errorf("failed to send log message; service "+
			"responded with status code: %v", res.StatusCode)
	}
	return len(data), nil
}

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: serviceURL})
}
