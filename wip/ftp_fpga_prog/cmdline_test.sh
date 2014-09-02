#!/bin/bash

echo "1. Test for -v" 
go run ffprog.go -v

if [ "$?" -ne "0" ]; then
  echo "Не понял флаг -v"
  exit 1
fi


echo "1. Test for -debug=config" 
go run ffprog.go -debug=config

if [ "$?" -ne "0" ]; then
  echo "Не понял флаг -debug"
  exit 1
fi

echo "1. Test for json config" 
go run ffprog.go -debug=config -json=./test_config.json

if [ "$?" -ne "0" ]; then
  echo "Не понял флаг -debug"
  exit 1
fi

