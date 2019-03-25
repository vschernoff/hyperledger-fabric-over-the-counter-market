#!/usr/bin/env bash

starttime=$(date +%s)

# defaults; export these variables before executing this script
composeTemplatesFolder="dockercompose-templates"
artifactsTemplatesFolder="artifacts-templates"
: ${FABRIC_STARTER_HOME:=$PWD}
: ${TEMPLATES_ARTIFACTS_FOLDER:=$FABRIC_STARTER_HOME/$artifactsTemplatesFolder}
: ${TEMPLATES_DOCKER_COMPOSE_FOLDER:=$FABRIC_STARTER_HOME/$composeTemplatesFolder}
: ${GENERATED_ARTIFACTS_FOLDER:=./artifacts}
: ${GENERATED_DOCKER_COMPOSE_FOLDER:=./dockercompose}

: ${DOMAIN:="example.com"}
: ${IP_ORDERER:="127.0.0.1"}
: ${ORG1:="a"}
: ${ORG2:="b"}
: ${ORG3:="c"}
: ${PEER0:="peer0"}
: ${PEER1:="peer1"}
: ${MAIN_ORG:=${ORG1}}
: ${IP1:="127.0.0.1"}
: ${IP2:="127.0.0.1"}
: ${IP3:="127.0.0.1"}

: ${FABRIC_VERSION:="1.4.0"}
: ${THIRDPARTY_VERSION:="0.4.14"}
: ${FABRIC_REST_VERSION:="0.13.0"}


echo "Use Fabric-Starter home: $FABRIC_STARTER_HOME"
echo "Use docker compose template folder: $TEMPLATES_DOCKER_COMPOSE_FOLDER"
echo "Use target artifact folder: $GENERATED_ARTIFACTS_FOLDER"
echo "Use target docker-compose folder: $GENERATED_DOCKER_COMPOSE_FOLDER"
echo "Use Fabric Version: $FABRIC_VERSION"
echo "Use Fabric REST Version: $FABRIC_REST_VERSION"
echo "Use 3rdParty Version: $THIRDPARTY_VERSION"


CLI_TIMEOUT=10000

CHAINCODE_VERSION="1.0"
CHAINCODE_COMMON_NAME=reference

CHAINCODE_COMMON_INIT='{"Args":["init"]}'
CHAINCODE_BILATERAL_INIT='{"Args":["init"]}'
COLLECTION_CONFIG='$GOPATH/src/reference/collections_config.json'
#Set default State Database
LITERAL_COUCHDB="couchdb"
LITERAL_LEVELDB="leveldb"
STATE_DATABASE="${LITERAL_COUCHDB}"

DEFAULT_ORDERER_PORT=7050
DEFAULT_WWW_PORT=4000
DEFAULT_API_PORT=8080
DEFAULT_CA_PORT=7054
DEFAULT_PEER0_PORT=7051
DEFAULT_PEER0_EVENT_PORT=7053
DEFAULT_PEER1_PORT=7056
DEFAULT_PEER1_EVENT_PORT=7058
DEFAULT_KAFKA_PORT=9092

DEFAULT_PEER_EXTRA_HOSTS="extra_hosts:[newline]      - orderer.$DOMAIN:$IP_ORDERER"
DEFAULT_CLI_EXTRA_HOSTS="extra_hosts:[newline]      - orderer.$DOMAIN:$IP_ORDERER[newline]      - www.$DOMAIN:$IP_ORDERER[newline]      - www.$ORG1.$DOMAIN:$IP1[newline]      - www.$ORG2.$DOMAIN:$IP2[newline]      - www.$ORG3.$DOMAIN:$IP3"
DEFAULT_API_EXTRA_HOSTS="extra_hosts:[newline]      - orderer.$DOMAIN:$IP_ORDERER[newline]      - peer0.$ORG1.$DOMAIN:$IP1[newline]      - peer0.$ORG2.$DOMAIN:$IP2[newline]      - peer0.$ORG3.$DOMAIN:$IP3[newline]      - peer1.$ORG1.$DOMAIN:$IP1[newline]      - peer1.$ORG2.$DOMAIN:$IP2[newline]      - peer1.$ORG3.$DOMAIN:$IP3"
DEFAULT_API_EXTRA_HOSTS1="extra_hosts:[newline]      - orderer.$DOMAIN:$IP_ORDERER[newline]      - peer0.$ORG2.$DOMAIN:$IP2[newline]      - peer0.$ORG3.$DOMAIN:$IP3[newline]      - peer1.$ORG2.$DOMAIN:$IP2[newline]      - peer1.$ORG3.$DOMAIN:$IP3"
DEFAULT_API_EXTRA_HOSTS2="extra_hosts:[newline]      - orderer.$DOMAIN:$IP_ORDERER[newline]      - peer0.$ORG1.$DOMAIN:$IP1[newline]      - peer0.$ORG3.$DOMAIN:$IP3[newline]      - peer1.$ORG1.$DOMAIN:$IP1[newline]      - peer1.$ORG3.$DOMAIN:$IP3"
DEFAULT_API_EXTRA_HOSTS3="extra_hosts:[newline]      - orderer.$DOMAIN:$IP_ORDERER[newline]      - peer0.$ORG1.$DOMAIN:$IP1[newline]      - peer0.$ORG2.$DOMAIN:$IP2[newline]      - peer1.$ORG1.$DOMAIN:$IP1[newline]      - peer1.$ORG2.$DOMAIN:$IP2"

