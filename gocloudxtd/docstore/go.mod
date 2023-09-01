// Ideally we don't need to subpackage this, but this is necessary to avoid circular imports when
// gostd/gocloudxtd/docstore/mongodb is imported.

module github.com/fensak-io/gostd/gocloudxtd/docstore

go 1.19

require gocloud.dev v0.34.0

require (
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/net v0.13.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.134.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230731193218-e0aa005b6bdf // indirect
	google.golang.org/grpc v1.57.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
