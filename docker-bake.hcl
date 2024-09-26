# VARIABLES
variable "REGISTRY" {
  default = "ghcr.io"
}

variable "REPO" {
  default = "layr-labs/eigenda"
}

variable "BUILD_TAG" {
  default = "latest"
}

variable "SEMVER" {
  default = "v0.0.0"
}

variable "GITCOMMIT" {
  default = "dev"
}

variable "GITDATE" {
  default = "0"
}

# GROUPS

group "default" {
  targets = ["all"]
}

group "all" {
  targets = ["node-group", "disperser-group", "retriever", "churner"]
}

group "node-group" {
  targets = ["node", "nodeplugin"]
}

group "disperser-group" {
  targets = ["batcher", "disperser", "encoder"]
}

group "node-group-release" {
  targets = ["node-release", "nodeplugin-release"]
}

# DOCKER METADATA TARGET
# See https://github.com/docker/metadata-action?tab=readme-ov-file#bake-definition

target "docker-metadata-action" {}

# DISPERSER TARGETS

target "batcher" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "batcher"
  tags       = ["${REGISTRY}/${REPO}/batcher:${BUILD_TAG}"]
}

target "disperser" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "apiserver"
  tags       = ["${REGISTRY}/${REPO}/apiserver:${BUILD_TAG}"]
}

target "encoder" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "encoder"
  tags       = ["${REGISTRY}/${REPO}/encoder:${BUILD_TAG}"]
}

target "retriever" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "retriever"
  tags       = ["${REGISTRY}/${REPO}/retriever:${BUILD_TAG}"]
}

target "churner" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "churner"
  tags       = ["${REGISTRY}/${REPO}/churner:${BUILD_TAG}"]
}

target "traffic-generator" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./trafficgenerator.Dockerfile"
  target     = "trafficgenerator"
  tags       = ["${REGISTRY}/${REPO}/trafficgenerator:${BUILD_TAG}"]
}

target "dataapi" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "dataapi"
  tags       = ["${REGISTRY}/${REPO}/dataapi:${BUILD_TAG}"]
}

# NODE TARGETS

target "node" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "node"
  args = {
    SEMVER    = "${SEMVER}"
    GITCOMMIT = "${GITCOMMIT}"
    GITDATE   = "${GITDATE}"
  }
  tags = ["${REGISTRY}/${REPO}/node:${BUILD_TAG}"]
}

target "nodeplugin" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "nodeplugin"
  tags       = ["${REGISTRY}/${REPO}/nodeplugin:${BUILD_TAG}"]
}

# RELEASE TARGETS

target "_release" {
  platforms = ["linux/amd64", "linux/arm64"]
}

target "node-release" {
  inherits = ["node", "_release"]
  tags     = ["${REGISTRY}/${REPO}/opr-node:${BUILD_TAG}"]
}

target "nodeplugin-release" {
  inherits = ["nodeplugin", "_release"]
  tags     = ["${REGISTRY}/${REPO}/opr-nodeplugin:${BUILD_TAG}"]
}

