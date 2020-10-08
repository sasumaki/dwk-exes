#!/bin/bash
cd "${0%/*}"
tagversion=$1
if [[ $tagversion -eq 0 ]] ; then
  tagversion=$(uuidgen)
fi
echo "building sasumaki/ponger:$tagversion"
docker build -t sasumaki/ponger:$tagversion ..
docker push sasumaki/ponger:$tagversion

sed -i '' -e "s/sasumaki\/ponger:.*/sasumaki\/ponger:$tagversion/g" ../manifests/deployment.yaml

kubectl apply -f ../manifests/deployment.yaml
kubectl apply -f ../manifests/ingress.yaml
kubectl apply -f ../manifests/service.yaml

echo "DEPLOYMENT DONE"
