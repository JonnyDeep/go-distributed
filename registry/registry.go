package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceUrl       string
	RequiredService  []ServiceName
	ServiceUpdateURL string
	HeartBeatURL     string
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Add    []patchEntry
	Remove []patchEntry
}
