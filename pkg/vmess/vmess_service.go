package vmess

import (
	"context"
	box "github.com/kelleygo/sing-box"
	"github.com/kelleygo/sing-box/constant"
	"github.com/kelleygo/sing-box/log"
	"github.com/kelleygo/sing-box/option"
	"github.com/kelleygo/sing-box/pkg"
	_ "github.com/kelleygo/sing-box/transport/dhcp"
	E "github.com/sagernet/sing/common/exceptions"
	"os"
	"os/signal"
	runtimeDebug "runtime/debug"
	"syscall"
	"time"
)

type vmessService struct {
	vmessGlobalProxy bool
	serviceAddr      string
	servicePort      uint64
	uuid             string
	instance         *box.Box
	cancel           context.CancelFunc
}

var defaultLogPtah = ".letsCore.log"

func ParseVmessSpeed(logPath string, serviceAddr string, servicePort uint64, uuid string) ([]byte, error) {
	if logPath == "" {
		logPath = CurrentPath() + "/" + defaultLogPtah
	}
	return TemplateContent(vmessSpeedTemplate, map[string]interface{}{
		"logPath":     logPath,
		"serviceAddr": serviceAddr,
		"servicePort": servicePort,
		"userUuid":    uuid,
	})
}

func ParseVmessGlobal(logPath string, serviceAddr string, servicePort uint64, uuid string) ([]byte, error) {
	if logPath == "" {
		logPath = CurrentPath() + "/" + defaultLogPtah
	}
	return TemplateContent(vmessGlobalTemplate, map[string]interface{}{
		"logPath":     logPath,
		"serviceAddr": serviceAddr,
		"servicePort": servicePort,
		"userUuid":    uuid,
	})
}

func (this *vmessService) create() (*box.Box, context.CancelFunc, error) {
	var err error
	var options option.Options
	options, err = this.readConfigAndMerge()
	if err != nil {
		return nil, nil, E.Cause(err, "readConfigAndMerge")
	}

	ctx, cancel := context.WithCancel(context.TODO())
	var instance *box.Box
	for retry := 0; retry < 3; retry++ {
		instance, err = box.New(box.Options{
			Context: ctx,
			Options: options,
		})
		if err != nil {
			log.Warn(err)
			continue
		}
	}

	if err != nil {
		cancel()
		return nil, nil, E.Cause(err, "create service")
	}
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer func() {
		signal.Stop(osSignals)
		close(osSignals)
	}()
	startCtx, finishStart := context.WithCancel(context.Background())
	go func() {
		_, loaded := <-osSignals
		if loaded {
			cancel()
			closeMonitor(startCtx)
		}
	}()
	err = instance.Start()
	finishStart()
	if err != nil {
		cancel()
		return nil, nil, E.Cause(err, "start service")
	}
	return instance, cancel, nil
}

func closeMonitor(ctx context.Context) {
	time.Sleep(constant.FatalStopTimeout)
	select {
	case <-ctx.Done():
		return
	default:
	}
	log.Fatal("sing-box did not close!")
}

func (this *vmessService) readConfigAndMerge() (option.Options, error) {
	content, err := ParseVmessSpeed("", this.serviceAddr, this.servicePort, this.uuid)
	if this.vmessGlobalProxy {
		content, err = ParseVmessGlobal("", this.serviceAddr, this.servicePort, this.uuid)
	}
	var mergedOptions option.Options
	err = mergedOptions.UnmarshalJSON(content)
	if err != nil {
		return option.Options{}, E.Cause(err, "unmarshal merged config")
	}
	return mergedOptions, nil
}

func (this *vmessService) Start(ctx context.Context) error {
	var err error
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(osSignals)
	for {
		this.instance, this.cancel, err = this.create()
		if err != nil {
			return err
		}
		runtimeDebug.FreeOSMemory()
		select {
		case <-ctx.Done():
			this.Close()
			return nil
		case <-osSignals:
			this.Close()
			return nil
		}
	}
}

func (this *vmessService) Close() error {
	closeCtx, closed := context.WithCancel(context.Background())
	go closeMonitor(closeCtx)
	closed()
	this.instance.Close()
	return nil
}

// NewVmessService
// 运行模式，默认0 极速，1 全局模式
func NewVmessService(runMode int, serviceAddr string, servicePort uint64, uuid string) pkg.Tun2VmessService {
	vmessGlobalProxy := false
	if runMode == 1 {
		vmessGlobalProxy = true
	}
	return &vmessService{
		serviceAddr:      serviceAddr,
		servicePort:      servicePort,
		uuid:             uuid,
		vmessGlobalProxy: vmessGlobalProxy,
	}
}
