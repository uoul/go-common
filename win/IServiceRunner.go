package win

type IServiceRunner interface {
	Run(command ServiceRunnerCommand) error
}
