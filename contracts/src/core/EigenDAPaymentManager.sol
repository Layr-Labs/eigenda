// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {EIP712} from "@openzeppelin/contracts/utils/cryptography/draft-EIP712.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {EIP1271SignatureUtils} from "eigenlayer-core/contracts/libraries/EIP1271SignatureUtils.sol";
import {IPaymentCoordinator} from "eigenlayer-core/contracts/interfaces/IPaymentCoordinator.sol";
import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {EigenDAPaymentManagerStorage} from "./EigenDAPaymentManagerStorage.sol";

/**
 * @title A `EigenDAPaymentManager` that takes payment for EigenDA bandwidth reservations and makes payments to EigenDA operators
 * @author Layr Labs, Inc.
 */
contract EigenDAPaymentManager is EigenDAPaymentManagerStorage, EIP712 {

    /// @notice checks if the caller is the payment approver
    modifier onlyPaymentApprover() {
        require(msg.sender == paymentApprover, "EigenDAPaymentManager.onlyPaymentProcessor: caller is not payment approver");
        _;
    }

    /// @notice checks if the caller is the payment processor
    modifier onlyPaymentProcessor() {
        require(msg.sender == paymentProcessor, "EigenDAPaymentManager.onlyPaymentProcessor: caller is not payment processor");
        _;
    }

    /// @notice checks if the caller is the payment distributer
    modifier onlyPaymentDistributer() {
        require(msg.sender == paymentDistributer, "EigenDAPaymentManager.onlyPaymentDistributer: caller is not payment distributer");
        _;
    }

    constructor(IServiceManager _serviceManager) 
        EigenDAPaymentManagerStorage(_serviceManager)
        EIP712("EigenDAPaymentManager", "v0.0.1") 
    {
        _disableInitializers();
    }

    /**
     * @param _initialOwner will hold the owner role
     * @param _paymentApprover will hold the payment approver role
     * @param _paymentProcessor will hold the payment processor role
     * @param _paymentDistributer will hold the payment distributer role
     */
    function initialize(
        address _initialOwner,
        address _paymentApprover,
        address _paymentProcessor,
        address _paymentDistributer
    ) external initializer {
        _transferOwnership(_initialOwner);
        _setPaymentApprover(_paymentApprover);
        _setPaymentProcessor(_paymentProcessor);
        _setPaymentDistributer(_paymentDistributer);
    }

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
    ) external payable {
        require(
            msg.sender == paymentApprover || 
            msg.sender == reservationPayment.payer, 
            "EigenDAPaymentManager.makeReservationPayment: caller is not payment approver or reservation payer"
        );
        require(reservationPayment.tokens.length == reservationPayment.amounts.length, "EigenDAPaymentManager.makeReservationPayment: length mismatch");

        _verifyPaymentApproverSignature(reservationPayment, paymentApproverSignature);

        for (uint256 i = 0; i < reservationPayment.tokens.length; ++i) {
            reservationPayment.tokens[i].transferFrom(reservationPayment.payer, address(this), reservationPayment.amounts[i]);
        }

        emit ReservationPaymentMade(reservationPayment);
    }

    /**
     * @notice Makes a range payment to EigenDA operators
     * @param rangePayments The range payment details
     * @dev only callable by the payment processor
     */
    function payForRange(
        IPaymentCoordinator.RangePayment[] calldata rangePayments
    ) external onlyPaymentProcessor {
        for (uint256 i = 0; i < rangePayments.length; ++i) {
            rangePayments[i].token.approve(address(serviceManager), rangePayments[i].amount);
        }
        serviceManager.payForRange(rangePayments);
    }

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
    ) external onlyOwner {
        require(tokens.length == amounts.length, "EigenDAPaymentManager.distribute: length mismatch");
        for (uint256 i = 0; i < tokens.length; ++i) {
            tokens[i].transfer(to, amounts[i]);
        }
    }

    /**
     * @notice Sets the payment approver which can approve payments via signature
     * @param _paymentApprover the address of the payment approver
     * @dev only callable by the owner
     */
    function setPaymentApprover(address _paymentApprover) external onlyOwner {
        _setPaymentApprover(_paymentApprover);
    }

    /**
     * @notice Sets the payment processor which can make range payments
     * @param _paymentProcessor the address of the payment processor
     * @dev only callable by the owner
     */
    function setPaymentProcessor(address _paymentProcessor) external onlyOwner {
        _setPaymentProcessor(_paymentProcessor);
    }

    /**
     * @notice Sets the payment distributer which can distribute funds outside of range payments
     * @param _paymentDistributer the address of the payment distributer
     * @dev only callable by the payment distributer
     */
    function setPaymentDistributer(address _paymentDistributer) external onlyPaymentDistributer() {
        _setPaymentDistributer(_paymentDistributer);
    }

    /// @notice internal function to set the payment approver
    function _setPaymentApprover(address _paymentApprover) internal {
        paymentApprover = _paymentApprover;
    }

    /// @notice internal function to set the payment processor
    function _setPaymentProcessor(address _paymentProcessor) internal {
        paymentProcessor = _paymentProcessor;
    }

    /// @notice internal function to set the payment distributer
    function _setPaymentDistributer(address _paymentDistributer) internal {
        paymentDistributer = _paymentDistributer;
    }

    /// @notice verifies paymentApprover's signature over reservation payment details
    function _verifyPaymentApproverSignature(
        ReservationPayment memory reservationPayment,
        SignatureWithSaltAndExpiry memory paymentApproverSignature
    ) internal {
        require(!isPaymentApproverSaltUsed[paymentApproverSignature.salt], "EigenDAPaymentManager._verifyPaymentApproverSignature: paymentApprover salt already used");
        require(paymentApproverSignature.expiry > block.timestamp, "EigenDAPaymentManager._verifyPaymentApproverSignature: paymentApprover signature expired");   

        isPaymentApproverSaltUsed[paymentApproverSignature.salt] = true;    

        EIP1271SignatureUtils.checkSignature_EIP1271(
            paymentApprover, 
            calculatePaymentApprovalDigestHash(reservationPayment, paymentApproverSignature.salt, paymentApproverSignature.expiry), 
            paymentApproverSignature.signature
        );
    }

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
    ) public view returns (bytes32) {
        return _hashTypedDataV4(keccak256(abi.encode(PAYMENT_APPROVAL_TYPEHASH, reservationPayment, salt, expiry)));
    }
    
}