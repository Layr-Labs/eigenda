// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "../MockEigenDADeployer.sol";

contract EigenDADisperserRegistryUnit is MockEigenDADeployer {

    event DisperserAdded(uint32 indexed key, address indexed disperser);

    function setUp() virtual public {
        _deployDA();
    }

    function test_initalize() public {
        assertEq(eigenDADisperserRegistry.owner(), registryCoordinatorOwner);
        vm.expectRevert("Initializable: contract is already initialized");
        eigenDADisperserRegistry.initialize(address(this));
    }

    function test_setDisperserInfo() public {
        uint32 disperserKey = 1;
        address disperserAddress = address(uint160(uint256(keccak256(abi.encodePacked("disperser")))));
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });

        vm.expectEmit(address(eigenDADisperserRegistry));
        emit DisperserAdded(disperserKey, disperserAddress);
        vm.prank(registryCoordinatorOwner);
        eigenDADisperserRegistry.setDisperserInfo(disperserKey, disperserInfo);

        assertEq(eigenDADisperserRegistry.disperserKeyToAddress(disperserKey), disperserAddress);
    }

    function test_setDisperserInfo_revert_notOwner() public {
        uint32 disperserKey = 1;
        address disperserAddress = address(uint160(uint256(keccak256(abi.encodePacked("disperser")))));
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
            });

        vm.expectRevert("Ownable: caller is not the owner");
        eigenDADisperserRegistry.setDisperserInfo(disperserKey, disperserInfo);
    }
}