GID=$(id -g)

function array_orgs_to_json() {
      local arr=("$@");
      local len=${#arr[@]}

      if [[ ${len} -eq 0 ]]; then
        >&2 echo "Error: Length of input array needs to be at least 2.";
         return 1;
      fi

      if [[ $((len%2)) -eq 1 ]]; then
         >&2 echo "Error: Length of input array needs to be even (key/value pairs).";
         return 1;
      fi

      local foo=0;
      for i in "${arr[@]}"; do
          local char="},{"
          if [ $((++foo%2)) -eq 0 ]; then
               char=":";
          fi

          local first="${i:0:1}";  # read first charc

          local app="\\\"$i\\\""

          if [[ "$first" == "^" ]]; then
            app="${i:1}"  # remove first char
          fi

          JSON_ORG="$JSON_ORG$char$app";

      done

      JSON_ORG="\"[{${JSON_ORG:3}}]\"";  # remove the first three chars
}

# Handle MacOS sed
ARCH=$(uname -s | grep Darwin)
SED_OPTS="-i"
if [ "$ARCH" == "Darwin" ]; then
    SED_OPTS="-it"
fi

function setDockerVersions() {
    echo "set Docker image versions for $1"
    sed $SED_OPTS -e "s/FABRIC_VERSION/$FABRIC_VERSION/g" -e "s/THIRDPARTY_VERSION/$THIRDPARTY_VERSION/g" -e "s/FABRIC_REST_VERSION/$FABRIC_REST_VERSION/g" $1
}

function removeUnwantedContainers() {
  docker ps -a -q -f "name=dev-*"|xargs docker rm -f
}

# Delete any images that were generated as a part of this setup
function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | grep "dev\|none\|test-vp\|peer[0-9]-" | awk '{print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "No images available for deletion"
  else
    echo "Removing docker images: $DOCKER_IMAGE_IDS"
    docker rmi -f ${DOCKER_IMAGE_IDS}
  fi
}

function generateArtifacts() {
  [[ -d $GENERATED_ARTIFACTS_FOLDER ]] || mkdir $GENERATED_ARTIFACTS_FOLDER
  [[ -d $GENERATED_DOCKER_COMPOSE_FOLDER ]] || mkdir $GENERATED_DOCKER_COMPOSE_FOLDER
  cp -f "$TEMPLATES_DOCKER_COMPOSE_FOLDER/base.yaml" "$GENERATED_DOCKER_COMPOSE_FOLDER"
  cp -f "$TEMPLATES_DOCKER_COMPOSE_FOLDER/base-intercept.yaml" "$GENERATED_DOCKER_COMPOSE_FOLDER"
  if [[ -d ./$composeTemplatesFolder ]]; then cp -f "./$composeTemplatesFolder/base-intercept.yaml" "$GENERATED_DOCKER_COMPOSE_FOLDER"; fi
}

