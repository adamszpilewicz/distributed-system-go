package registry

import (
	"bytes"
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
	mutex         *sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	err := r.sendRequiredServices(reg)
	r.notify(patch{
		Added: []patchEntry{
			patchEntry{Name: reg.ServiceName, URL: reg.ServiceURL},
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}
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

func (r *registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqService := range reg.RequiredServices {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					serviceReg.ServiceName, serviceReg.ServiceURL})
			}
		}
	}

	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil

}

func (r *registry) notify(fullPatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, regService := range reg.RequiredServices {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false

				// add service
				for _, added := range fullPatch.Added {
					if added.Name == regService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}

				// remove service
				for _, removed := range fullPatch.Removed {
					if removed.Name == regService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

var reg = registry{
	make([]Registration, 0),
	new(sync.RWMutex),
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
