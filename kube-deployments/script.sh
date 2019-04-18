if type -p kubectl; then
  kubectl apply -f cache-service-service.yaml,mongo-service.yaml,rabbitmq-service.yaml
  kubectl apply -f mongo-deployment.yaml,rabbitmq-deployment.yaml
  echo "sleeping 5 mins before starting cache-service"
  sleep 5s
  kubectl apply -f cache-service-pod.yaml
else
  echo "kubectl command not found in path variable"
fi
