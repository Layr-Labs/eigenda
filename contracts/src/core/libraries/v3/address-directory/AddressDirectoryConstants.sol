// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library AddressDirectoryConstants {
    // CORE
    bytes32 internal constant DISPERSER_REGISTRY_KEY = keccak256(abi.encodePacked(DISPERSER_REGISTRY_NAME));
    bytes32 internal constant RELAY_REGISTRY_KEY = keccak256(abi.encodePacked(RELAY_REGISTRY_NAME));
    bytes32 internal constant SERVICE_MANAGER_KEY = keccak256(abi.encodePacked(SERVICE_MANAGER_NAME));
    bytes32 internal constant THRESHOLD_REGISTRY_KEY = keccak256(abi.encodePacked(THRESHOLD_REGISTRY_NAME));
    bytes32 internal constant PAYMENT_VAULT_KEY = keccak256(abi.encodePacked(PAYMENT_VAULT_NAME));

    // MIDDLEWARE

    bytes32 internal constant REGISTRY_COORDINATOR_KEY = keccak256(abi.encodePacked(REGISTRY_COORDINATOR_NAME));
    bytes32 internal constant STAKE_REGISTRY_KEY = keccak256(abi.encodePacked(STAKE_REGISTRY_NAME));
    bytes32 internal constant INDEX_REGISTRY_KEY = keccak256(abi.encodePacked(INDEX_REGISTRY_NAME));
    bytes32 internal constant SOCKET_REGISTRY_KEY = keccak256(abi.encodePacked(SOCKET_REGISTRY_NAME));
    bytes32 internal constant PAUSER_REGISTRY_KEY = keccak256(abi.encodePacked(PAUSER_REGISTRY_NAME));
    bytes32 internal constant BLS_APK_REGISTRY_KEY = keccak256(abi.encodePacked(BLS_APK_REGISTRY_NAME));
    bytes32 internal constant EJECTION_MANAGER_KEY = keccak256(abi.encodePacked(EJECTION_MANAGER_NAME));

    /// PERIPHERY

    bytes32 internal constant OPERATOR_STATE_RETRIEVER_KEY = keccak256(abi.encodePacked(OPERATOR_STATE_RETRIEVER_NAME));
    bytes32 internal constant CERT_VERIFIER_KEY = keccak256(abi.encodePacked(CERT_VERIFIER_NAME));
    bytes32 internal constant CERT_VERIFIER_ROUTER_KEY = keccak256(abi.encodePacked(CERT_VERIFIER_ROUTER_NAME));

    /// LEGACY

    bytes32 internal constant CERT_VERIFIER_V1_KEY = keccak256(abi.encodePacked(CERT_VERIFIER_V1_NAME));
    bytes32 internal constant CERT_VERIFIER_V2_KEY = keccak256(abi.encodePacked(CERT_VERIFIER_V2_NAME));

    /// NAMES

    string internal constant DISPERSER_REGISTRY_NAME = "DISPERSER_REGISTRY";
    string internal constant RELAY_REGISTRY_NAME = "RELAY_REGISTRY";
    string internal constant SERVICE_MANAGER_NAME = "SERVICE_MANAGER";
    string internal constant THRESHOLD_REGISTRY_NAME = "THRESHOLD_REGISTRY";
    string internal constant PAYMENT_VAULT_NAME = "PAYMENT_VAULT";

    string internal constant REGISTRY_COORDINATOR_NAME = "REGISTRY_COORDINATOR";
    string internal constant STAKE_REGISTRY_NAME = "STAKE_REGISTRY";
    string internal constant INDEX_REGISTRY_NAME = "INDEX_REGISTRY";
    string internal constant SOCKET_REGISTRY_NAME = "SOCKET_REGISTRY";
    string internal constant PAUSER_REGISTRY_NAME = "PAUSER_REGISTRY";
    string internal constant BLS_APK_REGISTRY_NAME = "BLS_APK_REGISTRY";
    string internal constant EJECTION_MANAGER_NAME = "EJECTION_MANAGER";

    string internal constant OPERATOR_STATE_RETRIEVER_NAME = "OPERATOR_STATE_RETRIEVER";

    string internal constant CERT_VERIFIER_NAME = "CERT_VERIFIER";
    string internal constant CERT_VERIFIER_ROUTER_NAME = "CERT_VERIFIER_ROUTER";

    string internal constant CERT_VERIFIER_V1_NAME = "CERT_VERIFIER_V1";
    string internal constant CERT_VERIFIER_V2_NAME = "CERT_VERIFIER_V2";
}
