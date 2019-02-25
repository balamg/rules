#!/bin/bash
go build -o capi.so -buildmode=c-shared github.com/project-flogo/rules/capi
