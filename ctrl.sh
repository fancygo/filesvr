#!/bin/bash
auth="fancy"
svrname=$1
#./mqsvr fancy.mqsvr >> log 2>&1 &
./${svrname} ${auth}.${svrname} >> log 2>&1 &
