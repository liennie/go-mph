#!/bin/bash

name=$1
if [ -z "$name" ]; then
	echo no name
	exit
fi

go test -benchmem -run=^$ -bench ^BenchmarkNew$ github.com/liennie/go-mph -count=100 | tee "$name.new.small.txt"
go test -benchmem -run=^$ -bench ^BenchmarkNew$ github.com/liennie/go-mph -count=100 -keys=bigkeys.txt | tee "$name.new.big.txt"
go test -run=^$ -bench ^BenchmarkMPH$ github.com/liennie/go-mph -count=100 | tee "$name.query.small.txt"
go test -run=^$ -bench ^BenchmarkMPH$ github.com/liennie/go-mph -count=100 -keys=bigkeys.txt | tee "$name.query.big.txt"
