package main

import (
	"GoSimpleProxy/lib"
	"io"
	"log"
	"net/http"
)

func init() {
	if err := lib.ENV.Init(); err != nil {
		lib.LogFatalf("Error init env: %v \n", err)
	}

	lib.SetLogFileName("proxy.log")
	lib.OpenLogFile()
}

func main() {
	lib.LogAppRun(lib.ENV.PROXY_PORT)
	log.Fatal(http.ListenAndServe(":"+lib.ENV.PROXY_PORT, proxy()))
}

func proxy() http.Handler {
	return lib.LogRequest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.URL.Host = lib.ENV.ProxyUrl.Host
		r.URL.Scheme = lib.ENV.ProxyUrl.Scheme

		client := &http.Client{}

		req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)

		//Proxy headers
		for name, value := range r.Header {
			req.Header.Set(name, value[0])
		}

		resp, err := client.Do(req)

		// combined for GET/POST
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		//Donor headers
		for k, v := range resp.Header {
			w.Header().Set(k, v[0])
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

	}))
}
