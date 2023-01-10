#!/bin/env sh

# commong functions are stored in this file and sourced on need

###################################################
# Asks user for confirmation
# Globals:
#   None
# Arguments:
#   Question to ask
###################################################
function ask_confirm {
    echo -n "$1 [y/N]? " ; read reply; case $reply in Y*|y*) true ;; *) false ;; esac 
}

###################################################
# Loops until the provided command is successful or
#   maximum number of retries have been attempted
# Globals:
#   None
# Arguments:
#   Maximum number of retry attempt
#   Wait time between an attempt and another
##################################################
function loop_for {
    local max_retry=$1
    local wait_time=$2

    local retry=1
    until ${@:3} || (( retry++ >= $max_retry )); do
        echo "Retrying ${retry}/${max_retry} in $wait_time sec..."
        sleep $wait_time
    done
    ${@:3}
}

###################################################
# Opens the provided URL with xdg-open or prints out the url to open
#  ! Only tested on Fedora 37
#  ! No input sanitization
# Globals:
#   None
# Arguments:
#   URL to be opened
###################################################
function open_url {
    [ $(command -v xdg-open) > /dev/null ] && xdg-open "$1" || echo "Open $1"
}

###################################################
# Waits for deployments rollout
# Globals:
#   None
# Arguments:
#   deployments to monitor
###################################################
function check_deployment_rollout {
    kubectl rollout status deployment "$@"
}

###################################################
# Waits for website availability
# Globals:
#   None
# Arguments:
#   website url
###################################################
function check_availability {
    curl -s -k -f -X GET "$1" > /dev/null
}

