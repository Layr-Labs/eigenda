// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library AddressDirectoryConstants {
    /// CORE

    string internal constant DISPERSER_REGISTRY_NAME = "DISPERSER_REGISTRY";
    string internal constant RELAY_REGISTRY_NAME = "RELAY_REGISTRY";
    string internal constant SERVICE_MANAGER_NAME = "SERVICE_MANAGER";
    string internal constant THRESHOLD_REGISTRY_NAME = "THRESHOLD_REGISTRY";
    string internal constant PAYMENT_VAULT_NAME = "PAYMENT_VAULT";

    /// MIDDLEWARE

    string internal constant REGISTRY_COORDINATOR_NAME = "REGISTRY_COORDINATOR";
    string internal constant STAKE_REGISTRY_NAME = "STAKE_REGISTRY";
    string internal constant INDEX_REGISTRY_NAME = "INDEX_REGISTRY";
    string internal constant SOCKET_REGISTRY_NAME = "SOCKET_REGISTRY";
    string internal constant PAUSER_REGISTRY_NAME = "PAUSER_REGISTRY";
    string internal constant BLS_APK_REGISTRY_NAME = "BLS_APK_REGISTRY";
    string internal constant EJECTION_MANAGER_NAME = "EJECTION_MANAGER";

    /// PERIPHERY

    string internal constant OPERATOR_STATE_RETRIEVER_NAME = "OPERATOR_STATE_RETRIEVER";

    string internal constant CERT_VERIFIER_NAME = "CERT_VERIFIER";
    string internal constant CERT_VERIFIER_ROUTER_NAME = "CERT_VERIFIER_ROUTER";

    /// LEGACY

    string internal constant CERT_VERIFIER_LEGACY_V1_NAME = "CERT_VERIFIER_LEGACY_V1";
    string internal constant CERT_VERIFIER_LEGACY_V2_NAME = "CERT_VERIFIER_LEGACY_V2";
    string internal constant RELAY_REGISTRY_LEGACY_NAME = "RELAY_REGISTRY_LEGACY";
}
