# VARIABLES
variable "REGISTRY" {
  default = "ghcr.io"
}

variable "REPO" {
  default = "layr-labs/eigenda"
}

# We use the `dev` tag for local development builds.
# CI builds will overwrite this with the `master` or `v*` tag.
variable "BUILD_TAG" {
  default = "dev"
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
  targets = [
    "node-group",
    "batcher",
    "disperser",
    "encoder",
    "retriever",
    "churner",
    "dataapi",
    "traffic-generator",
    "traffic-generator-v2",
    "controller",
    "relay",
    "blobapi",
    "proxy",
  ]
}

group "node-group" {
  targets = ["node", "nodeplugin"]
}

# Internal devops builds. These targets are used by the eigenda-devops CI pipeline.
# TODO: refactor the ECR repo to make the `${REGISTRY}/${REPO}` tags such that we can
# get rid of all of these internal targets.
group "internal-release" {
  targets = [
    "node-internal",
    "batcher-internal",
    "disperser-internal",
    "encoder-internal",
    "retriever-internal",
    "churner-internal",
    "dataapi-internal",
    "traffic-generator-internal",
    "traffic-generator-v2-internal",
    "controller-internal",
    "relay-internal",
    "blobapi-internal",
    "proxy-internal",
  ]
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
  tags     = [
    "${REGISTRY}/eigenda-batcher:${BUILD_TAG}",
    "${REGISTRY}/eigenda-batcher:${GIT_SHA}",
    "${REGISTRY}/eigenda-batcher:sha-${GIT_SHORT_SHA}"
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
  tags     = [
    "${REGISTRY}/eigenda-disperser:${BUILD_TAG}",
    "${REGISTRY}/eigenda-disperser:${GIT_SHA}",
    "${REGISTRY}/eigenda-disperser:sha-${GIT_SHORT_SHA}"
  ]
}

target "encoder" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "encoder"
  tags       = ["${REGISTRY}/${REPO}/encoder:${BUILD_TAG}"]
}

target "encoder-icicle" {
  context    = "."
  dockerfile = "./disperser/cmd/encoder/icicle.Dockerfile"
  tags       = ["${REGISTRY}/${REPO}/encoder-icicle:${BUILD_TAG}"]
}

target "encoder-internal" {
  inherits = ["encoder"]
  tags     = [
    "${REGISTRY}/eigenda-encoder:${BUILD_TAG}",
    "${REGISTRY}/eigenda-encoder:${GIT_SHA}",
    "${REGISTRY}/eigenda-encoder:sha-${GIT_SHORT_SHA}"
  ]
}

target "encoder-icicle-internal" {
  inherits = ["encoder-icicle"]
  tags     = [
    "${REGISTRY}/eigenda-encoder-icicle:${BUILD_TAG}",
    "${REGISTRY}/eigenda-encoder-icicle:${GIT_SHA}",
    "${REGISTRY}/eigenda-encoder-icicle:sha-${GIT_SHORT_SHA}"
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
  tags     = [
    "${REGISTRY}/eigenda-retriever:${BUILD_TAG}",
    "${REGISTRY}/eigenda-retriever:${GIT_SHA}",
    "${REGISTRY}/eigenda-retriever:sha-${GIT_SHORT_SHA}"
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
  tags     = [
    "${REGISTRY}/eigenda-churner:${BUILD_TAG}",
    "${REGISTRY}/eigenda-churner:${GIT_SHA}",
    "${REGISTRY}/eigenda-churner:sha-${GIT_SHORT_SHA}"
  ]
}

target "traffic-generator" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "generator"
  tags       = ["${REGISTRY}/${REPO}/traffic-generator:${BUILD_TAG}"]
}

target "traffic-generator-internal" {
  inherits = ["traffic-generator"]
  tags     = [
    "${REGISTRY}/eigenda-traffic-generator:${BUILD_TAG}",
    "${REGISTRY}/eigenda-traffic-generator:${GIT_SHA}",
    "${REGISTRY}/eigenda-traffic-generator:sha-${GIT_SHORT_SHA}"
  ]
}

target "traffic-generator-v2" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "generator2"
  tags       = ["${REGISTRY}/${REPO}/traffic-generator-v2:${BUILD_TAG}"]
}

