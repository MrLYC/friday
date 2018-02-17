#!/bin/sh

timeout 1 ${TARGET} run

code=$?

if [ "${code}" = 124 ] 
then
    code=0
fi

exit ${code}
