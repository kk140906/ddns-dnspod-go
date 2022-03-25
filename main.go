package main

import (
	"github.com/kk140906/ddns-dnspod-go/tencentcloud"
	"github.com/kk140906/ddns-dnspod-go/utils"
	"github.com/kk140906/ddns-dnspod-go/zap_wrapper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer func() { _ = recover() }()
	config := utils.DDNSConfig{}
	_ = utils.ReadConfig(&config)
	go schedule(config.FlushAfterMinutes)
	<-quit
}

func schedule(minutes int) {
	duration := time.Duration(int64(time.Minute) * int64(minutes))
	updateDnspod() // 启动时运行一次
	for {
		select {
		case <-time.After(duration):
			go updateDnspod()
			// default:
			// time.Sleep(100000)
		}
	}
}

func updateDnspod() {
	zap_wrapper.DefaultLogger.Info("prepare update ip")
	ip := utils.GetIp()
	if ip == "" {
		return
	}
	client := tencentcloud.NewClient()
	for k, v := range client.Domain {
		res, err := client.GetDomainRecordList(k)
		if err != nil {
			continue
		}
		for _, subdomain := range v {
			if _, ok := res[subdomain]; !ok {
				zap_wrapper.DefaultLogger.Error("subdomain not exist",
					zap.String("subdomain", subdomain),
					zap.String("domain", k))
				continue
			}
			if res[subdomain].Value == ip {
				zap_wrapper.DefaultLogger.Warn("subdomain ip not changed, skip update",
					zap.String("subdomain", subdomain),
					zap.String("domain", k),
					zap.String("ip", ip))
				continue
			}
			_, err := client.UpdateDomainRecord(k, subdomain, ip, res[subdomain])
			if err != nil {
				zap_wrapper.DefaultLogger.Warn("update subdomain ip failed",
					zap.String("subdomain", subdomain),
					zap.String("domain", k),
					zap.String("error", err.Error()))
				continue
			}
			zap_wrapper.DefaultLogger.Info("update ip success",
				zap.String("subdomain", subdomain),
				zap.String("domain", k),
				zap.String("previous ip", res[subdomain].Value),
				zap.String("new ip", ip))
		}
	}
}
