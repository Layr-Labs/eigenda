// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {PaymentVaultTypes} from "src/core/libraries/v3/PaymentVaultLib.sol";

interface IPaymentVault {
    function getOnDemandDeposit(uint64 quorumId, address account) external view returns (uint256);

    function getReservation(uint64 quorumId, address account)
        external
        view
        returns (PaymentVaultTypes.Reservation memory);

    function getQuorumProtocolConfig(uint64 quorumId)
        external
        view
        returns (PaymentVaultTypes.QuorumProtocolConfig memory);

    function getQuorumPaymentConfig(uint64 quorumId) external view returns (PaymentVaultTypes.QuorumConfig memory);

    function getQuorumReservedSymbols(uint64 quorumId, uint64 period) external view returns (uint64);
}
