package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3000"
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}

func (r *registry) remove(url string) error {
	for i, obj := range r.registrations {
		if string(obj.ServiceName) == url {
			r.mutex.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("service at url %v not found", url)
}

var reg = registry{
	make([]Registration, 0),
	new(sync.Mutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("received request")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var obj Registration
		err := dec.Decode(&obj)
		if err != nil {
			log.Println("err")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("adding service %v with the url %v", obj.ServiceName, obj.ServiceURL)
		err = reg.add(obj)
		if err != nil {
			log.Println("err")
			w.WriteHeader(http.StatusBadRequest)
		}

	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("err")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		url := string(payload)
		fmt.Printf("removing service with url %v", url)
		reg.remove(url)
		if err != nil {
			log.Println("err")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
