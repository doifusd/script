#!/bin/bash

#================================================================
#   Copyright (C) 2023 Sangfor Ltd. All rights reserved.
#   
#   创 建 者：sky
#   创建日期：2023年04月07日
#   描    述：
#
#================================================================


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY="https://goproxy.cn,direct" go build -o "./jf" -a -installsuffix cgo -ldflags '-s -w' ./

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GOPROXY="https://goproxy.cn,direct" go build -o "./jf" -a -installsuffix cgo -ldflags '-s -w' ./
