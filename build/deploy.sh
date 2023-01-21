#!/bin/bash

REGISTRY=cloud.canister.io:5000/elgendy1204/currency-grpc

main_folder_check() {
    if [[ ! -f "./go.mod" ]]; then
        echo "Error: should run from project's main folder"
        exit 1;
    fi
}

main_branch_check() {
    BRANCH=$(git rev-parse --abbrev-ref HEAD)
    if [[ "$BRANCH" != "main" ]]; then
      echo 'Error: deploying on main only!';
      exit 1;
    fi
}

env_file_source() {
    local current_dir="$(dirname "${BASH_SOURCE[0]}")"
    local ENV_FILE="$current_dir/../.env"
    if [[ ! -f "$ENV_FILE" ]]; then
        echo "Error: .env file for production should exists in build folder"
        exit 1;
    fi
    source $ENV_FILE
}

docker_image_build() {
    docker build \
        -t $REGISTRY:latest \
        -f build/Dockerfile .
    docker push $REGISTRY --all-tags
}

server_deploy() {
    ssh $SERVER_USER@$SERVER_ADDRESS <<EOF
    docker pull $REGISTRY:latest

    docker network create --driver bridge application-network || true

    docker container run -d --name migrations \
        --network application-network \
        --entrypoint /bin/sh \
        --rm $REGISTRY:latest \
        -c 'npm run migrations'

    docker stop currency-grpc || true
    docker rm currency-grpc || true
    docker container run -d --name currency-grpc \
        --network application-network \
        --restart unless-stopped \
        $REGISTRY:latest
EOF
}

main_branch_check
main_folder_check
env_file_source
docker_image_build
server_deploy
