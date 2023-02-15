perf:
	powerprofilesctl set performance

update:
	go get -u -t ./...

bench:
	go test -count 5 -test.bench=.* > bench.txt