package httprouter

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "Hello Httprouter")
	})

	server := http.Server{
		Handler: router,
		Addr:    "localhost:2323",
	}
	server.ListenAndServe()
}
