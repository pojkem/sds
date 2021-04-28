module github.com/stratosnet/sds

go 1.14

require (
	github.com/HuKeping/rbtree v0.0.0-20210106022122-8ad34838eb2b
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/alex023/clock v0.0.0-20191208111215-c265f1b2ab18
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/cosmos/go-bip39 v1.0.0
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.4
	github.com/gorilla/mux v1.8.0
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/onsi/ginkgo v1.16.1 // indirect
	github.com/onsi/gomega v1.11.0 // indirect
	github.com/pborman/uuid v1.2.1
	github.com/peterh/liner v1.2.1
	github.com/shirou/gopsutil v3.20.12+incompatible
	github.com/tendermint/classic v0.0.0-20201012085102-0a11024b2668
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/sys v0.0.0-20210113181707-4bcb84eeeb78 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1