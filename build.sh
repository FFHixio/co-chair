#!/bin/bash

set -e

function generate_cert {
    ./generate-cert.sh
}

function clean_proto {
    rm -f proto/client/* proto/server/* 
}

function clean_static {
    rm -f frontend/static/*
}

function build_proto {
    # the /tmp path is for CI only
    protoc -I. -I/tmp/protobuf/include -Ivendor/ ./proto/web.proto \
        --go_out=plugins=grpc:$GOPATH/src
}

function build_static {
    mkdir -p frontend/static
    # Do the static frontend build, unless if parcel is running.
    # We need the -f option here, since parcel is a node process.
    if ! pgrep -f "parcel" > /dev/null ; then
        # we could pass --public-url to parcel build
        (cd ui && parcel build index.html)
    fi
    cp ui/dist/* frontend/static/
    # TODO(cm) move these steps into assets_generate.go
    mkdir -p frontend/bundle
    touch frontend/bundle/bundle.go
    (cd frontend && go run assets_generate.go)
}

function build_all {

    if [ ! -f ./key.pem ]; then
        generate_cert
    fi

    clean_proto
    build_proto

    clean_static
    build_static

    go build 
    build_container
}

function build_default {
    go build
}

function build_container {
    if [ ! -f ./co-chair ]; then
        go build
    fi
    cp co-chair docker
    pushd docker
        bash container.sh
    popd
}

case "$1" in
    all)
        build_all 
    ;;
    clean)
        clean_static
        clean_proto
    ;;
    container)
        build_container
    ;;
    static)
        clean_static
        build_static
    ;;
    proto)
        clean_proto
        build_proto
    ;;
    help)
        echo "Usage:"
        echo "    ./build.sh [help|clean|all|static|proto]"
    ;;
    *)
        build_default
    ;;
esac
