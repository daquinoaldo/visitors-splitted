# Kubernetes Operator notes

## Kickstart

### Requirements
- OpenShift (`oc`) and probably also `kubectl` and Kubernets for the APIs.
- Docker for the build (podman should work too, but I had some trouble).
- Golang (for the Go operator) - [DigitalOcean guide](https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-ubuntu-18-04)
- Ansible:
```bash
# Ansible (either one of this command)
sudo dnf install ansible
pip install --user ansible

# Ansible Runner (either one of this command)
sudo dnf install python-ansible-runner
pip install --user ansible-runner

# Ansible Runner HTTP plugin
pip install --user ansible-runner-http
```

### Variables
```bash
OPERATOR_NAME="appname-operator"
APP_NAME="AppnameApp"
APP_DOMAIN="example.com"
USERNAME="aldo.daquino"
REPO="github.com/$USERNAME/$OPERATOR_NAME"
```

### Go skeleton generation
```bash
operator-sdk new $OPERATOR_NAME --repo $REPO
operator-sdk add api --api-version=$APP_DOMAIN/v1 --kind=$APP_NAME
operator-sdk add controller --api-version=$APP_DOMAIN/v1 --kind=$APP_NAME
```

### Ansible skeleton generation
```
operator-sdk new $OPERATOR_NAME --type=ansible --api-version=$APP_DOMAIN/v1 --kind=$APP_NAME
```


## Run locally
```bash
# Deploy the CRD (Custom Resources Definition)
oc apply -f deploy/crds/*_crd.yaml

# Run the operator locally
operator-sdk run --local

# Deploy a Custom Resource
oc apply -f deploy/crds/*_cr.yaml

# Remove the CR
oc delete -f deploy/crds/*_cr.yaml
```


## Build
Note: if you are building an Ansible operator and you have changed the `role` path in `watches.yaml`, restore it back to `/opt/ansible/roles/appname`.

```bash
# Variables
OPERATOR_NAME=${PWD##*/}
PROJECT=$(oc project -q)
REGISTRY="default-route-openshift-image-registry.apps-crc.testing"
IMAGE_NAME="$REGISTRY/$PROJECT/$OPERATOR_NAME"
INTERNAL_REGISTRY="image-registry.openshift-image-registry.svc:5000"
INTERNAL_IMAGE="$INTERNAL_REGISTRY/$PROJECT/$OPERATOR_NAME"

# Login to the OpenShift internal registry
oc extract secret/router-ca --keys=tls.crt -n openshift-ingress-operator
sudo mkdir -p /etc/docker/certs.d/$REGISTRY/
sudo mv tls.crt /etc/docker/certs.d/$REGISTRY/
docker login -u kubeadmin -p $(oc whoami -t) $REGISTRY

# Build
operator-sdk build $IMAGE_NAME

# Push
docker push $IMAGE_NAME

# Update the manifest (be careful to the OS and the implementation)
# Go operator
sed -i "s|REPLACE_IMAGE|$INTERNAL_IMAGE|g" deploy/operator.yaml     # linux
sed -i "" "s|REPLACE_IMAGE|$INTERNAL_IMAGE|g" deploy/operator.yaml  # macOS

# Ansible operator
sed -i "s|{{ REPLACE_IMAGE }}|$INTERNAL_IMAGE|g" deploy/operator.yaml         # linux
sed -i "s|{{ pull_policy\|default('Always') }}|Always|g" deploy/operator.yaml  # linux
sed -i "" "s|{{ REPLACE_IMAGE }}|$INTERNAL_IMAGE|g" deploy/operator.yaml         # macOS
sed -i "" "s|{{ pull_policy\|default('Always') }}|Always|g" deploy/operator.yaml  # macOS


# Setup Service Account, Role and Role Binding, and finally the CRD
oc apply -f deploy/service_account.yaml
oc apply -f deploy/role.yaml
oc apply -f deploy/role_binding.yaml
oc apply -f deploy/crds/*_crd.yaml

# Deploy the operator inside a container (actually 2 containers in the Ansible implementation)
oc apply -f deploy/operator.yaml

# Deploy the CR
oc apply -f deploy/crds/*_cr.yaml
```
