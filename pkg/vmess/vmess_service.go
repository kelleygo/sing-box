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

func (this *vmessService) create(ctx context.Context) (*box.Box, context.CancelFunc, error) {
	var err error
	var options option.Options
	options, err = this.readConfigAndMerge()
	if err != nil {
		return nil, nil, E.Cause(err, "readConfigAndMerge")
	}

	ctx, cancel := context.WithCancel(ctx)
	var instance *box.Box
	for retry := 0; retry < 3; retry++ {
		instance, err = box.New(box.Options{
			Context: ctx,
			Options: options,
		})
		if err != nil {
			log.Warn(err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	if err != nil {
		cancel()
		return nil, nil, E.Cause(err, "create service")
	}

	err = instance.Start()
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

func (this *vmessService) Create(ctx context.Context) error {
	var err error
	this.instance, this.cancel, err = this.create(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (this *vmessService) Start(ctx context.Context) error {
	var err error
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(osSignals)

	for {
		this.instance, this.cancel, err = this.create(ctx)
		if err != nil {
			return err
		}
		runtimeDebug.FreeOSMemory()
		for {
			osSignal := <-osSignals
			if osSignal == syscall.SIGHUP {
				err = this.check()
				if err != nil {
					log.Error(E.Cause(err, "reload service"))
					continue
				}
			}
			if osSignal != syscall.SIGHUP {
				return nil
			}
			break
		}
		select {
		case <-ctx.Done():
			this.Close()
			return nil
		}
	}
}

func (this *vmessService) Close() error {
	this.cancel()
	closeCtx, closed := context.WithCancel(context.Background())
	go closeMonitor(closeCtx)
	this.instance.Close()
	closed()
	return nil
}

func (this *vmessService) check() error {
	options, err := this.readConfigAndMerge()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	instance, err := box.New(box.Options{
		Context: ctx,
		Options: options,
	})
	if err == nil {
		instance.Close()
	}
	cancel()
	return err
}

// NewVmessService
// 运行模式，默认0 全局模式，1极速模式
func NewVmessService(runMode int, serviceAddr string, servicePort uint64, uuid string) pkg.Tun2VmessService {
	vmessGlobalProxy := false
	if runMode == 0 {
		vmessGlobalProxy = true
	}
	return &vmessService{
		serviceAddr:      serviceAddr,
		servicePort:      servicePort,
		uuid:             uuid,
		vmessGlobalProxy: vmessGlobalProxy,
	}
}
