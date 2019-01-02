# mfh
http handler


## HTTPCall
It is simple    
```golang
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/myfantasy/mfh"
)

func main() {
	d, sc, s, h, c, e := mfh.HTTPCall(http.MethodGet, "http://google.com", nil, nil, time.Millisecond*500, nil)
	fmt.Println(string(d))
	fmt.Println()
	fmt.Println(sc)
	fmt.Println()
	fmt.Println(s)
	fmt.Println()
	fmt.Println(h)
	fmt.Println()
	fmt.Println(c)
	fmt.Println()
	fmt.Println(e)
}

```

## ServeHTTP
Example
```golang

func main() {
	httpPort := 8080

	r := mfh.Route{}

    r.AddDefaultRoute(func(http.ResponseWriter, *http.Request) {

    })
    
    r.AddRoute("ping", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("true"))
	})

	api := http.Server{
		Addr:           fmt.Sprintf(":%d", httpPort),
		Handler:        &r,
		ReadTimeout:    5e9,
		WriteTimeout:   5e9,
		MaxHeaderBytes: 16 << 20, // 16Mb
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Infof("Listen and serve %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Can`t start server; %v", err)

	case <-osSignals:
		log.Infof("Start shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), 5e9)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			log.Infof("Graceful shutdown did not complete in 5s : %v", err)
			if err := api.Close(); err != nil {
				log.Fatalf("Could not stop http server: %v", err)
			}
		}
	}
}
```
