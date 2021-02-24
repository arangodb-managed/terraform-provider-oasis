module github.com/arangodb-managed/terraform-provider-oasis

require (
	github.com/arangodb-managed/apis v0.59.0
	github.com/arangodb-managed/log-helper v0.1.4
	github.com/gogo/protobuf v1.3.0
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.4.1
	github.com/rs/zerolog v1.17.2
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.29.1
)

go 1.16

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.0.1+incompatible

replace github.com/arangodb/kube-arangodb => github.com/arangodb/kube-arangodb v0.0.0-20201104155747-ca02a3a4c884

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a

replace github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.37.0

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

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

replace k8s.io/api => k8s.io/api v0.16.15

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.15

replace k8s.io/apimachinery => k8s.io/apimachinery v0.16.15

replace k8s.io/apiserver => k8s.io/apiserver v0.16.15

replace k8s.io/client-go => k8s.io/client-go v0.16.15

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.15

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.15

replace k8s.io/code-generator => k8s.io/code-generator v0.16.15

replace k8s.io/component-base => k8s.io/component-base v0.16.15

replace k8s.io/kubernetes => k8s.io/kubernetes v1.16.15

replace k8s.io/metrics => k8s.io/metrics v0.16.15

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.4.0

replace sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.0

replace github.com/hashicorp/terraform-plugin-sdk => github.com/hashicorp/terraform-plugin-sdk v1.4.1

replace github.com/cilium/cilium => github.com/cilium/cilium v1.7.2

replace github.com/optiopay/kafka => github.com/optiopay/kafka v2.0.4+incompatible

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.15

replace k8s.io/cri-api => k8s.io/cri-api v0.16.15

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.15

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.15

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.15

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.15

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.15

replace k8s.io/kubelet => k8s.io/kubelet v0.16.15

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.15

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.15

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20191206224255-0243a4be9c8f

replace k8s.io/kubectl => k8s.io/kubectl v0.16.15

replace github.com/nats-io/nats.go => github.com/nats-io/nats.go v1.8.1

replace github.com/nats-io/stan.go => github.com/nats-io/stan.go v0.4.5