function removeArtifacts() {
  generateArtifacts
  echo "Removing generated and downloaded artifacts from: $GENERATED_DOCKER_COMPOSE_FOLDER, $GENERATED_ARTIFACTS_FOLDER"
  rm $GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-*.yaml
  rm -rf $GENERATED_ARTIFACTS_FOLDER/crypto-config
  rm -rf $GENERATED_ARTIFACTS_FOLDER/channel
  rm -rf $GENERATED_ARTIFACTS_FOLDER/*block*
  rm -rf www/artifacts && mkdir www/artifacts
  rm -rf $GENERATED_ARTIFACTS_FOLDER/cryptogen-*.yaml
  rm -rf $GENERATED_ARTIFACTS_FOLDER/fabric-ca-server-config-*.yaml
  rm -rf $GENERATED_ARTIFACTS_FOLDER/network-config.json
  rm -rf $GENERATED_ARTIFACTS_FOLDER/configtx.yaml
  rm -rf $GENERATED_ARTIFACTS_FOLDER/*Config.json
  rm -rf $GENERATED_ARTIFACTS_FOLDER/*.pb
  rm -rf $GENERATED_ARTIFACTS_FOLDER/updated_config.*
  rm -rf $GENERATED_ARTIFACTS_FOLDER/update_in_envelope.*
  rm -rf $GENERATED_ARTIFACTS_FOLDER/update.*
  rm -rf $GENERATED_ARTIFACTS_FOLDER/config.*
  rm -rf $GENERATED_ARTIFACTS_FOLDER/hosts
}

function removeDockersWithDomain() {
  search="$DOMAIN"
  docker_ids=$(docker ps -a | grep ${search} | awk '{print $1}')
  if [ -z "$docker_ids" -o "$docker_ids" == " " ]; then
    echo "No docker instances available for deletion with $search"
  else
    echo "Removing docker instances found with $search: $docker_ids"
    docker rm -f ${docker_ids}
  fi

  docker_ids=$(docker volume ls -q | grep ${search})
  if [ -z "$docker_ids" -o "$docker_ids" == " " ]; then
    echo "No docker volumes available for deletion with $search"
  else
    echo "Removing docker volumes found with $search: $docker_ids"
    docker volume rm -f ${docker_ids}
  fi

}

function generateOrdererDockerCompose() {
    mainOrg=$1
    echo "Creating orderer docker compose yaml file with $DOMAIN, $ORG1, $ORG2, $ORG3, $DEFAULT_ORDERER_PORT, $DEFAULT_WWW_PORT"

    compose_template=$TEMPLATES_DOCKER_COMPOSE_FOLDER/docker-composetemplate-orderer.yaml
    f="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-$DOMAIN.yaml"

    cli_extra_hosts=${DEFAULT_CLI_EXTRA_HOSTS}

    sed -e "s/DOMAIN/$DOMAIN/g" -e "s/MAIN_ORG/$mainOrg/g" -e "s/CLI_EXTRA_HOSTS/$cli_extra_hosts/g" -e "s/ORDERER_PORT/$DEFAULT_ORDERER_PORT/g" -e "s/WWW_PORT/$DEFAULT_WWW_PORT/g" -e "s/ORG1/$ORG1/g" -e "s/ORG2/$ORG2/g" -e "s/ORG3/$ORG3/g" ${compose_template} | awk '{gsub(/\[newline\]/, "\n")}1' > ${f}

    setDockerVersions $f
}

function generateOrdererArtifacts() {
    org=$1

    echo "Creating orderer yaml files with $DOMAIN, $ORG1, $ORG2, $ORG3, $DEFAULT_ORDERER_PORT, $DEFAULT_WWW_PORT"

    f="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-$DOMAIN.yaml"

    mkdir -p "$GENERATED_ARTIFACTS_FOLDER/channel"


    if [[ -n "$org" ]]; then
        sed -e "s/DOMAIN/$DOMAIN/g" -e "s/ORG1/$org/g" "$TEMPLATES_ARTIFACTS_FOLDER/configtxtemplate-oneOrg-orderer.yaml" > $GENERATED_ARTIFACTS_FOLDER/configtx.yaml
    else
        # replace in configtx
        sed -e "s/DOMAIN/$DOMAIN/g" -e "s/ORG1/$ORG1/g" -e "s/ORG2/$ORG2/g" -e "s/ORG3/$ORG3/g" $TEMPLATES_ARTIFACTS_FOLDER/configtxtemplate.yaml > $GENERATED_ARTIFACTS_FOLDER/configtx.yaml
    fi
    createChannels=("common")

    for channel_name in ${createChannels[@]}
    do
        echo "Generating channel config transaction for $channel_name"
        docker-compose --file ${f} run --rm -e FABRIC_CFG_PATH=/etc/hyperledger/artifacts "cli.$DOMAIN" configtxgen -profile "$channel_name" -outputCreateChannelTx "./channel/$channel_name.tx" -channelID "$channel_name"
    done

    # replace in cryptogen
    sed -e "s/DOMAIN/$DOMAIN/g" $TEMPLATES_ARTIFACTS_FOLDER/cryptogentemplate-orderer.yaml > "$GENERATED_ARTIFACTS_FOLDER/cryptogen-$DOMAIN.yaml"

    echo "Generating crypto material with cryptogen"

    echo "docker-compose --file ${f} run --rm \"cli.$DOMAIN\" bash -c \"sleep 2 && cryptogen generate --config=cryptogen-$DOMAIN.yaml\""
    docker-compose --file ${f} run --rm "cli.$DOMAIN" bash -c "sleep 2 && cryptogen generate --config=cryptogen-$DOMAIN.yaml"

    echo "Generating orderer genesis block with configtxgen"
    docker-compose --file ${f} run --rm -e FABRIC_CFG_PATH=/etc/hyperledger/artifacts "cli.$DOMAIN" configtxgen -profile OrdererGenesis -outputBlock ./channel/genesis.block

    for channel_name in ${createChannels[@]}
    do
        echo "Generating channel config transaction for $channel_name"
        docker-compose --file ${f} run --rm -e FABRIC_CFG_PATH=/etc/hyperledger/artifacts "cli.$DOMAIN" configtxgen -profile "$channel_name" -outputCreateChannelTx "./channel/$channel_name.tx" -channelID "$channel_name"

        for myorg in ${ORG1} ${ORG2} ${ORG3}
        do
            echo "Generating anchor peers update for ${myorg}"
            docker-compose --file $GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-${myorg}.yaml run --rm -e FABRIC_CFG_PATH=/etc/hyperledger/artifacts "cli.$myorg.$DOMAIN" configtxgen -profile "$channel_name" --outputAnchorPeersUpdate "./channel/${myorg}MSPanchors-$channel_name.tx" -channelID "$channel_name" -asOrg ${myorg}MSP
        done

    done

    echo "changing ownership of channel block files"
    docker-compose --file ${f} run --rm "cli.$DOMAIN" bash -c "chown -R $UID:$GID ."
}

function addHostFiles() {
    org=$1

    mkdir -p $GENERATED_ARTIFACTS_FOLDER/hosts/${org}
    cp $TEMPLATES_ARTIFACTS_FOLDER/default_hosts $GENERATED_ARTIFACTS_FOLDER/hosts/${org}/api_hosts
    cp $TEMPLATES_ARTIFACTS_FOLDER/default_hosts $GENERATED_ARTIFACTS_FOLDER/hosts/${org}/cli_hosts
}

function generatePeerArtifacts() {
    org=$1

    [[ ${#} == 0 ]] && echo "missing required argument -o ORG" && exit 1

    if [ ${#} == 1 ]; then
      # if no port args are passed assume generating for multi host deployment
      peer_extra_hosts=${DEFAULT_PEER_EXTRA_HOSTS}
      cli_extra_hosts=${DEFAULT_CLI_EXTRA_HOSTS}
      if [ ${org} == ${ORG1} ]; then
        api_extra_hosts=${DEFAULT_API_EXTRA_HOSTS1}
      elif [ ${org} == ${ORG2} ]; then
        api_extra_hosts=${DEFAULT_API_EXTRA_HOSTS2}
      elif [ ${org} == ${ORG3} ]; then
        api_extra_hosts=${DEFAULT_API_EXTRA_HOSTS3}
      fi
    fi

    www_port=$2
    ca_port=$3
    peer0_port=$4
    peer0_event_port=$5
    peer1_port=$6
    peer1_event_port=$7

    if [ "${STATE_DATABASE}" == "couchdb" ]; then
      couchdb_port=${8}
    fi

    : ${www_port:=${DEFAULT_WWW_PORT}}
    : ${ca_port:=${DEFAULT_CA_PORT}}
    : ${peer0_port:=${DEFAULT_PEER0_PORT}}
    : ${peer0_event_port:=${DEFAULT_PEER0_EVENT_PORT}}
    : ${peer1_port:=${DEFAULT_PEER1_PORT}}
    : ${peer1_event_port:=${DEFAULT_PEER1_EVENT_PORT}}
    if [ "${STATE_DATABASE}" == "couchdb" ]; then
    echo "Creating peer yaml files with $DOMAIN, $org, $www_port, $ca_port, $peer0_port, $peer0_event_port, $peer1_port, $peer1_event_port, $couchdb_port"
    compose_template=$TEMPLATES_DOCKER_COMPOSE_FOLDER/docker-composetemplate-peer-couchdb.yaml
    else
    echo "Creating peer yaml files with $DOMAIN, $org, $www_port, $ca_port, $peer0_port, $peer0_event_port, $peer1_port, $peer1_event_port"
    compose_template=$TEMPLATES_DOCKER_COMPOSE_FOLDER/docker-composetemplate-peer.yaml
    fi

    f="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-$org.yaml"

    # cryptogen yaml
    sed -e "s/DOMAIN/$DOMAIN/g" -e "s/ORG/$org/g" $TEMPLATES_ARTIFACTS_FOLDER/cryptogentemplate-peer.yaml > $GENERATED_ARTIFACTS_FOLDER/"cryptogen-$org.yaml"

    # nginx proxy config
    sed -e "s/DOMAIN/$DOMAIN/g" \
        -e "s/ORG/$org/g" \
        -e "s/API_PORT/${DEFAULT_API_PORT}/g" \
        $TEMPLATES_ARTIFACTS_FOLDER/nginx.conf > $GENERATED_ARTIFACTS_FOLDER/"nginx-$org.conf"

    # API configs
    mkdir ${GENERATED_ARTIFACTS_FOLDER}/api-configs-${org}
    sed -e "s/DOMAIN/$DOMAIN/g" \
        -e "s/ORG/$org/g" \
        $TEMPLATES_ARTIFACTS_FOLDER/api-configs/api.yaml > ${GENERATED_ARTIFACTS_FOLDER}/api-configs-${org}/api.yaml
    sed -e "s/DOMAIN/$DOMAIN/g" \
        -e "s/ORG/$org/g" \
        $TEMPLATES_ARTIFACTS_FOLDER/api-configs/network.yaml > ${GENERATED_ARTIFACTS_FOLDER}/api-configs-${org}/network.yaml

    # docker-compose yaml
    if [ "${STATE_DATABASE}" == "couchdb" ]; then
      sed -e "s/PEER_EXTRA_HOSTS/$peer_extra_hosts/g" \
          -e "s/CLI_EXTRA_HOSTS/$cli_extra_hosts/g" \
          -e "s/API_EXTRA_HOSTS/$api_extra_hosts/g" \
          -e "s/DOMAIN/$DOMAIN/g" \
          -e "s/\([^ ]\)ORG/\1$org/g" \
          -e "s/WWW_PORT/$www_port/g" \
          -e "s/CA_PORT/$ca_port/g" \
          -e "s/PEER0_PORT/$peer0_port/g" \
          -e "s/PEER0_EVENT_PORT/$peer0_event_port/g" \
          -e "s/PEER1_PORT/$peer1_port/g" \
          -e "s/PEER1_EVENT_PORT/$peer1_event_port/g" \
          -e "s/COUCHDB_PORT/$couchdb_port/g" \
          ${compose_template} | awk '{gsub(/\[newline\]/, "\n")}1' > ${f}
    else
      sed -e "s/PEER_EXTRA_HOSTS/$peer_extra_hosts/g" \
          -e "s/CLI_EXTRA_HOSTS/$cli_extra_hosts/g" \
          -e "s/API_EXTRA_HOSTS/$api_extra_hosts/g" \
          -e "s/DOMAIN/$DOMAIN/g" \
          -e "s/\([^ ]\)ORG/\1$org/g" \
          -e "s/WWW_PORT/$www_port/g" \
          -e "s/CA_PORT/$ca_port/g" \
          -e "s/PEER0_PORT/$peer0_port/g" \
          -e "s/PEER0_EVENT_PORT/$peer0_event_port/g" \
          -e "s/PEER1_PORT/$peer1_port/g" \
          -e "s/PEER1_EVENT_PORT/$peer1_event_port/g" \
          ${compose_template} | awk '{gsub(/\[newline\]/, "\n")}1' > ${f}
    fi

    # fabric-ca-server-config yaml
    sed -e "s/ORG/$org/g" $TEMPLATES_ARTIFACTS_FOLDER/fabric-ca-server-configtemplate.yaml > $GENERATED_ARTIFACTS_FOLDER/"fabric-ca-server-config-$org.yaml"

    addHostFiles "${org}"
    setDockerVersions $f

    echo "Generating crypto material with cryptogen"

    echo "docker-compose --file ${f} run --rm \"cliNoCryptoVolume.$org.$DOMAIN\" bash -c \"cryptogen generate --config=cryptogen-$org.yaml\""
    docker-compose --file ${f} run --rm "cliNoCryptoVolume.$org.$DOMAIN" bash -c "sleep 2 && cryptogen generate --config=cryptogen-$org.yaml"

    echo "Changing artifacts ownership"
    docker-compose --file ${f} run --rm "cliNoCryptoVolume.$org.$DOMAIN" bash -c "chown -R $UID:$GID ."

    echo "Adding generated CA private keys filenames to $f"
    ca_private_key=$(basename `ls -t $GENERATED_ARTIFACTS_FOLDER/crypto-config/peerOrganizations/"$org.$DOMAIN"/ca/*_sk`)
    [[ -z  ${ca_private_key}  ]] && echo "empty CA private key" && exit 1
    sed $SED_OPTS -e "s/CA_PRIVATE_KEY/${ca_private_key}/g" ${f}

    # replace in configtx
    sed -e "s/DOMAIN/$DOMAIN/g" -e "s/ORG/$org/g" $TEMPLATES_ARTIFACTS_FOLDER/configtx-orgtemplate.yaml > $GENERATED_ARTIFACTS_FOLDER/configtx.yaml

    echo "Generating ${org}Config.json"
    echo "docker-compose --file ${f} run --rm \"cliNoCryptoVolume.$org.$DOMAIN\" bash -c \"FABRIC_CFG_PATH=./ configtxgen  -printOrg ${org}MSP > ${org}Config.json\""
    docker-compose --file ${f} run --rm "cliNoCryptoVolume.$org.$DOMAIN" bash -c "FABRIC_CFG_PATH=./ configtxgen  -printOrg ${org}MSP > ${org}Config.json"
}

function createChannel () {
    org=$1
    channel_name=$2
    f="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-${org}.yaml"

    info "creating channel $channel_name by $org using $f"

    echo "docker-compose --file ${f} run --rm \"cli.$org.$DOMAIN\" bash -c \"peer channel create -o orderer.$DOMAIN:7050 -c $channel_name -f /etc/hyperledger/artifacts/channel/$channel_name.tx --tls --cafile /etc/hyperledger/crypto/orderer/tls/ca.crt\""
    docker-compose --file ${f} run --rm "cli.$org.$DOMAIN" bash -c "peer channel create -o orderer.$DOMAIN:7050 -c $channel_name -f /etc/hyperledger/artifacts/channel/$channel_name.tx --tls --cafile /etc/hyperledger/crypto/orderer/tls/ca.crt"

    sleep 2

    echo "docker-compose --file ${f} run --rm \"cli.$org.$DOMAIN\" bash -c \"peer channel update -o orderer.$DOMAIN:7050 -c $channel_name -f /etc/hyperledger/artifacts/channel/${org}MSPanchors-$channel_name.tx --tls --cafile /etc/hyperledger/crypto/orderer/tls/ca.crt\""
    docker-compose --file ${f} run --rm "cli.$org.$DOMAIN" bash -c "peer channel update -o orderer.$DOMAIN:7050 -c $channel_name -f /etc/hyperledger/artifacts/channel/${org}MSPanchors-$channel_name.tx --tls --cafile /etc/hyperledger/crypto/orderer/tls/ca.crt"

    echo "changing ownership of channel block files"
    docker-compose --file ${f} run --rm "cli.$DOMAIN" bash -c "chown -R $UID:$GID ."

    d="$GENERATED_ARTIFACTS_FOLDER"
    echo "copying channel block file from ${d} to be served by www.$org.$DOMAIN"
    cp "${d}/$channel_name.block" "www/${d}"
}

function joinChannel() {
    org=$1
    channel_name=$2
    f="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-${org}.yaml"

    info "joining channel $channel_name by all peers of $org using $f"

    docker-compose --file ${f} run --rm "cli.$org.$DOMAIN" bash -c "CORE_PEER_ADDRESS=peer0.$org.$DOMAIN:7051 peer channel join -b $channel_name.block"
    docker-compose --file ${f} run --rm "cli.$org.$DOMAIN" bash -c "CORE_PEER_ADDRESS=peer1.$org.$DOMAIN:7051 peer channel join -b $channel_name.block"
}

function instantiateChaincode () {
    sleep 3
    org=$1
    channel_names=($2)
    n=$3
    i=$4
    cc=$5
    if [ -z "$cc" ]; then
    cc="";
    else
    cc="--collections-config $cc";
    fi

    f="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-${org}.yaml"

    for channel_name in ${channel_names[@]}; do
        info "instantiating chaincode $n on $channel_name by $org using $f with $i"

        c="CORE_PEER_ADDRESS=peer0.$org.$DOMAIN:7051 peer chaincode instantiate -n $n -v ${CHAINCODE_VERSION} -c '$i' -o orderer.$DOMAIN:7050 -C $channel_name  $cc --tls --cafile /etc/hyperledger/crypto/orderer/tls/ca.crt"
        d="cli.$org.$DOMAIN"

        echo "instantiating with $d by $c"
        docker-compose --file ${f} run --rm ${d} bash -c "${c}"
    done
}

function installChaincode() {
    org=$1
    n=$2
    v=$3
    f="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-${org}.yaml"
    # chaincode path is the same as chaincode name by convention: code of chaincode instruction lives in ./chaincode/go/reference mapped to docker path /opt/gopath/src/reference
    p=${n}

    lang=golang

    info "installing chaincode $n to peers of $org from ./chaincode/go/$p $v using $f"

    echo "docker-compose --file ${f} run --rm \"cli.$org.$DOMAIN\" bash -c \"CORE_PEER_ADDRESS=peer0.$org.$DOMAIN:7051 peer chaincode install -n $n -v $v -p $p -l $lang "
    echo " && CORE_PEER_ADDRESS=peer1.$org.$DOMAIN:7051 peer chaincode install -n $n -v $v -p $p -l $lang\""

    docker-compose --file ${f} run --rm "cli.$org.$DOMAIN" bash -c "CORE_PEER_ADDRESS=peer0.$org.$DOMAIN:7051 peer chaincode install -n $n -v $v -p $p -l $lang \
    && CORE_PEER_ADDRESS=peer1.$org.$DOMAIN:7051 peer chaincode install -n $n -v $v -p $p -l $lang"
}

function dockerComposeUp () {
  compose_file="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-$1.yaml"

  info "starting docker instances from $compose_file"

  TIMEOUT=${CLI_TIMEOUT} docker-compose -f ${compose_file} up -d 2>&1
  if [ $? -ne 0 ]; then
    echo "ERROR !!!! Unable to start network"
    logs ${1}
    exit 1
  fi
}

function dockerComposeDown () {
  compose_file="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-$1.yaml"

  if [ -f ${compose_file} ]; then
      info "stopping docker instances from $compose_file"
      docker-compose -f ${compose_file} down
  fi;
}

function dockerContainerRestart () {
  org=$1
  service=$2

  compose_file="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-$org.yaml"

  echo "Restart container: docker-compose -f ${compose_file} restart $service.$org.$DOMAIN"
  docker-compose -f ${compose_file} restart $service.$org.$DOMAIN
}

function installAll() {
  org=$1

  sleep 2

  for chaincode_name in ${CHAINCODE_COMMON_NAME}
  do
    installChaincode "${org}" "${chaincode_name}" "${CHAINCODE_VERSION}"
  done
}

function createJoinInstantiate() {
  org=${1}
  channel_name=${2}
  chaincode_name=${3}
  chaincode_init=${4}
  collections=${5}

  createChannel ${org} ${channel_name}
  joinChannel ${org} ${channel_name}
  instantiateChaincode ${org} ${channel_name} ${chaincode_name} ${chaincode_init} ${collections}
}

function makeCertDirs() {
   for certDirsOrg in "$@"
    do
        d="$GENERATED_ARTIFACTS_FOLDER/crypto-config/peerOrganizations/$certDirsOrg.$DOMAIN/peers/peer0.$certDirsOrg.$DOMAIN/tls"
        echo "mkdir -p ${d}"
        mkdir -p ${d}
    done
}

function info() {
    echo "*************************************************************************************************************"
    echo "$1"
    echo "*************************************************************************************************************"
}

function logs () {
  f="$GENERATED_DOCKER_COMPOSE_FOLDER/docker-compose-$1.yaml"

  TIMEOUT=${CLI_TIMEOUT} COMPOSE_HTTP_TIMEOUT=${CLI_TIMEOUT} docker-compose -f ${f} logs -f
}

function clean() {
  removeDockersWithDomain
  #removeUnwantedImages
}

function printArgs() {
  echo "$DOMAIN, $ORG1, $ORG2, $ORG3, $IP1, $IP2, $IP3"
}

function checkDocker() {
  if [ -n "$(which docker-compose)" ]; then
    echo "Found docker-compose, skipping install.."
  else
    if [ -n "$(which yum)" ]; then
      sh $FABRIC_STARTER_HOME/init-docker-centos.sh
    else
      sh $FABRIC_STARTER_HOME/init-docker.sh
    fi
    return
  fi
}

# Print the usage message
function printHelp () {
  echo "Usage: "
  echo "  network.sh -m up|down|restart|generate"
  echo "  network.sh -h|--help (print this message)"
  echo "    -m <mode> - one of 'up', 'down', 'restart' or 'generate'"
  echo "      - 'up' - bring up the network with docker-compose up"
  echo "      - 'down' - clear the network with docker-compose down"
  echo "      - 'restart' - restart the network"
  echo "      - 'generate' - generate required certificates and genesis block"
  echo "      - 'logs' - print and follow all docker instances log files"
  echo
  echo "    -s <state> - one of 'couchdb' or 'leveldb'"
  echo "      - 'couchdb' - set CouchDB as State Database"
  echo "      - 'leveldb' - set LevelDB as State Database"

  echo "Typically, one would first generate the required certificates and "
  echo "genesis block, then bring up the network. e.g.:"
  echo
  echo "	sudo network.sh -m generate"
  echo "	network.sh -m up"
  echo "	network.sh -m logs"
  echo "	network.sh -m down"
}

# Parse commandline args
while getopts "h?m:o:a:w:c:0:1:2:3:k:v:i:n:M:I:R:P:s:" opt; do
  case "$opt" in
    h|\?)
      printHelp
      exit 0
    ;;
    m)  MODE=$OPTARG
    ;;
    s)  STATE_DATABASE=$OPTARG
    ;;
    v)  CHAINCODE_VERSION=$OPTARG
    ;;
    o)  ORG=$OPTARG
    ;;
    M)  MAIN_ORG=$OPTARG
    ;;
    a)  API_PORT=$OPTARG
    ;;
    w)  WWW_PORT=$OPTARG
    ;;
    c)  CA_PORT=$OPTARG
    ;;
    0)  PEER0_PORT=$OPTARG
    ;;
    1)  PEER0_EVENT_PORT=$OPTARG
    ;;
    2)  PEER1_PORT=$OPTARG
    ;;
    3)  PEER1_EVENT_PORT=$OPTARG
    ;;
    k)  CHANNELS=$OPTARG
    ;;
    i) IP=$OPTARG
    ;;
    n) CHAINCODE=$OPTARG
    ;;
    I) CHAINCODE_INIT_ARG=$OPTARG
    ;;
    R) REMOTE_ORG=$OPTARG
    ;;
    P) ENDORSEMENT_POLICY=$OPTARG
    ;;
  esac
done

checkDocker

if [ "${MODE}" == "up" -a "${ORG}" == "" ]; then

  #Building array from Organizations
  JSON_ORG="";
  ARRAY_ORG=()
  for org in ${ORG3} ${ORG1} ${ORG2}
  do
    ARRAY_ORG+=("name" ${org})
  done

  #Building $JSON_ORG
  array_orgs_to_json "${ARRAY_ORG[@]}"

  CHAINCODE_COMMON_INIT='{"Args":["init",'"${JSON_ORG}"']}'

  for org in ${DOMAIN} ${ORG1} ${ORG2} ${ORG3}
  do
    dockerComposeUp ${org}
    sleep 2
  done

  for org in ${ORG1} ${ORG2} ${ORG3}
  do
    installAll ${org}
  done

  createJoinInstantiate ${ORG1} common ${CHAINCODE_COMMON_NAME} ${CHAINCODE_COMMON_INIT} ${COLLECTION_CONFIG}

  joinChannel ${ORG2} common

  joinChannel ${ORG3} common

elif [ "${MODE}" == "down" ]; then
  for org in ${DOMAIN} ${ORG1} ${ORG2} ${ORG3}
  do
    dockerComposeDown ${org}
  done

  removeUnwantedContainers
  #removeUnwantedImages
  removeDockersWithDomain

elif [ "${MODE}" == "clean" ]; then
  clean

elif [ "${MODE}" == "generate" ]; then
  clean
  removeArtifacts

  [[ -d $GENERATED_ARTIFACTS_FOLDER ]] || mkdir $GENERATED_ARTIFACTS_FOLDER
  [[ -d $GENERATED_DOCKER_COMPOSE_FOLDER ]] || mkdir $GENERATED_DOCKER_COMPOSE_FOLDER
  cp -f "$TEMPLATES_DOCKER_COMPOSE_FOLDER/base.yaml" "$GENERATED_DOCKER_COMPOSE_FOLDER"
  cp -f "$TEMPLATES_DOCKER_COMPOSE_FOLDER/base-intercept.yaml" "$GENERATED_DOCKER_COMPOSE_FOLDER"
  if [[ -d ./$composeTemplatesFolder ]]; then cp -f "./$composeTemplatesFolder/base-intercept.yaml" "$GENERATED_DOCKER_COMPOSE_FOLDER"; fi

  file_base="$GENERATED_DOCKER_COMPOSE_FOLDER/base.yaml"
  file_base_intercept="$GENERATED_DOCKER_COMPOSE_FOLDER/base-intercept.yaml"

  setDockerVersions $file_base
  setDockerVersions $file_base_intercept

  #                     org     www_port ca_port peer0_port peer0_event_port peer1_port peer1_event_port couchdb_port
  generatePeerArtifacts ${ORG1} 4000     7054    7051       7053             7056       7058             5984
  generatePeerArtifacts ${ORG2} 4001     8054    8051       8053             8056       8058             6984
  generatePeerArtifacts ${ORG3} 4002     9054    9051       9053             9056       9058             7984
  generateOrdererDockerCompose ${ORG1}
  generateOrdererArtifacts

elif [ "${MODE}" == "restart-api" ]; then # params:  -o ORG
  dockerContainerRestart $ORG api

elif [ "${MODE}" == "logs" ]; then
  logs ${ORG}

elif [ "${MODE}" == "printArgs" ]; then
  printArgs

elif [ "${MODE}" == "removeArtifacts" ]; then
  removeArtifacts

else
  printHelp
  exit 1
fi

endtime=$(date +%s)
info "Finished in $(($endtime - $starttime)) seconds"

# Delete yamlt files in MacOS
if [ "$ARCH" == "Darwin" ]; then
    find . -name "*.yamlt" -type f -delete
fi
