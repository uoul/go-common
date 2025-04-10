package health

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type HealthMonitor struct {
	readynessChecks map[string]func() error
}

var instance *HealthMonitor
var lock = &sync.Mutex{}

// GetHealthMonitor() returns an instance of the health monitor
func GetHealthMonitor() *HealthMonitor {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &HealthMonitor{
				readynessChecks: map[string]func() error{},
			}
		}
	}
	return instance
}

// RegisterReadynessCheck(name string, check func() error) registers checks, that will be executed
// on endpointcall for /startup or /ready. If any function returns an error, the probe will fail,
// otherwise it will pass with http code 200
//
// Params:
// * name - name of given check (will be part of error message if check fails)
// * check - function, that executes the given check (return nil if check pass)
func (m *HealthMonitor) RegisterReadynessCheck(name string, check func() error) {
	m.readynessChecks[name] = check
}

// UnregisterReadynessCheck(name string) removes a predefined check
func (m *HealthMonitor) UnregisterReadynessCheck(name string) {
	delete(m.readynessChecks, name)
}

// RegisterEndpointsDefault(...) registers three endpoints on a given http mutex behind a given prefix
// Endpoints:
// - <prefix>/startup (Startup probe)
// - <prefix>/alive (Liveness probe)
// - <prefix>/ready (Readyness probe)
func (m *HealthMonitor) RegisterEndpointsDefault(mux *http.ServeMux, urlPrefix string) {
	mux.HandleFunc(fmt.Sprintf("%s/startup", urlPrefix), m.readynessProbeDefault)
	mux.HandleFunc(fmt.Sprintf("%s/alive", urlPrefix), func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc(fmt.Sprintf("%s/ready", urlPrefix), m.readynessProbeDefault)
}

func (m *HealthMonitor) readynessProbeDefault(w http.ResponseWriter, r *http.Request) {
	for name, check := range m.readynessChecks {
		if err := check(); err != nil {
			http.Error(w, fmt.Sprintf("%s: %v", name, err), http.StatusServiceUnavailable)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// RegisterEndpoints(...) registers three endpoints on a given gin router
// Endpoints:
// - <prefix>/startup (Startup probe)
// - <prefix>/alive (Liveness probe)
// - <prefix>/ready (Readyness probe)
func (m *HealthMonitor) RegisterEndpointsGin(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/alive", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
	routerGroup.GET("/startup", m.readynessProbeGin)
	routerGroup.GET("/ready", m.readynessProbeGin)
}

func (m *HealthMonitor) readynessProbeGin(ctx *gin.Context) {
	for name, check := range m.readynessChecks {
		if err := check(); err != nil {
			ctx.String(http.StatusServiceUnavailable, "%s: %v", name, err)
			return
		}
	}
	ctx.Status(http.StatusOK)
}

// DoReadynessChecks() executes all registerd readyness checks and returns
// a slice of all occured errors
func (m *HealthMonitor) DoReadynessChecks() []error {
	r := []error{}
	for name, check := range m.readynessChecks {
		if err := check(); err != nil {
			r = append(r, fmt.Errorf("ReadynessProbe %s failed - %v", name, err))
		}
	}
	return r
}