target "traffic-generator-v2-internal" {
  inherits = ["traffic-generator-v2"]
  tags     = [
    "${REGISTRY}/eigenda-traffic-generator-v2:${BUILD_TAG}",
    "${REGISTRY}/eigenda-traffic-generator-v2:${GIT_SHA}",
    "${REGISTRY}/eigenda-traffic-generator-v2:sha-${GIT_SHORT_SHA}"
  ]
}

target "relay" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "relay"
  tags       = ["${REGISTRY}/${REPO}/relay:${BUILD_TAG}"]
}

target "relay-internal" {
  inherits = ["relay"]
  tags     = [
    "${REGISTRY}/eigenda-relay:${BUILD_TAG}",
    "${REGISTRY}/eigenda-relay:${GIT_SHA}",
    "${REGISTRY}/eigenda-relay:sha-${GIT_SHORT_SHA}"
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
  tags     = [
    "${REGISTRY}/eigenda-dataapi:${BUILD_TAG}",
    "${REGISTRY}/eigenda-dataapi:${GIT_SHA}",
    "${REGISTRY}/eigenda-dataapi:sha-${GIT_SHORT_SHA}"
  ]
}

target "controller" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "controller"
  tags       = ["${REGISTRY}/${REPO}/controller:${BUILD_TAG}"]
}

target "controller-internal" {
  inherits = ["controller"]
  tags     = [
    "${REGISTRY}/eigenda-controller:${BUILD_TAG}",
    "${REGISTRY}/eigenda-controller:${GIT_SHA}",
    "${REGISTRY}/eigenda-controller:sha-${GIT_SHORT_SHA}"
  ]
}

target "blobapi" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "blobapi"
  tags       = ["${REGISTRY}/${REPO}/blobapi:${BUILD_TAG}"]
}

target "blobapi-internal" {
  inherits = ["blobapi"]
  tags     = [
    "${REGISTRY}/eigenda-blobapi:${BUILD_TAG}",
    "${REGISTRY}/eigenda-blobapi:${GIT_SHA}",
    "${REGISTRY}/eigenda-blobapi:sha-${GIT_SHORT_SHA}"
  ]
}

target "proxy" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "proxy"
  # We push to layr-labs/ directly instead of layr-labs/eigenda/ for historical reasons,
  # since proxy was previously in its own repo: https://github.com/Layr-Labs/eigenda-proxy
  tags       = ["${REGISTRY}/layr-labs/eigenda-proxy:${BUILD_TAG}"]
}

target "proxy-internal" {
  inherits = ["proxy"]
  tags     = [
    "${REGISTRY}/eigenda-proxy:${BUILD_TAG}",
    "${REGISTRY}/eigenda-proxy:${GIT_SHA}",
    "${REGISTRY}/eigenda-proxy:sha-${GIT_SHORT_SHA}"
  ]
}

# NODE TARGETS
target "node" {
  context    = "."
  dockerfile = "./Dockerfile"
  target     = "node"
  args       = {
    SEMVER    = "${SEMVER}"
    GITCOMMIT = "${GIT_SHORT_SHA}"
    GITDATE   = "${GITDATE}"
  }
  tags = ["${REGISTRY}/${REPO}/node:${BUILD_TAG}"]
}

target "node-internal" {
  inherits = ["node"]
  tags     = [
    "${REGISTRY}/eigenda-node:${BUILD_TAG}",
    "${REGISTRY}/eigenda-node:${GIT_SHA}",
    "${REGISTRY}/eigenda-node:sha-${GIT_SHORT_SHA}"
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

group "node-group-release" {
  targets = ["node-release", "nodeplugin-release"]
}

target "node-release" {
  inherits = ["node", "_release"]
  # We overwrite the tag with a opr- prefix for public releases.
  tags     = ["${REGISTRY}/${REPO}/opr-node:${BUILD_TAG}"]
}

target "nodeplugin-release" {
  inherits = ["nodeplugin", "_release"]
  # We overwrite the tag with a opr- prefix for public releases.
  tags     = ["${REGISTRY}/${REPO}/opr-nodeplugin:${BUILD_TAG}"]
}

target "proxy-release" {
  inherits = ["proxy", "_release"]
  # We push to layr-labs/ directly instead of layr-labs/eigenda/ for historical reasons,
  # since proxy was previously in its own repo: https://github.com/Layr-Labs/eigenda-proxy
  tags     = ["${REGISTRY}/eigenda-proxy:${BUILD_TAG}"]
}