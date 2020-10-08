#!/bin/bash
cd "${0%/*}"
tagversion=$1
if [[ $tagversion -eq 0 ]] ; then
  tagversion=$(uuidgen)
fi
echo "building sasumaki/reader:$tagversion"
docker build -t sasumaki/reader:$tagversion ../reader
docker push sasumaki/reader:$tagversion

sed -i '' -e "s/sasumaki\/reader:.*/sasumaki\/reader:$tagversion/g" ../manifests/deployment.yaml

echo "sasumaki/reader:$tagversion done"

echo "building sasumaki/hashemitter:$tagversion"
docker build -t sasumaki/hashemitter:$tagversion ../emitter
docker push sasumaki/hashemitter:$tagversion

sed -i '' -e "s/sasumaki\/hashemitter:.*/sasumaki\/hashemitter:$tagversion/g" ../manifests/deployment.yaml

echo "sasumaki/hashemitter:$tagversion done"
kubectl apply -f ../manifests/deployment.yaml
kubectl apply -f ../manifests/ingress.yaml
kubectl apply -f ../manifests/service.yaml

echo "DEPLOYMENT DONE"
