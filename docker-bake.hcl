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

variable "GIT_SHA" {
  default = ""
}

variable "GIT_SHORT_SHA" {
  default = ""
}

variable "GITDATE" {
  default = "0"
}

# GROUPS

group "default" {
  targets = ["all"]
}

group "all" {
  targets = ["node-group", "batcher", "disperser", "encoder", "retriever", "churner", "dataapi"]
}

group "node-group" {
  targets = ["node", "nodeplugin"]
}

group "node-group-release" {
  targets = ["node-release", "nodeplugin-release"]
}

# CI builds
group "ci-release" {
  targets = ["node-group", "batcher", "disperser", "encoder", "retriever", "churner", "dataapi"]
}

# internal devops builds
group "internal-release" {
  targets = ["batcher-release", "disperser-release", "encoder-release", "retriever-release", "churner-release", "dataapi-release"]
}


# DOCKER METADATA TARGET
# See https://github.com/docker/metadata-action?tab=readme-ov-file#bake-definition

target "docker-metadata-action" {}

# DISPERSER TARGETS

target "batcher" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "batcher"
  tags       = ["${REGISTRY}/${REPO}/batcher:${BUILD_TAG}"]
}

target "batcher-release" {
  inherits = ["batcher"]
  tags       = ["${REGISTRY}/${REPO}/eigenda-batcher:${BUILD_TAG}",
                "${REGISTRY}/${REPO}/eigenda-batcher:${GIT_SHA}",
                "${REGISTRY}/${REPO}/eigenda-batcher:sha-${GIT_SHORT_SHA}",
               ]
}

target "disperser" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "apiserver"
  tags       = ["${REGISTRY}/${REPO}/apiserver:${BUILD_TAG}"]
}

target "disperser-release" {
  inherits = ["disperser"]
  tags       = ["${REGISTRY}/${REPO}/eigenda-disperser:${BUILD_TAG}",
                "${REGISTRY}/${REPO}/eigenda-disperser:${GIT_SHA}",
                "${REGISTRY}/${REPO}/eigenda-disperser:sha-${GIT_SHORT_SHA}",
               ]
}

target "encoder" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "encoder"
  tags       = ["${REGISTRY}/${REPO}/encoder:${BUILD_TAG}"]
}

target "encoder-release" {
  inherits = ["encoder"]
  tags       = ["${REGISTRY}/${REPO}/eigenda-encoder:${BUILD_TAG}",
                "${REGISTRY}/${REPO}/eigenda-encoder:${GIT_SHA}",
                "${REGISTRY}/${REPO}/eigenda-encoder:sha-${GIT_SHORT_SHA}",
               ]
}

target "retriever" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "retriever"
  tags       = ["${REGISTRY}/${REPO}/retriever:${BUILD_TAG}"]
}

target "retriever-release" {
  inherits = ["retriever"]
  tags       = ["${REGISTRY}/${REPO}/eigenda-retriever:${BUILD_TAG}",
                "${REGISTRY}/${REPO}/eigenda-retriever:${GIT_SHA}",
                "${REGISTRY}/${REPO}/eigenda-retriever:sha-${GIT_SHORT_SHA}",
               ]
}

target "churner" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "churner"
  tags       = ["${REGISTRY}/${REPO}/churner:${BUILD_TAG}"]
}

target "churner-release" {
  inherits = ["churner"]
  tags       = ["${REGISTRY}/${REPO}/eigenda-churner:${BUILD_TAG}",
                "${REGISTRY}/${REPO}/eigenda-churner:${GIT_SHA}",
                "${REGISTRY}/${REPO}/eigenda-churner:sha-${GIT_SHORT_SHA}",
               ]
}

target "traffic-generator" {
  context    = "."
  dockerfile = "./trafficgenerator.Dockerfile"
  target     = "trafficgenerator"
  tags       = ["${REGISTRY}/${REPO}/traffic-generator:${BUILD_TAG}"]
}

target "traffic-generator-release" {
  inherits = ["traffic-generator"]
  tags       = ["${REGISTRY}/${REPO}/eigenda-traffic-generator:${BUILD_TAG}",
                "${REGISTRY}/${REPO}/eigenda-traffic-generator:${GIT_SHA}",
                "${REGISTRY}/${REPO}/eigenda-traffic-generator:sha-${GIT_SHORT_SHA}",
               ]
}

target "dataapi" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "dataapi"
  tags       = ["${REGISTRY}/${REPO}/dataapi:${BUILD_TAG}"]
}

target "dataapi-release" {
  inherits = ["dataapi"]
  tags       = ["${REGISTRY}/${REPO}/eigenda-dataapi:${BUILD_TAG}",
                "${REGISTRY}/${REPO}/eigenda-dataapi:${GIT_SHA}",
                "${REGISTRY}/${REPO}/eigenda-dataapi:sha-${GIT_SHORT_SHA}",
               ]
}

# NODE TARGETS

target "node" {
  inherits = ["docker-metadata-action"]
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "node"
  args = {
    SEMVER    = "${SEMVER}"
    GITCOMMIT = "${GIT_SHORT_SHA}"
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
