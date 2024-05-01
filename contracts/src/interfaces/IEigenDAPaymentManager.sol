// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";
import {IPaymentCoordinator} from "eigenlayer-core/contracts/interfaces/IPaymentCoordinator.sol";

/**
 * @title Interface for a contract that manages payments for EigenDA bandwidth reservations
 * @author Layr Labs, Inc.
 */
interface IEigenDAPaymentManager is ISignatureUtils{

    /**
     * @notice Data structure for bandwidth reservation payment information
     */
    struct ReservationPayment {
        IERC20[] tokens; //addresses of the tokens to be paid
        uint256[] amounts; //amounts of the tokens to be paid
        address payer; //address of the payer
        uint32 duration; //duration of the reservation
        uint16 bandwith; //bandwidth of the reservation
    }

    /// @notice emitted when a reservation payment is made
    event ReservationPaymentMade(ReservationPayment reservationPayment);

    /**
     * @notice Makes a reservation payment for EigenDA bandwith
     * @param reservationPayment The reservation payment details
     * @param paymentApproverSignature The payment approver's signature over the reservation payment details
     * @dev assumes that the payer has made appropiate approvals for the tokens to be transferred
     * @dev only callable by the payment approver or the reservation payer
     */
    function makeReservationPayment(
        ReservationPayment memory reservationPayment,
        SignatureWithSaltAndExpiry memory paymentApproverSignature
    ) external payable;

    /**
     * @notice Makes a range payment to EigenDA operators
     * @param rangePayments The range payment details
     * @dev only callable by the payment processor
     */
    function payForRange(
        IPaymentCoordinator.RangePayment[] calldata rangePayments
    ) external;

    /**
     * @notice Distributes funds to a recipient
     * @param to The recipient of the funds
     * @param tokens The tokens to distribute
     * @param amounts The amount of each to distribute
     * @dev only callable by the payment distributer
     */
    function distribute(
        address to,
        IERC20[] memory tokens,
        uint256[] memory amounts
    ) external;

    /**
     * @notice Getter for the the paymentApprover signature hash calculation 
     * @param reservationPayment The reservation payment details
     * @param salt The salt to use for the paymentApprover's signature
     * @param expiry The desired expiry time of the paymentApprover's signature
     */
    function calculatePaymentApprovalDigestHash(
        ReservationPayment memory reservationPayment,
        bytes32 salt,
        uint256 expiry
    ) external view returns (bytes32);
    
}