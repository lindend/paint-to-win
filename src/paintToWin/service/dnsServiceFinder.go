package service

type DnsServiceFinder struct {
}

func (d DnsServiceFinder) Find(serviceName string) (Location, error) {
	return Location{}, nil
}
