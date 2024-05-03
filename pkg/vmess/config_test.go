package vmess

import (
	"context"
	_ "github.com/kelleygo/sing-box/transport/dhcp"
	"net/netip"
	"net/url"
	"testing"
	"time"
)

var (
	serviceAddr        = "3.77.232.6"
	servicePort uint64 = 443
	uuid               = "47b35bd5-7cfc-4687-a9ab-56e75e5881b8"
	logPath            = ".letsCore.log"
)

func TestParseVmessSpeed(t *testing.T) {
	logPath = CurrentPath() + "/" + logPath
	out, err := ParseVmessSpeed(logPath, serviceAddr, servicePort, uuid)
	t.Log(err)
	t.Log(string(out))
}

func TestParseVmessGlobal(t *testing.T) {
	out, err := ParseVmessGlobal(logPath, serviceAddr, servicePort, uuid)
	t.Log(err)
	t.Log(string(out))
}

func TestSpeedVmess(t *testing.T) {
	var runMode int = 1
	service := NewVmessService(runMode, serviceAddr, servicePort, uuid)
	err := service.Start(context.TODO())
	t.Log(err)
}

func TestSpeedV2Vmess(t *testing.T) {
	var runMode int = 1
	service := NewVmessService(runMode, serviceAddr, servicePort, uuid)
	err := service.Create(context.TODO())
	t.Log(err)

	time.Sleep(time.Second * 30)
}

func TestGlobalVmess(t *testing.T) {
	var runMode int = 0
	service := NewVmessService(runMode, serviceAddr, servicePort, uuid)
	err := service.Start(context.TODO())
	t.Log(err)
}

func TestDnsRoute(t *testing.T) {
	dnsAddress := "127.0.0.1"
	serverURL, _ := url.Parse(dnsAddress)
	var serverAddress string
	if serverURL != nil {
		serverAddress = serverURL.Hostname()
	}
	t.Log(dnsAddress)
	t.Log(serverURL)
	t.Log(serverAddress)
	_, notIpAddress := netip.ParseAddr(serverAddress)
	t.Log(notIpAddress)
}
