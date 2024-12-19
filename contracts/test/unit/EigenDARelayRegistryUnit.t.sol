// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "../MockEigenDADeployer.sol";

contract EigenDARelayRegistryUnit is MockEigenDADeployer {

    event RelayAdded(address indexed relay, uint32 indexed key, string relayURL);

    function setUp() virtual public {
        _deployDA();
    }

    function test_initalize() public {
        require(eigenDARelayRegistry.owner() == registryCoordinatorOwner, "EigenDARelayRegistry: owner is not set");
        vm.expectRevert("Initializable: contract is already initialized");
        eigenDARelayRegistry.initialize(address(this));
    }

    function test_addRelayInfo() public {
        RelayInfo memory relayInfo = RelayInfo({
            relayAddress: address(uint160(uint256(keccak256(abi.encodePacked("relay"))))),
            relayURL: "https://relay.com"
        });

        vm.expectEmit(address(eigenDARelayRegistry));
        emit RelayAdded(relayInfo.relayAddress, eigenDARelayRegistry.nextRelayKey(), relayInfo.relayURL);
        vm.prank(registryCoordinatorOwner);
        eigenDARelayRegistry.addRelayInfo(relayInfo);

        assertEq(eigenDARelayRegistry.relayKeyToAddress(eigenDARelayRegistry.nextRelayKey() - 1), relayInfo.relayAddress);
        assertEq(eigenDARelayRegistry.relayKeyToUrl(eigenDARelayRegistry.nextRelayKey() - 1), relayInfo.relayURL);
    }

    function test_addRelayInfo_revert_notOwner() public {
        RelayInfo memory relayInfo = RelayInfo({
            relayAddress: address(uint160(uint256(keccak256(abi.encodePacked("relay"))))),
            relayURL: "https://relay.com"
        });

        vm.expectRevert("Ownable: caller is not the owner");
        eigenDARelayRegistry.addRelayInfo(relayInfo);
    }
}
