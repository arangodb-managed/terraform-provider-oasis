module github.com/arangodb-managed/terraform-provider-oasis

require (
	github.com/arangodb-managed/apis v0.76.0
	github.com/arangodb-managed/log-helper v0.2.5
	github.com/gogo/protobuf v1.3.2
	github.com/hashicorp/terraform-plugin-docs v0.8.1
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.13.0
	github.com/rs/zerolog v1.26.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.47.0
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-checkpoint v0.5.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320 // indirect
	github.com/hashicorp/go-hclog v1.2.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.4.3 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/go-version v1.4.0 // indirect
	github.com/hashicorp/hc-install v0.3.2 // indirect
	github.com/hashicorp/hcl/v2 v2.11.1 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-exec v0.16.1 // indirect
	github.com/hashicorp/terraform-json v0.13.0 // indirect
	github.com/hashicorp/terraform-plugin-go v0.8.0 // indirect
	github.com/hashicorp/terraform-plugin-log v0.3.0 // indirect
	github.com/hashicorp/terraform-registry-address v0.0.0-20210412075316-9b2996cce896 // indirect
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/cli v1.1.3 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/posener/complete v1.1.1 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.1 // indirect
	github.com/zclconf/go-cty v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20211215165025-cf75a172585e // indirect
	golang.org/x/net v0.0.0-20211209124913-491a49abca63 // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220126215142-9970aeb2e350 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

go 1.17

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible

replace github.com/arangodb/kube-arangodb => github.com/arangodb/kube-arangodb v0.0.0-20220224031947-96c0978d52b5

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

replace k8s.io/api => k8s.io/api v0.21.8

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.8

replace k8s.io/apimachinery => k8s.io/apimachinery v0.21.8

replace k8s.io/apiserver => k8s.io/apiserver v0.21.8

replace k8s.io/client-go => k8s.io/client-go v0.21.8

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.8

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.8

replace k8s.io/code-generator => k8s.io/code-generator v0.21.8

replace k8s.io/component-base => k8s.io/component-base v0.21.8

replace k8s.io/kubernetes => k8s.io/kubernetes v1.21.8

replace k8s.io/metrics => k8s.io/metrics v0.21.8

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.9.7

replace sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.0

replace github.com/hashicorp/terraform-plugin-sdk => github.com/hashicorp/terraform-plugin-sdk/v2 v2.13.0

replace github.com/cilium/cilium => github.com/cilium/cilium v1.9.5

replace github.com/optiopay/kafka => github.com/optiopay/kafka v2.0.4+incompatible

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.8

replace k8s.io/cri-api => k8s.io/cri-api v0.21.8

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.8

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.8

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.8

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.8

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.8

replace k8s.io/kubelet => k8s.io/kubelet v0.21.8

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.21.8

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.8

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20201211151036-40ec1c210f7a

replace k8s.io/kubectl => k8s.io/kubectl v0.21.8

replace github.com/nats-io/nats.go => github.com/nats-io/nats.go v1.10.0

replace github.com/nats-io/stan.go => github.com/nats-io/stan.go v0.8.3

replace github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring => github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.47.1

replace github.com/prometheus-operator/prometheus-operator/pkg/client => github.com/prometheus-operator/prometheus-operator/pkg/client v0.47.1

replace github.com/prometheus-operator/prometheus-operator => github.com/prometheus-operator/prometheus-operator v0.47.1

replace go.uber.org/multierr => go.uber.org/multierr v1.6.1-0.20201027220001-0eb6eb5383b9

replace k8s.io/component-helpers => k8s.io/component-helpers v0.21.8

replace k8s.io/controller-manager => k8s.io/controller-manager v0.21.8

replace k8s.io/mount-utils => k8s.io/mount-utils v0.21.8
