package resource

type IResourceManager interface {
	Register(r IResource)
	Unregister(r IResource)
	Wait()
}
