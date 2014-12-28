package service

const (
	UdpTransport = "udp"
	TcpTransport = "tcp"
)

const (
	HttpProtocol  = "http"
	HttpsProtocol = "https"
)

type Location struct {
	Address string
	Port    int

	Protocol  string
	Transport string

	Priority int
	Weight   int
}

type ServiceManager interface {
	Find(serviceName string) ([]Location, error)
	Register(serviceName string, location Location) error
}

var serviceManager ServiceManager

func InitFinder(manager ServiceManager) {
	serviceManager = manager
}

func Find(serviceName string) ([]Location, error) {
	if serviceManager != nil {
		return serviceManager.Find(serviceName)
	}
	panic("No service locator configured")
	return []Location{}, nil
}

func FindByDef(serviceDefinition ServiceOperation) ([]Location, error) {
	return Find(serviceDefinition.ServiceName)
}
