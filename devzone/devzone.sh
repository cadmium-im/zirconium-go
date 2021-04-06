#!/usr/bin/env bash

# Actually this isn't needed as we extract script path later and use it
# for everything, but anyway!
if [ "$0" != "./devzone.sh" ]; then
    echo "This script should be launched as './devzone.sh'!"
    exit 1
fi

# Check for OS first. On macOS greadlink should be used instead of
# readlink.
READLINK="/bin/readlink"
OS=$(uname -s)
if [ "${OS}" == "Darwin" ]; then
    READLINK="/usr/local/bin/greadlink"
    if [ ! -f "${READLINK}" ]; then
        echo "GNU readlink is required for macOS. Please, install coreutils with brew."
        exit 1
    fi
fi

SCRIPT_PATH=$(dirname "`${READLINK} -f "${BASH_SOURCE}"`")
echo "devzone script path: ${SCRIPT_PATH}"

# This values might or might not be filled by OS.
# And it MUST NOT be filled on macOS. Docker is so convenient...
if [ "${OS}" != "Darwin" ]; then
    USER_ID=$(id -u)
    GROUP_ID=$(id -g)
else
    echo "macOS users have no need in setting user and group ID"
fi

down() {
    echo "Cleaning up development environment..."
    docker-compose down --remove-orphans
}

run() {
    echo "Getting development environment up and running with docker-compose..."
    # ToDo: checks?
    USER_ID=$USER_ID GROUP_ID=$GROUP_ID docker-compose -p devzone_zirconium up
    if [ $? -ne 0 ]; then
        echo "Something went wrong. Read previous messages carefully!"
        exit 1
    fi
    echo "Development zone shutted down."
}


help() {
    echo "Developers helper script."
    echo ""
    echo "Available subcommands:"
    echo -e "\tdown\t\t\t\tClear development environment from data."
    echo -e "\trun\t\t\t\tStart development zone required servers (databases, etc.)."
}

case $1 in
    down)
        down
    ;;
    run)
        run
    ;;
    *)
        help
    ;;
esac
