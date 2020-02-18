test:
	go test . -v

bench:
	go test -bench=. -run=^B -v