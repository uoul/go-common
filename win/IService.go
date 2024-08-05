//go:build windows
// +build windows

package win

import "golang.org/x/sys/windows/svc/debug"

type IService interface {
	GetServiceName() string
	GetServiceDescription() string
	OnStart(elog debug.Log, args []string)
	OnStop()
	OnPause()
	OnContinue()
	OnInterrogate()
	OnShutdown()
}
