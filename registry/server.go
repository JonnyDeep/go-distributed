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
const ServerURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mu            *sync.RWMutex
}

var reg = registry{
	registrations: make([]Registration, 0),
	mu:            new(sync.RWMutex),
}

/*********************************************
*
*
*添加服务
*
*
***********************************************/
func (r *registry) add(re Registration) error {
	r.mu.RLock()
	r.registrations = append(r.registrations, re)
	r.mu.RUnlock()
	err := r.sendRequiredServices(re)
	r.notify(patch{
		Add: []patchEntry{
			{
				Name: re.ServiceName,
				URL:  re.ServiceUrl,
			},
		},
	})
	return err
}

func (r *registry) notify(fullPatch patch) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, reqService := range reg.RequiredService {
				p := patch{Add: []patchEntry{}, Remove: []patchEntry{}}
				sendUpdate := false
				for _, added := range fullPatch.Add {
					if added.Name == reqService {
						p.Add = append(p.Add, added)
						sendUpdate = true
					}
				}

				for _, removeed := range fullPatch.Remove {
					if removeed.Name == reqService {
						p.Remove = append(p.Remove, removeed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPath(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

func (r *registry) sendRequiredServices(re Registration) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqService := range re.RequiredService {
			if serviceReg.ServiceName == reqService {
				p.Add = append(p.Add, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceUrl,
				})
			}
		}
	}
	err := r.sendPath(p, re.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) sendPath(p patch, url string) error {
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

func (r *registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceUrl == url {
			r.notify(patch{
				Remove: []patchEntry{
					{
						Name: r.registrations[i].ServiceName,
						URL:  r.registrations[i].ServiceUrl,
					},
				},
			})
			r.mu.Lock()
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			r.mu.Unlock()
			return nil
		}
	}
	return fmt.Errorf("Service at url %s not found", url)
}

func (r *registry) hearbeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, reg := range r.registrations {
			wg.Add(1)
			go func(reg Registration) {
				defer wg.Done()
				success := true
				for attemps := 0; attemps < 3; attemps++ {
					res, err := http.Get(reg.HeartBeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("Heartbeat check passed for %v", reg.ServiceName)
						if !success {
							r.add(reg)
						}
						break
					}
					log.Printf("Heartbeat check failed for %v", reg.ServiceName)
					if success {
						success = false
						r.remove(reg.ServiceUrl)
					}
					time.Sleep(1 * time.Second)

				}
			}(reg)
			wg.Wait()
		}
		time.Sleep(freq)
	}
}

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.hearbeat(3 * time.Second)
	})
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Add servce:%v with %s\n", r.ServiceName, r.ServiceUrl)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		url := string(payload)
		log.Printf("Removing service at url:%s", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
