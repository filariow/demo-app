#!/bin/env sh

## Style Guide: https://google.github.io/styleguide/shellguide.html

DIR=$(dirname $(readlink -f "$0"))
source "$DIR/common.sh" || exit 1

[ -z "$LOCAL_SERVICES" ] && LOCAL_SERVICES="false"
[ -z "$MAX_CONCURRENCY" ] && MAX_CONCURRENCY=4
[ -z "$NAMESPACE" ] && NAMESPACE="demo-soa"

MANIFESTS="config/manifests/overlays/crc"$([ "$LOCAL_SERVICES" != "true" ] && echo "-aws")

SERVICE_MAPPER_PATH=$(realpath "$DIR/../deploy/kubernetes/operators/service-mapper")
SERVICES=( demo-soa-catalog demo-soa-frontend demo-soa-orders demo-soa-orders-events-consumer )
NAMESPACE_SERVICE_MAPPER="service-mapper-system"

###################################################
# checks if Host's Docker is configured correctly to push images onto Local Openshift insecure registry
#  and helps configure Docker Daemon. (Only tested on Fedora 37)
# Globals:
#   None
# Arguments:
#   URL to use for connecting to rtarget Openshift
###################################################
interactively_configure_insecure_registry() {
    # check for local openshift's insecure-regisry in docker's daemon.json
    local docker_daemon_file="/etc/docker/daemon.json"
    local openshift_local_url=$1

    if [ ! -f "$docker_daemon_file" ]; then
        cat << EOF
"$docker_daemon_file" not found. Create it with the following content:

$(echo "{}" | jq '."insecure-registry" += ["'$openshift_local_url'"]')
EOF

        ask_confirm "Do you want me to create the file $docker_daemon_filE (needs sudo)" && \
            (echo "{}" | jq '."insecure-registry" += ["'$openshift_local_url'"]' | sudo tee $docker_daemon_file) || \
            exit 1
    fi

    fr=$(cat $docker_daemon_file |jq '."insecure-registries" as $r | "'$openshift_local_url'" | IN($r[])')
    if [ "$fr" = "false" ]; then
        ask_confirm "Do you want me to enrich your '$docker_daemon_file' with insecure registRY '$openshift_local_url' (needs sudo)" || exit 1

        sudo cp "$docker_daemon_file{,.bkp}"
        cat $docker_daemon_file | jq '."insecure-registry += ["'$openshift_local_url'"]' | sudo tee $docker_daemon_file

        ask_confirm "Docker service needs to be restarted. Do you want me to restart it (relies on systemd and sudo)" || exit 1
        sudo systemctl restart docker.{service,socket}
    fi
}

prepare_cluster() {
    ## start the OpenShift cluster
    crc start || return 1
    oc login -u kubeadmin "https://api.crc.testing:6443" || return 1
    oc patch configs.imageregistry.operator.openshift.io/cluster --patch '{"spec":{"defaultRoute":true}}' --type=merge || return 1
    local openshift_local_url=$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')
    interactively_configure_insecure_registry "$openshift_local_url"

    ## create target namespace where to push images
    kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
    kubectl config set-context --current --namespace $NAMESPACE
}

install_service_mapper() {
    local openshift_local_url=$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')

    oc registry login || return 1

    kubectl create namespace $NAMESPACE_SERVICE_MAPPER --dry-run=client -o yaml | kubectl apply -f -

    ( cd "$SERVICE_MAPPER_PATH" && \
        make docker-build docker-push IMG="$openshift_local_url/$NAMESPACE_SERVICE_MAPPER/service-mapper:latest" && \
        make install deploy IMG="$openshift_local_url/$NAMESPACE_SERVICE_MAPPER/service-mapper:latest" )
    kubectl rollout status deployment service-mapper-controller-manager -n $NAMESPACE_SERVICE_MAPPER
}

build_and_push_demo() {
    local openshift_local_url=$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')

    oc registry login || return 1
    make -j $MAX_CONCURRENCY docker-push-all REPOSITORY_REF=$openshift_local_url/$NAMESPACE || \
        (echo "error building and pushing images into OpenShift Local cluster" && return 1)
}

main() {
    [ $(crc status > /dev/null) ] || prepare_cluster || return 1

    [ "$LOCAL_SERVICES" == "false" ] && ( install_service_mapper || return 1 )

    ## install the demo app
    build_and_push_demo || return 1
    loop_for 20 20 make install MANIFESTS_FOLDER=$MANIFESTS || return 1

    ## wait for rollout
    loop_for 20 20 check_deployment_rollout "$SERVICES" || return 1

    ## open url
    open_url "https://demo-soa.apps-crc.testing"
}

main || exit 1
