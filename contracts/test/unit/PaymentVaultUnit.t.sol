// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "src/payments/PaymentVault.sol";
import "src/interfaces/IPaymentVault.sol";
import "forge-std/Test.sol";
import "forge-std/StdStorage.sol";

contract PaymentVaultUnit is Test {
    using stdStorage for StdStorage;

    event ReservationUpdated(address indexed account, IPaymentVault.Reservation reservation);
    event OnDemandPaymentUpdated(address indexed account, uint80 onDemandPayment, uint80 totalDeposit);
    event GlobalSymbolsPerPeriodUpdated(uint64 previousValue, uint64 newValue);
    event ReservationPeriodIntervalUpdated(uint64 previousValue, uint64 newValue);
    event GlobalRatePeriodIntervalUpdated(uint64 previousValue, uint64 newValue);
    event PriceParamsUpdated(
        uint64 previousMinNumSymbols, 
        uint64 newMinNumSymbols, 
        uint64 previousPricePerSymbol, 
        uint64 newPricePerSymbol, 
        uint64 previousPriceUpdateCooldown, 
        uint64 newPriceUpdateCooldown
    );

    PaymentVault paymentVault;
    PaymentVault paymentVaultImplementation;
    ERC20 mockToken;

    address proxyAdmin = address(uint160(uint256(keccak256(abi.encodePacked("proxyAdmin")))));
    address initialOwner = address(uint160(uint256(keccak256(abi.encodePacked("initialOwner")))));
    address user = address(uint160(uint256(keccak256(abi.encodePacked("user")))));
    address user2 = address(uint160(uint256(keccak256(abi.encodePacked("user2")))));

    uint64 minNumSymbols = 1;
    uint64 globalSymbolsPerPeriod = 2;
    uint64 pricePerSymbol = 3;
    uint64 reservationPeriodInterval = 4;
    uint64 globalRatePeriodInterval = 5;
    uint64 priceUpdateCooldown = 6 days;

    bytes quorumNumbers = hex"0001";
    bytes quorumSplits = hex"3232";

    function setUp() virtual public {
        paymentVaultImplementation = new PaymentVault();

        paymentVault = PaymentVault(
            payable(
                address(
                    new TransparentUpgradeableProxy(
                        address(paymentVaultImplementation),
                        address(proxyAdmin),
                        abi.encodeWithSelector(
                            PaymentVault.initialize.selector,
                            initialOwner,
                            minNumSymbols,
                            pricePerSymbol,
                            priceUpdateCooldown,
                            globalSymbolsPerPeriod,
                            reservationPeriodInterval,
                            globalRatePeriodInterval
                        )
                    )
                )
            )
        );

        mockToken = new ERC20("Mock Token", "MOCK");
    }

    function test_initialize() public {
        require(paymentVault.owner() == initialOwner, "Owner is not set");
        assertEq(paymentVault.minNumSymbols(), minNumSymbols);
        assertEq(paymentVault.globalSymbolsPerPeriod(), globalSymbolsPerPeriod);
        assertEq(paymentVault.pricePerSymbol(), pricePerSymbol);
        assertEq(paymentVault.reservationPeriodInterval(), reservationPeriodInterval);
        assertEq(paymentVault.priceUpdateCooldown(), priceUpdateCooldown);
        assertEq(paymentVault.globalRatePeriodInterval(), globalRatePeriodInterval);

        vm.expectRevert("Initializable: contract is already initialized");
        paymentVault.initialize(address(0), 0, 0, 0, 0, 0, 0);
    }

    function test_setReservation(uint56 _seed) public {
        uint64 _symbolsPerSecond = uint64(_seed);
        uint64 _startTimestamp = uint64(_seed) + 1;
        uint64 _endTimestamp = uint64(_seed) + 2;

        address _account = address(uint160(_seed));

        IPaymentVault.Reservation memory reservation = IPaymentVault.Reservation({
            symbolsPerSecond: _symbolsPerSecond,
            startTimestamp: _startTimestamp,
            endTimestamp: _endTimestamp,
            quorumNumbers: quorumNumbers,
            quorumSplits: quorumSplits
        });

        vm.expectEmit(address(paymentVault));
        emit ReservationUpdated(_account, reservation);
        vm.prank(initialOwner);
        paymentVault.setReservation(_account, reservation);

        assertEq(keccak256(abi.encode(paymentVault.getReservation(_account))), keccak256(abi.encode(reservation)));
    }

    function test_setReservation_revertInvalidQuorumSplits() public {
        IPaymentVault.Reservation memory reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 100,
            startTimestamp: 101,
            endTimestamp: 102,
            quorumNumbers: hex"0001",
            quorumSplits: hex"3233"
        });

        vm.expectRevert("sum of quorumSplits must be 100");
        vm.prank(initialOwner);
        paymentVault.setReservation(user, reservation);

        reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 100,
            startTimestamp: 101,
            endTimestamp: 102,
            quorumNumbers: hex"0001",
            quorumSplits: hex"3231"
        });

        vm.expectRevert("sum of quorumSplits must be 100");
        vm.prank(initialOwner);
        paymentVault.setReservation(user, reservation);

        reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 100,
            startTimestamp: 101,
            endTimestamp: 102,
            quorumNumbers: hex"0001",
            quorumSplits: hex"323334"
        });

        vm.expectRevert("arrays must have the same length");
        vm.prank(initialOwner);
        paymentVault.setReservation(user, reservation);
    }

    function test_setReservation_revertInvalidTimestamps() public {
        IPaymentVault.Reservation memory reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 100,
            startTimestamp: 101,
            endTimestamp: 100,
            quorumNumbers: quorumNumbers,
            quorumSplits: quorumSplits
        });

        vm.expectRevert("end timestamp must be greater than start timestamp");
        vm.prank(initialOwner);
        paymentVault.setReservation(user, reservation);
    }

    function test_depositOnDemand() public {
        vm.deal(user, 200 ether);

        vm.expectEmit(address(paymentVault));
        emit OnDemandPaymentUpdated(user, 100 ether, 100 ether);
        vm.prank(user);
        paymentVault.depositOnDemand{value: 100 ether}(user);
        assertEq(paymentVault.getOnDemandTotalDeposit(user), 100 ether);

        vm.expectEmit(address(paymentVault));
        emit OnDemandPaymentUpdated(user, 100 ether, 200 ether);
        vm.prank(user);
        paymentVault.depositOnDemand{value: 100 ether}(user);
        assertEq(paymentVault.getOnDemandTotalDeposit(user), 200 ether);
    }

    function test_depositOnDemand_forOtherUser() public {
        vm.deal(user, 100 ether);
        address otherUser = address(uint160(420));

        vm.expectEmit(address(paymentVault));
        emit OnDemandPaymentUpdated(user2, 100 ether, 100 ether);
        vm.prank(user);
        paymentVault.depositOnDemand{value: 100 ether}(user2);
        assertEq(paymentVault.getOnDemandTotalDeposit(user2), 100 ether);
        assertEq(paymentVault.getOnDemandTotalDeposit(user), 0);
    }

    function test_depositOnDemand_fallback() public {
        vm.deal(user, 100 ether);

        vm.expectEmit(address(paymentVault));
        emit OnDemandPaymentUpdated(user, 100 ether, 100 ether);
        vm.prank(user);
        payable(paymentVault).call{value: 100 ether}(hex"69");
        assertEq(paymentVault.getOnDemandTotalDeposit(user), 100 ether);
    }

    function test_depositOnDemand_recieve() public {
        vm.deal(user, 100 ether);

        vm.expectEmit(address(paymentVault));
        emit OnDemandPaymentUpdated(user, 100 ether, 100 ether);
        vm.prank(user);
        payable(paymentVault).call{value: 100 ether}("");
        assertEq(paymentVault.getOnDemandTotalDeposit(user), 100 ether);
    }

    function test_depositOnDemand_revertUint80Overflow() public {
        vm.deal(user, uint256(type(uint80).max) + 1);
        vm.expectRevert("amount must be less than or equal to 80 bits");
        vm.prank(user);
        paymentVault.depositOnDemand{value: uint256(type(uint80).max) + 1}(user);
    }

    function test_setPriceParams() public {
        vm.warp(block.timestamp + priceUpdateCooldown);

        vm.expectEmit(address(paymentVault));
        emit PriceParamsUpdated(minNumSymbols, minNumSymbols + 1, pricePerSymbol, pricePerSymbol + 1, priceUpdateCooldown, priceUpdateCooldown + 1);
        vm.prank(initialOwner);
        paymentVault.setPriceParams(minNumSymbols + 1, pricePerSymbol + 1, priceUpdateCooldown + 1);

        assertEq(paymentVault.minNumSymbols(), minNumSymbols + 1);
        assertEq(paymentVault.pricePerSymbol(), pricePerSymbol + 1);
        assertEq(paymentVault.priceUpdateCooldown(), priceUpdateCooldown + 1);
        assertEq(paymentVault.lastPriceUpdateTime(), block.timestamp);
    }

    function test_setPriceParams_revertCooldownNotSurpassed() public {
        vm.warp(block.timestamp + priceUpdateCooldown - 1);

        vm.expectRevert("price update cooldown not surpassed");
        vm.prank(initialOwner);
        paymentVault.setPriceParams(minNumSymbols + 1, pricePerSymbol + 1, priceUpdateCooldown + 1);
    }

    function test_setGlobalRatePeriodInterval() public {
        vm.expectEmit(address(paymentVault));
        emit GlobalRatePeriodIntervalUpdated(globalRatePeriodInterval, globalRatePeriodInterval + 1);
        vm.prank(initialOwner);
        paymentVault.setGlobalRatePeriodInterval(globalRatePeriodInterval + 1);
        assertEq(paymentVault.globalRatePeriodInterval(), globalRatePeriodInterval + 1);
    }

    function test_setGlobalSymbolsPerPeriod() public {
        vm.expectEmit(address(paymentVault));
        emit GlobalSymbolsPerPeriodUpdated(globalSymbolsPerPeriod, globalSymbolsPerPeriod + 1);
        vm.prank(initialOwner);
        paymentVault.setGlobalSymbolsPerPeriod(globalSymbolsPerPeriod + 1);
        assertEq(paymentVault.globalSymbolsPerPeriod(), globalSymbolsPerPeriod + 1);
    }

    function test_setReservationPeriodInterval() public {
        vm.expectEmit(address(paymentVault));
        emit ReservationPeriodIntervalUpdated(reservationPeriodInterval, reservationPeriodInterval + 1);
        vm.prank(initialOwner);
        paymentVault.setReservationPeriodInterval(reservationPeriodInterval + 1);
        assertEq(paymentVault.reservationPeriodInterval(), reservationPeriodInterval + 1);
    }

    function test_withdraw() public {
        test_depositOnDemand();
        vm.prank(initialOwner);
        paymentVault.withdraw(100 ether);
        assertEq(address(paymentVault).balance, 100 ether);
    }

    function test_withdrawERC20() public {
        deal(address(mockToken), address(paymentVault), 100 ether);
        vm.prank(initialOwner);
        paymentVault.withdrawERC20(address(mockToken), 100 ether);
        assertEq(mockToken.balanceOf(address(initialOwner)), 100 ether);
    }

    function test_ownedFunctions() public {
        IPaymentVault.Reservation memory reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 100,
            startTimestamp: 101,
            endTimestamp: 102,
            quorumNumbers: quorumNumbers,
            quorumSplits: quorumSplits
        });

        vm.expectRevert("Ownable: caller is not the owner");
        paymentVault.setReservation(user, reservation);
        vm.expectRevert("Ownable: caller is not the owner");
        paymentVault.withdraw(100 ether);
        vm.expectRevert("Ownable: caller is not the owner");
        paymentVault.withdrawERC20(address(mockToken), 100 ether);
        vm.expectRevert("Ownable: caller is not the owner");
        paymentVault.setPriceParams(minNumSymbols + 1, pricePerSymbol + 1, priceUpdateCooldown + 1);
        vm.expectRevert("Ownable: caller is not the owner");
        paymentVault.setGlobalRatePeriodInterval(globalRatePeriodInterval + 1);
        vm.expectRevert("Ownable: caller is not the owner");
        paymentVault.setGlobalSymbolsPerPeriod(globalSymbolsPerPeriod + 1);
        vm.expectRevert("Ownable: caller is not the owner");
        paymentVault.setReservationPeriodInterval(reservationPeriodInterval + 1);
    }

    function test_getReservations() public {
        IPaymentVault.Reservation memory reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 100,
            startTimestamp: 101,
            endTimestamp: 102,
            quorumNumbers: quorumNumbers,
            quorumSplits: quorumSplits
        });

        IPaymentVault.Reservation memory reservation2 = IPaymentVault.Reservation({
            symbolsPerSecond: 200,
            startTimestamp: 201,
            endTimestamp: 202,
            quorumNumbers: hex"0203",
            quorumSplits: hex"0163"
        });

        vm.startPrank(initialOwner);
        paymentVault.setReservation(user, reservation);
        paymentVault.setReservation(user2, reservation2);
        vm.stopPrank();

        address[] memory accounts = new address[](2);
        accounts[0] = user;
        accounts[1] = user2;
        IPaymentVault.Reservation[] memory reservations = paymentVault.getReservations(accounts);
        assertEq(keccak256(abi.encode(reservations[0])), keccak256(abi.encode(reservation)));
        assertEq(keccak256(abi.encode(reservations[1])), keccak256(abi.encode(reservation2)));
    }

    function test_getOnDemandAmounts() public {
        vm.deal(user, 300 ether);

        vm.startPrank(user);
        paymentVault.depositOnDemand{value: 100 ether}(user);
        paymentVault.depositOnDemand{value: 200 ether}(user2);
        vm.stopPrank();

        address[] memory accounts = new address[](2);
        accounts[0] = user;
        accounts[1] = user2;

        uint80[] memory payments = paymentVault.getOnDemandTotalDeposits(accounts);
        assertEq(payments[0], 100 ether);
        assertEq(payments[1], 200 ether);
    }
}