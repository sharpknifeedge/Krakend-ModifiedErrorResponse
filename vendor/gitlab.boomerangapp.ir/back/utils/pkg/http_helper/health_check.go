package http_helper

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/kataras/iris/v12"
	"gitlab.boomerangapp.ir/back/utils/pkg/env"
)

//IrisHelthCheck is a health check for the iris router
func IrisHelthCheck(m iris.Party, version string) {
	m.Get("/health", func(ctx iris.Context) {
		ctx.ResponseWriter().Header().Add("Content-Type", "application/json")
		ctx.JSON(GetHealthCheckObject(version))
	})
}

//MuxHelthCheck is a health check for the gorilla mux router
func MuxHelthCheck(m *mux.Router, version string) {
	m.Methods(http.MethodGet).
		Path("/health").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(GetHealthCheckObject(version))
		})
}

//serverStartTime is the time when the server started
var serverStartTime = time.Now()

//HealthCheck is the object that will be sent to the client through the http
type HealthCheck struct {
	Status        bool   `json:"status"`
	ServerName    string `json:"server_name"`
	ServerVersion string `json:"server_version"`
	ServerTime    struct {
		Time     time.Time `json:"time"`
		UnixTime int64     `json:"unix_time"`
	} `json:"server_time"`
	MemProfileRate    int     `json:"memory_allocation"`
	NumCPU            int     `json:"num_cpu"`
	NumGoroutine      int     `json:"num_goroutine"`
	GoVersion         string  `json:"go_version"`
	GoOS              string  `json:"go_os"`
	GoArch            string  `json:"go_arch"`
	ProccessUptimeSec float64 `json:"process_uptime_sec"`
}

//GetHealthCheckObject returns the health check object for http
func GetHealthCheckObject(ver string) HealthCheck {
	return HealthCheck{
		Status:        true,
		ServerName:    env.Str("SERVER_NAME", "Not defined"),
		ServerVersion: ver,
		ServerTime: struct {
			Time     time.Time `json:"time"`
			UnixTime int64     `json:"unix_time"`
		}{
			Time:     time.Now(),
			UnixTime: time.Now().Unix(),
		},
		MemProfileRate:    runtime.MemProfileRate,
		NumCPU:            runtime.NumCPU(),
		NumGoroutine:      runtime.NumGoroutine(),
		GoVersion:         runtime.Version(),
		GoOS:              runtime.GOOS,
		GoArch:            runtime.GOARCH,
		ProccessUptimeSec: time.Now().Sub(serverStartTime).Seconds(),
	}
}
