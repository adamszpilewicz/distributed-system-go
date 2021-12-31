package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
	if err != nil {
		return err
	}
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to register service"+
			"\nregistry service responded with code: %v", res.StatusCode)
	}
	return nil
}

func ShutdownService(serviceUrl string) error {
	req, err := http.NewRequest(http.MethodDelete,
		ServicesURL, bytes.NewBuffer([]byte(serviceUrl)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != 200 || err != nil {
		return fmt.Errorf("error while trying to remove service"+
			"\nexit status code: %v", res.StatusCode)
	}
	return nil
}
