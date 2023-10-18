package prometheus

import (
	"das-account-indexer/config"
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"net"
	"sync"
	"time"
)

var (
	log          = logger.NewLogger("prometheus", logger.LevelDebug)
	PromRegister = prometheus.NewRegistry()
)

var Tools *Prometheus

type Prometheus struct {
	pusher  *push.Pusher
	Metrics Metric
}

type Metric struct {
	l         sync.Mutex
	api       *prometheus.SummaryVec
	errNotify *prometheus.CounterVec
}

func (m *Metric) Api() *prometheus.SummaryVec {
	if m.api == nil {
		m.l.Lock()
		defer m.l.Unlock()
		m.api = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: "api",
		}, []string{"method", "http_status", "err_no", "err_msg"})
		PromRegister.MustRegister(m.api)
	}
	return m.api
}

func (m *Metric) ErrNotify() *prometheus.CounterVec {
	if m.errNotify == nil {
		m.l.Lock()
		defer m.l.Unlock()
		m.errNotify = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "notify",
		}, []string{"title", "text"})
		PromRegister.MustRegister(m.errNotify)
	}
	return m.errNotify
}

func Init() {
	Tools = &Prometheus{}
}

func (t *Prometheus) Run() {
	if config.Cfg.Server.PrometheusPushGateway != "" && config.Cfg.Server.Name != "" {
		t.pusher = push.New(config.Cfg.Server.PrometheusPushGateway, config.Cfg.Server.Name)
		t.pusher.Gatherer(PromRegister)
		t.pusher.Grouping("env", fmt.Sprint(config.Cfg.Server.Net))
		t.pusher.Grouping("instance", GetLocalIp("eth0"))

		go func() {
			ticker := time.NewTicker(time.Second * 5)
			defer ticker.Stop()

			for range ticker.C {
				_ = t.pusher.Push()
			}
		}()
	}
}

func GetLocalIp(interfaceName string) string {
	ief, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Error("GetLocalIp: ", err)
		return ""
	}
	addrs, err := ief.Addrs()
	if err != nil {
		log.Error("GetLocalIp: ", err)
		return ""
	}

	var ipv4Addr net.IP
	for _, addr := range addrs {
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		log.Errorf("GetLocalIp interface %s don't have an ipv4 address", interfaceName)
		return ""
	}
	return ipv4Addr.String()
}
