// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {IEigenDAPaymentManager} from "../interfaces/IEigenDAPaymentManager.sol";

/**
 * @title Storage variables for the `EigenDAPaymentManager` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract EigenDAPaymentManagerStorage is IEigenDAPaymentManager, OwnableUpgradeable {

    /// @notice The EIP-712 typehash for the `PaymentApproval` struct used by the contract
    bytes32 public constant PAYMENT_APPROVAL_TYPEHASH =
        keccak256("PaymentApproval(ReservationPayment reservationPayment,bytes32 salt,uint256 expiry)ReservationPayment(address[] tokens,uint256[] amounts,address payer,uint32 duration,uint16 bandwith)");


    /// @notice the basis point denominator for bandwidth as Mb/s
    uint16 public constant BANDWITH_BIPS_DENOMINATOR = 10000;

    /// @notice the ServiceManager for this AVS, which forwards calls onto EigenLayer's core contracts
    IServiceManager public immutable serviceManager;

    constructor(IServiceManager _serviceManager) {
        serviceManager = _serviceManager;
    }

    /// @notice the address of the entity that can approve payments via signature
    address paymentApprover;
    /// @notice the address of the entity that can make range payments
    address paymentProcessor;
    /// @notice the address of the entity that can distribute funds outside of range payments
    address paymentDistributer;

    /// @notice whether the salt has been used for a payment approval
    mapping(bytes32 => bool) public isPaymentApproverSaltUsed;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[46] private __GAP;
}