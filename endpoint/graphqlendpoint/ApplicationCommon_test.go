package graphqlendpoint_test

import (
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/business/domain"
)

func createApplicationInfo() domain.Application {
	randomValue, _ := system.RandomUUID()
	return domain.Application{Name: randomValue.String()}
}
