package apiserver

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sakshamsharma/sarga/common/dht"
	"github.com/sakshamsharma/sarga/common/iface"
)

type handleFuncType func(rw http.ResponseWriter, req *http.Request)

// TODO(sakshams): Should have a shutdown channel for integration tests.
func StartAPIServer(args iface.CommonArgs, dht dht.DHT) {
	h := &proxyHandler{dht}
	fs := http.FileServer(http.Dir("static"))

	http.HandleFunc("/sarga/upload/", prefixHandler("/sarga/upload", h.uploadHandler))
	http.HandleFunc("/sarga/files/", prefixHandler("/sarga/files", h.filesHandler))
	http.Handle("/sarga/ui", http.StripPrefix("/sarga/ui/", fs))

	addr := iface.GetAddress(args.IP, args.Port).String()
	log.Println("Listening on", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println(err)
	}
	dht.Shutdown()
}

// prefixHandler creates a handler which receives a stripped req.URL.Path.
func prefixHandler(prefix string, handler handleFuncType) handleFuncType {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.URL.Path = req.URL.Path[len(prefix):]
		handler(rw, req)
	}
}

func (h *proxyHandler) uploadHandler(rw http.ResponseWriter, req *http.Request) {
	// Upload file.
	data, err := ioutil.ReadAll(req.Body)
	if err == nil {
		err = uploadFile(req.URL.Path, data, h.dht)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		} else {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("File uploaded at " + req.URL.Path))
		}
	}
}

func (h *proxyHandler) filesHandler(rw http.ResponseWriter, req *http.Request) {
	// Download file.
	if req.Method == "GET" {
		// Fetch file.
		data, err := downloadFile(req.URL.Path, h.dht)
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(err.Error()))
			log.Println(err)
		} else {
			rw.WriteHeader(http.StatusOK)
			rw.Write(data)
		}
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		_, err := rw.Write([]byte("Unsupported method. Allowed methods: GET"))
		if err != nil {
			log.Println(err)
		}
	}
}

type proxyHandler struct {
	dht dht.DHT
}

func (h *proxyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO(sakshams): Forward regular requests to the internet.
	rw.WriteHeader(http.StatusBadRequest)
	_, err := rw.Write([]byte("Forwarding non-/sarga/{ui,files,upload} requests is not supported yet."))
	if err != nil {
		log.Println(err)
	}
}
