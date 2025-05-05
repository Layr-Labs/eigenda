// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {AccessControlLib} from "src/core/libraries/AccessControlLib.sol";
import {PaymentVaultLib, PaymentVaultStorage} from "src/core/libraries/v3/PaymentVaultLib.sol";
import {EigenDATypesV3} from "src/core/libraries/v3/EigenDATypesV3.sol";
import {Constants} from "src/core/libraries/Constants.sol";
import {IPaymentVault} from "src/core/interfaces/IPaymentVault.sol";

contract PaymentVault is IPaymentVault {
    modifier onlyOwner() {
        _onlyOwner();
        _;
    }

    modifier onlyQuorumOwner(uint64 quorumId) {
        _onlyQuorumOwner(quorumId);
        _;
    }

    // TODO: ADD INITIALIZER MODIFIER
    function initialize(address owner, uint64 schedulePeriod) external {
        AccessControlLib.grantRole(Constants.OWNER_ROLE, owner);
        PaymentVaultStorage.layout().schedulePeriod = schedulePeriod;
    }

    function createReservation(uint64 quorumId, EigenDATypesV3.Reservation memory reservation)
        external
        onlyQuorumOwner(quorumId)
    {
        PaymentVaultLib.createReservation(
            quorumId, msg.sender, reservation, PaymentVaultStorage.layout().schedulePeriod
        );
    }

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "New owner is the zero address");
        AccessControlLib.transferRole(Constants.OWNER_ROLE, msg.sender, newOwner);
    }

    function initializeQuorum(
        uint64 quorumId,
        address newOwner,
        EigenDATypesV3.QuorumPaymentProtocolConfig memory protocolCfg
    ) external onlyOwner {
        require(
            AccessControlLib.getRoleMemberCount(Constants.QUORUM_OWNER_ROLE(quorumId)) == 0, "Quorum owner already set"
        );
        AccessControlLib.grantRole(Constants.QUORUM_OWNER_ROLE(quorumId), newOwner);
        PaymentVaultStorage.layout().quorum[quorumId].protocolCfg = protocolCfg;
    }

    function setQuorumPaymentConfig(uint64 quorumId, EigenDATypesV3.QuorumPaymentConfig memory paymentConfig)
        external
        onlyQuorumOwner(quorumId)
    {
        PaymentVaultStorage.layout().quorum[quorumId].cfg = paymentConfig;
    }

    function transferQuorumOwnership(uint64 quorumId, address newOwner) external onlyQuorumOwner(quorumId) {
        require(newOwner != address(0), "New owner is the zero address");
        AccessControlLib.transferRole(Constants.QUORUM_OWNER_ROLE(quorumId), msg.sender, newOwner);
    }

    function _onlyOwner() internal view virtual {
        require(AccessControlLib.hasRole(Constants.OWNER_ROLE, msg.sender), "Not owner");
    }

    function _onlyQuorumOwner(uint64 quorumId) internal view virtual {
        require(AccessControlLib.hasRole(Constants.QUORUM_OWNER_ROLE(quorumId), msg.sender), "Not quorum owner");
    }
}
