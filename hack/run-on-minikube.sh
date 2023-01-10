#!/bin/env sh

## Style Guide: https://google.github.io/styleguide/shellguide.html

DIR=$(dirname $(readlink -f "$0"))
source "$DIR/common.sh" || exit 1

[ -z "$LOCAL_SERVICES" ] && LOCAL_SERVICES="false"
[ -z "$MAX_CONCURRENCY" ] && MAX_CONCURRENCY=4

PROJECT_ROOTS=( frontend/eshop/ services/orders/deploy/* services/orders-events-consumer services/catalog/deploy/* deploy/kubernetes/operators/service-mapper )
SERVICES=( demo-soa-catalog demo-soa-frontend demo-soa-orders demo-soa-orders-events-consumer )
OLM_OPERATOR_INSTALLER_URL=https://github.com/operator-framework/operator-lifecycle-manager/releases/download/v0.22.0/install.sh
LOCAL_SERVICES_IMAGES=( amazon/dynamodb-local:latest pafortin/goaws demo-soa postgres:15.1 )
SERVICE_MAPPER_PATH=$(realpath "$DIR/../deploy/kubernetes/operators/service-mapper")

MANIFESTS="config/manifests/overlays/minikube"$([ "$LOCAL_SERVICES" != "true" ] && echo "-aws")

###################################################
# loads needed images onto the minikube cluster
# Globals:
#   LOCAL_SERVICES: whether to load also local service images
#   PROJECT_ROOTS: where to look for Dockerfiles to parse to get base images
# Arguments:
#   None
###################################################
function load_images_in_minikube {
  local images=()
  [ "$LOCAL_SERVICES" = "true" ] && images+=$LOCAL_SERVICES_IMAGES

  for r in ${PROJECT_ROOTS[@]}; do
      images+=($(cat "$r/Dockerfile" | grep FROM | cut -d' ' -f 2))
  done

  local sorted=()
  IFS=$'\n'; sorted=($(sort -u <<<"${images[*]}")); unset IFS

  local pids=()
  for r in ${sorted[@]}; do
      echo "loading image $r in minikube"
      (time minikube image load --profile demo-soa $r && echo "image $r loaded in minikube")& pids+=("$!")
  done

  wait "${pids[@]}" || true
}

###################################################
# Starts and prepare a minikube cluster
# Globals:
#   PROJECT_ROOTS: where to look for Dockerfiles to parse to get base images
#   OLM_OPERATOR_INSTALLER_URL: url to OLM Operator
#   MAX_CONCURRENCY: max concurrency to use for docker builds
# Arguments:
#   None
###################################################
function prepare_cluster {
    minikube start \
        --insecure-registry 0.0.0.0/0 \
        --profile demo-soa \
        --namespace demo-soa

    local pids=()
    # enable ingress addon
    minikube addons enable ingress --profile demo-soa & pids+=("$!")

    # install olm
    (curl -sL "$OLM_OPERATOR_INSTALLER_URL" | bash -s v0.22.0 && echo "OLM installed" )& pids+=("$!")

    # populate cache for next builds
    local load_pid=''
    (time load_images_in_minikube && "all images loaded in minikube")& load_pid="$!"
    wait "$load_pid" || (echo "error filling minikube's cache" && (kill "${pids[@]}" || true ) && return 1)

    eval $(minikube docker-env --profile demo-soa)

    (cd $SERVICE_MAPPER_PATH && make install docker-build deploy) & pids+=("$!")
    make -j $MAX_CONCURRENCY docker-build-all & pids+=("$!")

    echo "waiting for pids ${pids[@]}"
    wait "${pids[@]}" || (echo "error preparing cluster" && return 1)
}

###################################################
# Install manifests and waits for service rollout
# Globals:
#   SERVICES: services to wait for rollout
#   MANIFESTS: where to find manifests to install
# Arguments:
#   demo-soa app url
###################################################
function install_demo_soa_app {
    # install manifests
    loop_for 10 20 make install MANIFESTS_FOLDER=$MANIFESTS || return 1

    # wait for rollouts and service availability
    loop_for 10 20 check_deployment_rollout "$SERVICES" || return 1
    loop_for 10 20 check_availability "$1" || return 1
}

function main {
    (minikube --profile demo-soa status > /dev/null && echo "cluster demo-soa is still running") || \
        prepare_cluster || (echo "error preparing cluster" && exit 1)

    local minikube_url="https://$(minikube ip --profile demo-soa)"
    install_demo_soa_app "$minikube_url" || (echo "error installing demo soa app" && exit 1)

    open_url "$minikube_url"
}

main || exit 1
