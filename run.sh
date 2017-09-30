#!/bin/sh
echo "creating hamming code"
go build main.go
echo $1
./main $1
echo "creating pieces"
python3 pieces.py
