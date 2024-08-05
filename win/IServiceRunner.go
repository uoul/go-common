//go:build windows
// +build windows

package win

type IServiceRunner interface {
	Run(command ServiceRunnerCommand) error
}
