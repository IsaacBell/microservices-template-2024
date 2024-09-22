package gateway

import (
	"net/http"
	"net/url"
	"os"

	"github.com/vulcand/oxy/v2/buffer"
	"github.com/vulcand/oxy/v2/forward"
	"github.com/vulcand/oxy/v2/roundrobin"
)

// Forwards incoming requests to whatever location a URL points to, adding proper forwarding headers
func RProxy(port int) (*http.Server, error) {
	fwd := forward.New(false)
	lb, _ := roundrobin.New(fwd)

	// buffer will read the request body and will replay the request again in case if forward returned status
	// corresponding to nework error (e.g. Gateway Timeout)
	buffer, _ := buffer.New(lb, buffer.Retry(`IsNetworkError() && Attempts() < 2`))

	coreUrl := url.URL{Host: os.Getenv("CORE_SERVICE_ADDRESS")}
	lb.UpsertServer(&coreUrl)
	// lb.UpsertServer(url2)

	s := &http.Server{
		Addr:    ":"+string(port),
		Handler: buffer,
	}
	err := s.ListenAndServe()
	return s, err
}

func main() {
	RProxy(9211)
}
