// +build integration

package service_test

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/contract"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadAll method behaviour", func() {
	var (
		mockCtrl                 *gomock.Controller
		tenantDataService        *service.TenantDataService
		mockUUIDGeneratorService *MockUUIDGeneratorService
		tenantID                 system.UUID
		clusterConfig            *gocql.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = NewMockUUIDGeneratorService(mockCtrl)

		tenantDataService = &service.TenantDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: clusterConfig}

		tenantID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should return error if tenant does not exist", func() {
		_, err := tenantDataService.Read(tenantID)

		Expect(err).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", tenantID.String())))
	})

	It("should return the existing tenant", func() {
		mockUUIDGeneratorService.
			EXPECT().
			GenerateRandomUUID().
			Return(tenantID, nil)

		randomValue, _ := system.RandomUUID()
		expectedTenant := contract.Tenant{SecretKey: randomValue.String()}
		returnedTenantID, err := tenantDataService.Create(expectedTenant)

		Expect(err).To(BeNil())

		returnedTenant, err := tenantDataService.Read(returnedTenantID)

		Expect(err).To(BeNil())
		Expect(expectedTenant).To(Equal(returnedTenant))
	})
})

func TestReadAllBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadAll method behaviour")
}
