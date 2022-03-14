#!/bin/bash

# --------------------------------------------
# Options that must be configured by app owner
# --------------------------------------------
export APP_NAME="vmaas-go"  # name of app-sre "application" folder this component lives in
export COMPONENT_NAME="vmaas-go"  # name of app-sre "resourceTemplate" in deploy.yaml for this component
export IMAGE="quay.io/cloudservices/vmaas-go-app"
export DOCKERFILE="Dockerfile"
export COMPONENTS_W_RESOURCES="vmaas"

export IQE_PLUGINS="vmaas-go"
export IQE_MARKER_EXPRESSION="vmaas_go_smoke"
export IQE_FILTER_EXPRESSION=""
export IQE_CJI_TIMEOUT="30m"


# Install bonfire repo/initialize
CICD_URL=https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd
curl -s $CICD_URL/bootstrap.sh > .cicd_bootstrap.sh && source .cicd_bootstrap.sh

source $CICD_ROOT/build.sh

# Deploy to an ephemeral namespace for testing
source $CICD_ROOT/deploy_ephemeral_env.sh

# Run iqe somoke tests with ClowdJobInvocation
source $CICD_ROOT/cji_smoke_test.sh
