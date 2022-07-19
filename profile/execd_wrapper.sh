#!/bin/bash

# Source the .profile script
source .profile
# Iterate through the list of exported env vars
for var in $(compgen -e); do
	# quote the output values replacing double quotes with an escaped version
	value=$(echo ${!var} | sed 's/"/\\"/g')
	# output the variables to FD3
    echo "${var} = \"${value}\"" 1>&3
done