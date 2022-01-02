package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
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
			{Name: reg.ServiceName, URL: reg.ServiceURL},
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *registry) remove(url string) error {
	for i, reg := range r.registrations {
		if reg.ServiceURL == url {
			log.Println("notifying about removal", reg.ServiceName, reg.ServiceURL)
			r.notify(patch{
				Removed: []patchEntry{
					{Name: reg.ServiceName, URL: reg.ServiceURL},
				},
			})
			r.mutex.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
		}
		return nil
	}
	return fmt.Errorf("service at url %v not found", url)
}

func (r *registry) heartbeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, reg := range r.registrations {
			wg.Add(1)
			go func(reg Registration) {
				defer wg.Done()
				for attemps := 0; attemps < 3; attemps++ {
					//success := true
					res, err := http.Get(reg.HeartbeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("Hearbeat passed for %v\n", reg.ServiceName)

						//if !success {
						//	r.add(reg)
						//}
						break
					}
					log.Printf("Hearbeat failed for %v\n", reg.ServiceName)
					//success = false
					r.remove(reg.ServiceURL)
					time.Sleep(1 * time.Second)
				}
			}(reg)
			wg.Wait()
			time.Sleep(freq)
		}
	}
}

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.heartbeat(3 * time.Second)
	})
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
		err = reg.remove(url)
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
