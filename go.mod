module github.com/arangodb-managed/terraform-provider-oasis

require (
	github.com/arangodb-managed/apis v0.40.3
	github.com/arangodb-managed/log-helper v0.1.4
	github.com/gogo/protobuf v1.3.0
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.4.1
	github.com/rs/zerolog v1.17.2
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.27.1
)

go 1.13

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.0+incompatible

replace github.com/arangodb/kube-arangodb => github.com/arangodb/kube-arangodb v0.0.0-20200325103723-b449baf82554

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a

replace github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.31.1

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3

replace github.com/hashicorp/vault/api => github.com/hashicorp/vault/api v1.0.2-0.20190424005855-e25a8a1c7480

replace github.com/hashicorp/vault/sdk => github.com/hashicorp/vault/sdk v0.1.10-0.20190506194144-8fc8af3199a1

replace github.com/hashicorp/vault => github.com/hashicorp/vault v1.1.2

replace github.com/kamilsk/retry => github.com/kamilsk/retry/v3 v3.4.2

replace github.com/nats-io/go-nats-streaming => github.com/nats-io/go-nats-streaming v0.4.4

replace github.com/nats-io/go-nats => github.com/nats-io/go-nats v1.7.2

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.0-20190108154635-47c0da630f72

replace github.com/ugorji/go => github.com/ugorji/go v0.0.0-20181204163529-d75b2dcb6bc8

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20191204072324-ce4227a45e2e

replace google.golang.org/api => google.golang.org/api v0.7.0

replace google.golang.org/grpc => google.golang.org/grpc v1.21.1

replace k8s.io/api => k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190409022649-727a075fdec8

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d

replace k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190409021813-1ec86e4da56c

replace k8s.io/client-go => k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190409023720-1bc0c81fa51d

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20190409023614-027c502bb854

replace k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190311093542-50b561225d70

replace k8s.io/component-base => k8s.io/component-base v0.0.0-20190409021516-bd2732e5c3f7

replace k8s.io/kubernetes => k8s.io/kubernetes v1.14.1

replace k8s.io/metrics => k8s.io/metrics v0.0.0-20190409022812-850dadb8b49c

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.2.0-beta.2

replace sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.0

replace github.com/hashicorp/terraform-plugin-sdk => github.com/hashicorp/terraform-plugin-sdk v1.4.1
