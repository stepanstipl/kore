#!/bin/sh

set -e

## For string output from function
detect_service_endpoint_out=""
detect_service_endpoint () {
  external_ip=""
  echo "Finding endpoint for $1"
  while [ -z $external_ip ]; do
    echo "Waiting for end point..."
    external_ip=$(kubectl get svc $1 --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}")
    [ -z "$external_ip" ] && sleep 2
  done
  service_port=$(kubectl get svc $1 --template='{{ (index .spec.ports 0).port }}')
  echo "Found ip: $external_ip"
  echo "Found port: $service_port"

  detect_service_endpoint_out="$external_ip:$service_port"
}

echo "Dectecting service endpoints for $1"
echo "This may take some time..."

command="kubectl create secret generic $1-discovered-endpoints"

# The first var is the release name, we want to skip it
shift

for var in "$@"
do
    detect_service_endpoint $var
    command="$command --from-literal=$var=http://$detect_service_endpoint_out"
done

eval $command
