module github.com/arangodb-managed/terraform-provider-oasis

require (
	github.com/arangodb-managed/apis v0.71.2
	github.com/arangodb-managed/log-helper v0.2.0
	github.com/gogo/protobuf v1.3.0
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/rs/zerolog v1.19.0
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.33.2
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

go 1.16

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible

replace github.com/arangodb/kube-arangodb => github.com/arangodb/kube-arangodb v0.0.0-20210528082542-41972ad9b013

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a

replace github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.37.0

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3

replace github.com/hashicorp/vault/api => github.com/hashicorp/vault/api v1.0.5-0.20201006192546-7a875e245472

replace github.com/hashicorp/vault/sdk => github.com/hashicorp/vault/sdk v0.1.14-0.20201110183819-7ce0bd969199

replace github.com/hashicorp/vault => github.com/hashicorp/vault v1.6.0

replace github.com/kamilsk/retry => github.com/kamilsk/retry/v3 v3.4.2

replace github.com/nats-io/go-nats-streaming => github.com/nats-io/go-nats-streaming v0.4.4

replace github.com/nats-io/go-nats => github.com/nats-io/go-nats v1.7.2

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.0-20190108154635-47c0da630f72

replace github.com/ugorji/go => github.com/ugorji/go v0.0.0-20181204163529-d75b2dcb6bc8

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20210503173754-0981d6026fa6

replace google.golang.org/api => google.golang.org/api v0.36.0

replace google.golang.org/grpc => google.golang.org/grpc v1.36.0

replace k8s.io/api => k8s.io/api v0.19.8

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.8

replace k8s.io/apimachinery => k8s.io/apimachinery v0.19.8

replace k8s.io/apiserver => k8s.io/apiserver v0.19.8

replace k8s.io/client-go => k8s.io/client-go v0.19.8

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.8

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.8

replace k8s.io/code-generator => k8s.io/code-generator v0.19.8

replace k8s.io/component-base => k8s.io/component-base v0.19.8

replace k8s.io/kubernetes => k8s.io/kubernetes v1.19.8

replace k8s.io/metrics => k8s.io/metrics v0.19.8

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.8.3

replace sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.0

replace github.com/hashicorp/terraform-plugin-sdk => github.com/hashicorp/terraform-plugin-sdk v1.4.1

replace github.com/cilium/cilium => github.com/cilium/cilium v1.9.5

replace github.com/optiopay/kafka => github.com/optiopay/kafka v2.0.4+incompatible

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.8

replace k8s.io/cri-api => k8s.io/cri-api v0.19.8

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.8

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.8

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.8

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.8

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.8

replace k8s.io/kubelet => k8s.io/kubelet v0.19.8

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.8

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.8

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20201211151036-40ec1c210f7a

replace k8s.io/kubectl => k8s.io/kubectl v0.19.8

replace github.com/nats-io/nats.go => github.com/nats-io/nats.go v1.10.0

replace github.com/nats-io/stan.go => github.com/nats-io/stan.go v0.8.3

replace github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring => github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.47.1

replace github.com/prometheus-operator/prometheus-operator/pkg/client => github.com/prometheus-operator/prometheus-operator/pkg/client v0.47.1

replace github.com/prometheus-operator/prometheus-operator => github.com/prometheus-operator/prometheus-operator v0.47.1

replace go.uber.org/multierr => go.uber.org/multierr v1.6.1-0.20201027220001-0eb6eb5383b9
