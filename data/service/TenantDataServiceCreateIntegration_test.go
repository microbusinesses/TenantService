// +build integration

package service_test

import (
	"errors"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/contract"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create method behaviour", func() {
	var (
		mockCtrl                 *gomock.Controller
		tenantDataService        *service.TenantDataService
		mockUUIDGeneratorService *MockUUIDGeneratorService
		tenantID                 system.UUID
		validTenant              contract.Tenant
		clusterConfig            *gocql.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = NewMockUUIDGeneratorService(mockCtrl)

		tenantDataService = &service.TenantDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: clusterConfig}

		tenantID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		validTenant = contract.Tenant{SecretKey: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when UUID generator service succeeds to create the new UUID", func() {
		It("should return the new UUID as tenant uniuqe identifier and no error", func() {
			expectedTenantID, _ := system.RandomUUID()
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(expectedTenantID, nil)

			newTenantID, err := tenantDataService.Create(validTenant)

			Expect(expectedTenantID).To(Equal(newTenantID))
			Expect(err).To(BeNil())
		})
	})

	Context("when UUID generator service fails to create the new UUID", func() {
		It("should return tenant unique identifier as empty UUID and the returned error by tenant data service", func() {
			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(system.EmptyUUID, expectedError)

			newTenantID, err := tenantDataService.Create(validTenant)

			Expect(newTenantID).To(Equal(system.EmptyUUID))
			Expect(err).To(Equal(expectedError))
		})
	})

	Context("when creating new tenant", func() {
		It("should insert the record into tenant table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(tenantID, nil)

			returnedTenantID, err := tenantDataService.Create(validTenant)

			Expect(tenantID).To(Equal(returnedTenantID))
			Expect(err).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			Expect(err).To(BeNil())

			defer session.Close()

			iter := session.Query(
				"SELECT secret_key"+
					" FROM tenant"+
					" WHERE"+
					" tenant_id = ?",
				tenantID.String()).Iter()

			defer iter.Close()

			var secretKey string

			for iter.Scan(&secretKey) {
				Expect(validTenant.SecretKey).To(Equal(secretKey))
			}

		})
	})
})

func TestCreateBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method behaviour")
}
