# VARIABLES
variable "REGISTRY" {
  default = "ghcr.io/layr-labs/eigenda"
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

variable "BATCHER_PATH" {
  default =  "batcher"
}

variable "DISPERSER_PATH" {
  default =  "disperser"
}

variable "DATAAPI_PATH" {
  default =  "dataapi"
}

variable "ENCODER_PATH" {
  default =  "encoder"
}

variable "RETRIEVER_PATH" {
  default =  "retriever"
}

variable "CHURNER_PATH" {
  default =  "churner"
}

variable "NODE_PATH" {
  default =  "node"
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
  tags       = ["${REGISTRY}/${BATCHER_PATH}:${BUILD_TAG}"]
}

target "disperser" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "apiserver"
  tags       = ["${REGISTRY}/${DISPERSER_PATH}:${BUILD_TAG}"]
}

target "encoder" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "encoder"
  tags       = ["${REGISTRY}/${ENCODER_PATH}:${BUILD_TAG}"]
}

target "retriever" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "retriever"
  tags       = ["${REGISTRY}/${RETRIEVER_PATH}:${BUILD_TAG}"]
}

target "churner" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "churner"
  tags       = ["${REGISTRY}/${CHURNER_PATH}:${BUILD_TAG}"]
}

target "traffic-generator" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./traffic-generator.Dockerfile"
  target     = "traffic-generator"
  tags       = []
}

target "dataapi" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "dataapi"
  tags       = ["${REGISTRY}/${DATAAPI_PATH}:${BUILD_TAG}"]
}

# NODE TARGETS

target "node" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "node"
  tags       = ["${REGISTRY}/${NODE_PATH}:${BUILD_TAG}"]
  args = {
    SEMVER    = "${SEMVER}"
    GITCOMMIT = "${GITCOMMIT}"
    GITDATE   = "${GITDATE}"
  }
}

target "nodeplugin" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "nodeplugin"
  tags       = ["${REGISTRY}/nodeplugin:${BUILD_TAG}"]
}

# RELEASE TARGETS

target "_release" {
  platforms = ["linux/amd64", "linux/arm64"]
}

target "node-release" {
  inherits = ["node", "_release"]
  tags     = ["${REGISTRY}/opr-node:${BUILD_TAG}"]
}

target "nodeplugin-release" {
  inherits = ["nodeplugin", "_release"]
  tags     = ["${REGISTRY}/opr-nodeplugin:${BUILD_TAG}"]
}

