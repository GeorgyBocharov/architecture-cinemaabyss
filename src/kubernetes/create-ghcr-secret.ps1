$NAMESPACE = "cinemaabyss"
$DOCKER_CONFIG_PATH = "..\..\.docker\config.json"
$SECRET_NAME = "dockerconfigjson"


kubectl create secret generic ghcr-pull-secret `
    --namespace $NAMESPACE `
    --from-file=.dockerconfigjson=$DOCKER_CONFIG_PATH `
    --type=kubernetes.io/dockerconfigjson

Write-Host "Секрет $SECRET_NAME создан в namespace: $NAMESPACE"