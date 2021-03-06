package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type Info struct {
	Hostname      string `json:"hostname"`
	KernelVersion string `json:"kernel_version"`
	KernelArch    string `json:"kernel_arch"`
	Uptime        string `json:"uptime"`
	MemoryFree    uint64 `json:"memory_free"`
	MemoryTotal   uint64 `json:"memory_total"`
	MemoryUsed    string `json:"memory_used"`
}

type Server struct {
	listenAddr    string
	appVersion    string
	traceProvider *tracesdk.TracerProvider
}

func NewServer(addr, appVersion, traceEndpoint string) (*Server, error) {
	// enable tracing if specified
	var traceProvider *tracesdk.TracerProvider
	if traceEndpoint != "" {
		tp, err := newTraceProvider(traceEndpoint, "sysinfo", appVersion)
		if err != nil {
			return nil, err
		}
		traceProvider = tp
	}

	return &Server{
		listenAddr:    addr,
		appVersion:    appVersion,
		traceProvider: traceProvider,
	}, nil
}

func (s *Server) Run() error {
	r := mux.NewRouter()
	r.Handle("/", otelhttp.NewHandler(http.HandlerFunc(s.infoHandler), "info"))
	r.Handle("/version", otelhttp.NewHandler(http.HandlerFunc(s.versionHandler), "version"))

	srv := &http.Server{
		Handler:      r,
		Addr:         s.listenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}

	logrus.Infof("starting server on %s", s.listenAddr)

	return srv.ListenAndServe()
}

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info, err := host.Info()
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	memStat, err := mem.VirtualMemory()
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	up := time.Duration(info.Uptime) * time.Second
	sysinfo := Info{
		Hostname:      hostname,
		KernelVersion: info.KernelVersion,
		KernelArch:    info.KernelArch,
		Uptime:        humanize.Time(time.Now().Add(-up)),
		MemoryFree:    memStat.Free,
		MemoryTotal:   memStat.Total,
		MemoryUsed:    fmt.Sprintf("%.2f", memStat.UsedPercent),
	}

	if err := json.NewEncoder(w).Encode(sysinfo); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(s.appVersion + "\n"))
}
