package domainmanager

import (
	"fmt"
	"net/http"

	"github.com/sekiseigumi/dattebayo/internal/logger"
)

type DomainManagerServer struct {
	Port   int
	logger *logger.Logger
}

func NewDomainManagerServer(port int, log *logger.Logger) *DomainManagerServer {
	return &DomainManagerServer{
		Port:   port,
		logger: log,
	}
}

func (d *DomainManagerServer) Start() error {
	http.HandleFunc("/", d.handleRequest)

	d.logger.Log("DOMAIN MANAGER", fmt.Sprintf("Starting API server on port %d", d.Port))
	go func() {
		if err := http.ListenAndServe(":"+string(d.Port), nil); err != nil {
			d.logger.Log("DOMAIN MANAGER", fmt.Sprintf("API server error: %v", err))
		}
	}()

	return nil
}

func (d *DomainManagerServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
