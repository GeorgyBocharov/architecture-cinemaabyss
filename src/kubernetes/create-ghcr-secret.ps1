$NAMESPACE = "cinemaabyss"
$DOCKER_CONFIG_PATH = "..\..\.docker\config.json"
$SECRET_NAME = "dockerconfigjson"


kubectl create secret generic $SECRET_NAME `
    --namespace $NAMESPACE `
    --from-file=.dockerconfigjson=$DOCKER_CONFIG_PATH `
    --type=kubernetes.io/dockerconfigjson

Write-Host "Secret $SECRET_NAME created in namespace: $NAMESPACE"