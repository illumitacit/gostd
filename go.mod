module github.com/fensak-io/gostd

go 1.19

replace github.com/fensak-io/gostd/gocloudxtd/docstore => ./gocloudxtd/docstore

require (
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.2.2
	github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus v1.2.1
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/alexedwards/scs/v2 v2.5.1
	github.com/coreos/go-oidc/v3 v3.5.0
	github.com/fensak-io/httpzaplog v0.1.3
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-resty/resty/v2 v2.7.0
	github.com/hashicorp/go-getter/gcs/v2 v2.2.1
	github.com/hashicorp/go-getter/s3/v2 v2.2.1
	github.com/hashicorp/go-getter/v2 v2.2.1
	github.com/hashicorp/hcl/v2 v2.16.2
	github.com/hashicorp/terraform-config-inspect v0.0.0-20221020162138-81db043ad408
	github.com/hashicorp/terraform-registry-address v0.2.0
	github.com/in-toto/in-toto-golang v0.7.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo/v2 v2.9.7
	github.com/onsi/gomega v1.27.7
	github.com/ory/nosurf v1.2.7
	github.com/rabbitmq/amqp091-go v1.8.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.15.0
	go.uber.org/zap v1.24.0
	gocloud.dev v0.29.0
	gocloud.dev/pubsub/rabbitpubsub v0.29.0
	golang.org/x/oauth2 v0.7.0
	google.golang.org/protobuf v1.28.1
)

require (
	cloud.google.com/go v0.109.0 // indirect
	cloud.google.com/go/compute v1.18.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v0.10.0 // indirect
	cloud.google.com/go/storage v1.29.0 // indirect
	github.com/Azure/azure-amqp-common-go/v3 v3.2.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.2.0 // indirect
	github.com/Azure/go-amqp v0.18.1 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v0.9.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go v1.44.200 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/pprof v0.0.0-20230111200839-76d1ae5aea2b // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-svchost v0.0.1 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.5.0 // indirect
	github.com/shibumi/go-pathspec v1.3.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/zclconf/go-cty v1.12.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sync v0.2.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.9.1 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.110.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230209215440-0dfe4f8abfcc // indirect
	google.golang.org/grpc v1.53.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
