package main

import "github.com/eventscompass/service-framework/service"

type ServiceName struct {
	service.BaseService
}

func main() {
	service.Start(&ServiceName{})
}
