// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";

contract EigenDARelayRegistry is IEigenDARelayRegistry, OwnableUpgradeable {

    mapping(uint32 => string) public relayIdToURL;
    mapping(address => uint32) public relayAddressToId;
    mapping(uint32 => address) public relayIdToAddress;

    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _initialOwner
    ) external initializer {
        _transferOwnership(_initialOwner);
    }

    function setRelayURL(address relay, uint32 id, string memory relayURL) external onlyOwner {
        require(relayIdToAddress[id] == address(0), "EigenDARelayRegistry: relay id already taken");

        relayIdToURL[id] = relayURL;
        relayAddressToId[relay] = id;
        relayIdToAddress[id] = relay;
        emit RelayAdded(relay, id, relayURL);
    }

    function getRelayURL(uint32 id) external view returns (string memory) {
        return relayIdToURL[id];
    }

    function getRelayId(address relay) external view returns (uint32) {
        return relayAddressToId[relay];
    }

    function getRelayAddress(uint32 id) external view returns (address) {
        return relayIdToAddress[id];
    }
}
