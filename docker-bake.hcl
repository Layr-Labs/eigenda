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
}

target "disperser" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "apiserver"
}

target "encoder" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "encoder"
}

target "retriever" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "retriever"
}

target "churner" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "churner"
}

target "traffic-generator" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./trafficgenerator.Dockerfile"
  target     = "traffic-generator"
}

target "dataapi" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "dataapi"
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
}

target "nodeplugin" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "nodeplugin"
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

