default:
	echo 'Hello, world!'
build:
        @go build -o ./bin/rattle-cli ./cmd/rattle-cli
start0:
        ./bin/rattle-cli start --port 7800 --cluster-port 4000
start1:
        ./bin/rattle-cli start --port 7801 --cluster-port 4001
start2:
        ./bin/rattle-cli start --port 7802 --cluster-port 4002


