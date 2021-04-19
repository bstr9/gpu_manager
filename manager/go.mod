module gpu_manager/manager

replace gpu_manager/proto => ../proto

go 1.16

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/onsi/ginkgo v1.16.1 // indirect
	github.com/onsi/gomega v1.11.0 // indirect
	gpu_manager/proto v0.0.0-00010101000000-000000000000
)
