//go:build (infra.rke2k3s || validation || recurring) && !infra.any && !infra.aks && !infra.eks && !infra.gke && !infra.rke1 && !stress && !sanity && !extended

package rke2k3s

import (
	"os"
	"testing"

	"github.com/rancher/shepherd/clients/rancher"
	extClusters "github.com/rancher/shepherd/extensions/clusters"
	"github.com/rancher/shepherd/pkg/config"
	"github.com/rancher/shepherd/pkg/config/operations"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/rancher/tests/actions/clusters"
	"github.com/rancher/tests/actions/config/defaults"
	"github.com/rancher/tests/actions/provisioninginput"
	resources "github.com/rancher/tests/validation/provisioning/resources/provisioncluster"
	standard "github.com/rancher/tests/validation/provisioning/resources/standarduser"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DeleteInitMachineTestSuite struct {
	suite.Suite
	client        *rancher.Client
	session       *session.Session
	cattleConfig  map[string]any
	rke2ClusterID string
	k3sClusterID  string
}

func (d *DeleteInitMachineTestSuite) TearDownSuite() {
	d.session.Cleanup()
}

func (d *DeleteInitMachineTestSuite) SetupSuite() {
	testSession := session.NewSession()
	d.session = testSession

	client, err := rancher.NewClient("", testSession)
	require.NoError(d.T(), err)

	d.client = client

	standardUserClient, err := standard.CreateStandardUser(d.client)
	require.NoError(d.T(), err)

	d.cattleConfig = config.LoadConfigFromFile(os.Getenv(config.ConfigEnvironmentKey))

	rke2ClusterConfig := new(clusters.ClusterConfig)
	operations.LoadObjectFromMap(defaults.ClusterConfigKey, d.cattleConfig, rke2ClusterConfig)

	k3sClusterConfig := new(clusters.ClusterConfig)
	operations.LoadObjectFromMap(defaults.ClusterConfigKey, d.cattleConfig, k3sClusterConfig)

	nodeRolesStandard := []provisioninginput.MachinePools{
		provisioninginput.EtcdMachinePool,
		provisioninginput.ControlPlaneMachinePool,
		provisioninginput.WorkerMachinePool,
	}

	nodeRolesStandard[0].MachinePoolConfig.Quantity = 3
	nodeRolesStandard[1].MachinePoolConfig.Quantity = 2
	nodeRolesStandard[2].MachinePoolConfig.Quantity = 3

	rke2ClusterConfig.MachinePools = nodeRolesStandard
	k3sClusterConfig.MachinePools = nodeRolesStandard

	d.rke2ClusterID, err = resources.ProvisionRKE2K3SCluster(d.T(), standardUserClient, extClusters.RKE2ClusterType.String(), rke2ClusterConfig, true, false)
	require.NoError(d.T(), err)

	d.k3sClusterID, err = resources.ProvisionRKE2K3SCluster(d.T(), standardUserClient, extClusters.K3SClusterType.String(), k3sClusterConfig, true, false)
	require.NoError(d.T(), err)
}

func (d *DeleteInitMachineTestSuite) TestDeleteInitMachine() {
	tests := []struct {
		name      string
		clusterID string
	}{
		{"RKE2_Delete_Init_Machine", d.rke2ClusterID},
		{"K3S_Delete_Init_Machine", d.k3sClusterID},
	}

	for _, tt := range tests {
		d.Run(tt.name, func() {
			err := deleteInitMachine(d.client, tt.clusterID)
			require.NoError(d.T(), err)
		})
	}
}

func TestDeleteInitMachineTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteInitMachineTestSuite))
}
