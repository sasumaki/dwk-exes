#!/bin/bash
cd "${0%/*}"
tagversion=$1
if [[ $tagversion -eq 0 ]] ; then
  tagversion=$(uuidgen)
fi
echo "building sasumaki/htmler:$tagversion"
docker build -t sasumaki/htmler:$tagversion ../htmler
docker push sasumaki/htmler:$tagversion

sed -i '' -e "s/sasumaki\/htmler:.*/sasumaki\/htmler:$tagversion/g" ../manifests/deployment.yaml

echo "sasumaki/htmler:$tagversion done"

echo "building sasumaki/backend:$tagversion"
docker build -t sasumaki/backend:$tagversion ../backend
docker push sasumaki/backend:$tagversion

sed -i '' -e "s/sasumaki\/backend:.*/sasumaki\/backend:$tagversion/g" ../manifests/deployment.yaml

echo "sasumaki/backend:$tagversion done"
kubectl apply -f ../manifests


echo "DEPLOYMENT DONE"
