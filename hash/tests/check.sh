#!/bin/bash

shasum -c $1 1>/dev/null

if [ $? = 0 ]; then
 echo "ok"
else
 echo "something wrong"
fi


