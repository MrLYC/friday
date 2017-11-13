#!/bin/sh

find ${TESTLUAROOT} -name '*.lua' | while read lua;
do
    echo ${lua}
    ${TARGET} vm -path ${lua} || exit $?
done