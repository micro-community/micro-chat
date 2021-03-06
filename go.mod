module github.com/micro-community/micro-chat

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/micro/dev v0.0.0-20201026103917-a7b0e7877fa5
	github.com/micro/micro/v3 v3.0.0-beta.7.0.20201026143853-bf049ed6c478
	github.com/urfave/cli/v2 v2.2.0
	google.golang.org/protobuf v1.25.0
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.29.0
