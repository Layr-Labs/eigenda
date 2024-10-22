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

# Release targets will fail if GIT_SHA env is not exported. See Makefile:docker-release-build
variable "GIT_SHA" {
  default = "$GIT_SHA NOT DEFINED"
}

# Release targets will fail if GIT_SHORT_SHA env is not exported. See Makefile:docker-release-build
variable "GIT_SHORT_SHA" {
  default = "$GIT_SHORT_SHA NOT DEFINED"
}

variable "GITDATE" {
  default = "0"
}

# GROUPS

group "default" {
  targets = ["all"]
}

group "all" {
  targets = ["node-group", "batcher", "disperser", "encoder", "retriever", "churner", "dataapi", "traffic-generator"]
}

group "node-group" {
  targets = ["node", "nodeplugin"]
}

# Github public releases
group "node-group-release" {
  targets = ["node-release", "nodeplugin-release"]
}

# Github CI builds
group "ci-release" {
  targets = ["node-group", "batcher", "disperser", "encoder", "retriever", "churner", "dataapi"]
}

# Internal devops builds
group "internal-release" {
  targets = ["node-internal", "batcher-internal", "disperser-internal", "encoder-internal", "retriever-internal", "churner-internal", "dataapi-internal", "traffic-generator-internal"]
}


# DISPERSER TARGETS

target "batcher" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "batcher"
  tags       = ["${REGISTRY}/${REPO}/batcher:${BUILD_TAG}"]
}

target "batcher-internal" {
  inherits = ["batcher"]
  tags       = ["${REGISTRY}/eigenda-batcher:${BUILD_TAG}",
                "${REGISTRY}/eigenda-batcher:${GIT_SHA}",
                "${REGISTRY}/eigenda-batcher:sha-${GIT_SHORT_SHA}",
               ]
}

target "disperser" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "apiserver"
  tags       = ["${REGISTRY}/${REPO}/apiserver:${BUILD_TAG}"]
}

target "disperser-internal" {
  inherits = ["disperser"]
  tags       = ["${REGISTRY}/eigenda-disperser:${BUILD_TAG}",
                "${REGISTRY}/eigenda-disperser:${GIT_SHA}",
                "${REGISTRY}/eigenda-disperser:sha-${GIT_SHORT_SHA}",
               ]
}

target "encoder" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "encoder"
  tags       = ["${REGISTRY}/${REPO}/encoder:${BUILD_TAG}"]
}

target "encoder-internal" {
  inherits = ["encoder"]
  tags       = ["${REGISTRY}/eigenda-encoder:${BUILD_TAG}",
                "${REGISTRY}/eigenda-encoder:${GIT_SHA}",
                "${REGISTRY}/eigenda-encoder:sha-${GIT_SHORT_SHA}",
               ]
}

target "retriever" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "retriever"
  tags       = ["${REGISTRY}/${REPO}/retriever:${BUILD_TAG}"]
}

target "retriever-internal" {
  inherits = ["retriever"]
  tags       = ["${REGISTRY}/eigenda-retriever:${BUILD_TAG}",
                "${REGISTRY}/eigenda-retriever:${GIT_SHA}",
                "${REGISTRY}/eigenda-retriever:sha-${GIT_SHORT_SHA}",
               ]
}

target "churner" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "churner"
  tags       = ["${REGISTRY}/${REPO}/churner:${BUILD_TAG}"]
}

target "churner-internal" {
  inherits = ["churner"]
  tags       = ["${REGISTRY}/eigenda-churner:${BUILD_TAG}",
                "${REGISTRY}/eigenda-churner:${GIT_SHA}",
                "${REGISTRY}/eigenda-churner:sha-${GIT_SHORT_SHA}",
               ]
}

target "traffic-generator" {
  context    = "."
  dockerfile = "./trafficgenerator.Dockerfile"
  target     = "generator"
  tags       = ["${REGISTRY}/${REPO}/traffic-generator:${BUILD_TAG}"]
}

target "traffic-generator-internal" {
  inherits = ["traffic-generator"]
  tags       = ["${REGISTRY}/eigenda-traffic-generator:${BUILD_TAG}",
                "${REGISTRY}/eigenda-traffic-generator:${GIT_SHA}",
                "${REGISTRY}/eigenda-traffic-generator:sha-${GIT_SHORT_SHA}",
               ]
}

target "traffic-generator2" {
  context    = "."
  dockerfile = "./trafficgenerator2.Dockerfile"
  target     = "generator2"
  tags       = ["${REGISTRY}/${REPO}/traffic-generator2:${BUILD_TAG}"]
}

target "traffic-generator2-internal" {
  inherits = ["traffic-generator2"]
  tags       = ["${REGISTRY}/eigenda-traffic-generator2:${BUILD_TAG}",
    "${REGISTRY}/eigenda-traffic-generator2:${GIT_SHA}",
    "${REGISTRY}/eigenda-traffic-generator2:sha-${GIT_SHORT_SHA}",
  ]
}

target "dataapi" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "dataapi"
  tags       = ["${REGISTRY}/${REPO}/dataapi:${BUILD_TAG}"]
}

target "dataapi-internal" {
  inherits = ["dataapi"]
  tags       = ["${REGISTRY}/eigenda-dataapi:${BUILD_TAG}",
                "${REGISTRY}/eigenda-dataapi:${GIT_SHA}",
                "${REGISTRY}/eigenda-dataapi:sha-${GIT_SHORT_SHA}",
               ]
}

# NODE TARGETS

target "node" {
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

target "node-internal" {
  inherits = ["node"]
  tags       = ["${REGISTRY}/eigenda-node:${BUILD_TAG}",
                "${REGISTRY}/eigenda-node:${GIT_SHA}",
                "${REGISTRY}/eigenda-node:sha-${GIT_SHORT_SHA}",
               ]
}

target "nodeplugin" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "nodeplugin"
  tags       = ["${REGISTRY}/${REPO}/nodeplugin:${BUILD_TAG}"]
}

# PUBLIC RELEASE TARGETS

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
