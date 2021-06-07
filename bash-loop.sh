#!/bin/bash

for i in {20..50}
do 
  echo === 
  echo $i 
  echo ===
  go run create-lint-problems.go $i 
  time jsonnet-lint caller.jsonnet
done

