# Stackstate rancher integration

The purpose of these tests is to ensure that the StackState integration in Rancher functions correctly. Specifically, when the StackState agent chart is installed, Rancher should be able to successfully communicate with the StackState server.

## Pre-requisites

- Ensure you have an existing cluster that the user can access. If no downstream cluster is available in Rancher, it is required that the user creates one before running this test.
- A StackState API key in the StackState UI is essential to proceed with the tests. Navigate to the StackState UI and create an instance using the cluster name from the previous step which in turn returns an api key.
- For the tests TestDynamicUpgradeStackstateAgentChart, make sure to provide the version stackstate needs to be upgraded to. Otherwise the test will be skipped.
- For StackState Server tests the following pre-requisites are needed
  - The minimum instance size in AWS is t3a.xlarge and the cluster size needs to be 1 Control Plane, ETCD Node and 3 Worker Nodes
  - The admin password and clusterAPIKey is created when generating the template yaml for installing the Stackstate server
  - The license key has to be obtained from the Observability team
- Documentation for generating the chart value yaml files is located at <https://docs.stackstate.com/self-hosted-setup/install-stackstate/kubernetes_openshift/kubernetes_install>

## Test Setup

Your GO suite should be set to `-run ^Test<TestSuite>$`. For example to run the stackstate_test.go, set the GO suite to `-run ^TestStackStateTestSuite$` You can find specific tests by checking the test file you plan to run.

In your config file, set the following:

```yaml
rancher: 
  host: "rancher_server_address"
  adminToken: "rancher_admin_token"
  insecure: True #optional
  cleanup: True #optional
  clusterName: "downstream_cluster_name"

stackstateConfigs: 
  url: "stackstate_ui_server_address"
  serviceToken: "stackstate_user_service_token"
  clusterApiKey: "stackstate_cluster_apikey"
  upgradeVersion: "Stackstate_agent_version" #optional
  adminPassword: "Encrypted Password generated from Helm Chart"
  license: "Stackstate License Key"
```
