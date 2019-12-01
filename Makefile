clean: .clean

build: .build

push: .push

push_hub: .push_hub

build_squzy: .build_squzy

run_squzy: .run_squzy

test: .test

test_d: .test_debug

test_cover: .test_cover

dep: .dep

lint: .lint

default: build

.lint:
	golangci-lint run

.run_squzy:
	bazel run //apps/squzy:squzy_app

.build_squzy:
	bazel build //apps/squzy:squzy

.push:
	bazel run //apps/squzy:squzy_push

.push_hub:
	bazel run //apps/squzy:squzy_push_hub

.build:
	bazel build //apps/...

.test:
	bazel test //apps/...

.test_debug:
	bazel test //apps/...:all --sandbox_debug

.dep:
	bazel run //:gazelle -- update-repos -from_file=go.mod

.test_cover:
	# bazel coverage --test_arg="-test.coverprofile=c.out" //apps/...
	go test ./... -coverprofile=c.out
