package win

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var ErrUnknownCommand = errors.New("unknown command")
var ErrNotRunningAsService = errors.New("not running as windows service")

type ServiceRunnerCommand uint

const (
	SVC_EXEC      = ServiceRunnerCommand(0)
	SVC_START     = ServiceRunnerCommand(1)
	SVC_STOP      = ServiceRunnerCommand(2)
	SVC_INSTALL   = ServiceRunnerCommand(3)
	SVC_UNINSTALL = ServiceRunnerCommand(4)
	SVC_PAUSE     = ServiceRunnerCommand(5)
	SVC_CONTINUE  = ServiceRunnerCommand(6)
	SVC_DEBUG     = ServiceRunnerCommand(7)
)

type ServiceRunner struct {
	svc  IService
	elog debug.Log
}

// Run implements IServiceRunner.
func (sr *ServiceRunner) Run(command ServiceRunnerCommand) error {
	isRunningAsService, err := svc.IsWindowsService()
	if err != nil {
		return err
	}
	if isRunningAsService {
		sr.runService(false)
	} else {
		switch command {
		case SVC_START:
			return sr.startService()
		case SVC_STOP:
			return sr.controlService(svc.Stop, svc.Stopped)
		case SVC_INSTALL:
			return sr.installService()
		case SVC_UNINSTALL:
			return sr.removeService()
		case SVC_PAUSE:
			return sr.controlService(svc.Pause, svc.Paused)
		case SVC_CONTINUE:
			return sr.controlService(svc.Continue, svc.Running)
		case SVC_DEBUG:
			sr.runService(true)
		case SVC_EXEC:
			return ErrNotRunningAsService
		default:
			return ErrUnknownCommand
		}
	}
	return nil
}

func NewServiceRunner(svc IService) IServiceRunner {
	return &ServiceRunner{
		svc:  svc,
		elog: nil,
	}
}

func (sr *ServiceRunner) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	sr.svc.OnStart(sr.elog, args)
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		changeRequest := <-r
		switch changeRequest.Cmd {
		case svc.Stop:
			changes <- svc.Status{State: svc.StopPending}
			sr.svc.OnStop()
			break loop
		case svc.Pause:
			changes <- svc.Status{State: svc.PausePending}
			sr.svc.OnPause()
			changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
		case svc.Continue:
			changes <- svc.Status{State: svc.ContinuePending}
			sr.svc.OnContinue()
			changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
		case svc.Interrogate:
			sr.svc.OnInterrogate()
			changes <- changeRequest.CurrentStatus
		case svc.Shutdown:
			changes <- svc.Status{State: svc.StopPending}
			sr.svc.OnShutdown()
			break loop
		default:
			sr.elog.Error(1, fmt.Sprintf("unexpected control request #%d", changeRequest))
		}
	}
	changes <- svc.Status{State: svc.Stopped}
	return
}

func (sr *ServiceRunner) exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			return "", fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

func (sr *ServiceRunner) installService() error {
	exepath, err := sr.exePath()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(sr.svc.GetServiceName())
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", sr.svc.GetServiceName())
	}
	s, err = m.CreateService(sr.svc.GetServiceName(), exepath, mgr.Config{DisplayName: sr.svc.GetServiceDescription()}, "is", "auto-started")
	if err != nil {
		return err
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(sr.svc.GetServiceName(), eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}

func (sr *ServiceRunner) removeService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(sr.svc.GetServiceName())
	if err != nil {
		return fmt.Errorf("service %s is not installed", sr.svc.GetServiceName())
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(sr.svc.GetServiceName())
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}

func (sr *ServiceRunner) startService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(sr.svc.GetServiceName())
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

func (sr *ServiceRunner) controlService(c svc.Cmd, to svc.State) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(sr.svc.GetServiceName())
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}

func (sr *ServiceRunner) runService(isDebug bool) {
	var err error
	if isDebug {
		sr.elog = debug.New(sr.svc.GetServiceName())
	} else {
		sr.elog, err = eventlog.Open(sr.svc.GetServiceName())
		if err != nil {
			return
		}
	}
	defer sr.elog.Close()

	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(sr.svc.GetServiceName(), sr)
	if err != nil {
		sr.elog.Error(1, fmt.Sprintf("%s service failed: %v", sr.svc.GetServiceName(), err))
		return
	}
}
