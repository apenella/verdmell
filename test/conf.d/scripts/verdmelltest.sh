#!/bin/bash

EXITCODE="$1"
WAITTIME="$2"
OUTPUT="$3"

sleep $WAITTIME
echo "$OUTPUT. (exit: $EXITCODE)"

exit $EXITCODE
