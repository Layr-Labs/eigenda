// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {Test} from "forge-std/Test.sol";

import {UsageAuthorizationLib} from "src/core/libraries/v3/usage-authorization/UsageAuthorizationLib.sol";
import {UsageAuthorizationTypes} from "src/core/libraries/v3/usage-authorization/UsageAuthorizationTypes.sol";
import {UsageAuthorizationStorage} from "src/core/libraries/v3/usage-authorization/UsageAuthorizationStorage.sol";
import {UsageAuthorizationRegistry, IUsageAuthorizationRegistry} from "src/core/UsageAuthorizationRegistry.sol";

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract ERC20Mintable is ERC20 {
    constructor(string memory name, string memory symbol) ERC20(name, symbol) {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract UsageAuthorizationRegistryTestHarness is UsageAuthorizationRegistry {
    constructor(uint64 schedulePeriod) UsageAuthorizationRegistry(schedulePeriod) {}

    function checkReservation(
        uint64 quorumId,
        UsageAuthorizationTypes.Reservation memory reservation,
        uint64 schedulePeriod
    ) external view {
        UsageAuthorizationLib.checkReservation(quorumId, reservation, schedulePeriod);
    }

    function increaseReservedSymbols(
        uint64 quorumId,
        uint64 startTimestamp,
        uint64 endTimestamp,
        uint64 symbolsPerSecond,
        uint64 schedulePeriod
    ) external {
        UsageAuthorizationLib.increaseReservedSymbols(
            quorumId, startTimestamp, endTimestamp, symbolsPerSecond, schedulePeriod
        );
    }

    function decreaseReservedSymbols(
        uint64 quorumId,
        uint64 startTimestamp,
        uint64 endTimestamp,
        uint64 symbolsPerSecond,
        uint64 schedulePeriod
    ) external {
        UsageAuthorizationLib.decreaseReservedSymbols(
            quorumId, startTimestamp, endTimestamp, symbolsPerSecond, schedulePeriod
        );
    }
}

contract UsageAuthorizationRegistryUnit is Test {
    UsageAuthorizationRegistryTestHarness usageAuthorizationRegistry;

    address constant OWNER = address(uint160(uint256(keccak256("EIGEN_DA_USAGE_AUTHORIZATION_REGISTRY_OWNER"))));
    address constant QUORUM_OWNER_0 =
        address(uint160(uint256(keccak256("EIGEN_DA_USAGE_AUTHORIZATION_REGISTRY_QUORUM_OWNER_0"))));

    uint64 constant SCHEDULE_PERIOD = 1 minutes;
    uint64 constant START_PERIOD = 100;
    uint64 constant MAX_NUM_PERIODS = 100000; // schedule period * maxNumPeriods should not overflow uint64
    uint64 constant RESERVATION_ADVANCE_PERIODS = 100;
    uint64 constant RESERVATION_ADVANCE_WINDOW = SCHEDULE_PERIOD * RESERVATION_ADVANCE_PERIODS;
    uint64 constant MIN_NUM_SYMBOLS = 2;
    uint64 constant RESERVATION_RATE_LIMIT_WINDOW = 3;
    uint64 constant ON_DEMAND_RATE_LIMIT_WINDOW = 5;
    uint64 constant RESERVATION_SYMBOLS_PER_SECOND = 100;
    uint64 constant ON_DEMAND_SYMBOLS_PER_SECOND = 200;
    uint64 constant ON_DEMAND_PRICE_PER_SYMBOL = 7;

    ERC20Mintable public token;
    address constant TEST_RECIPIENT =
        address(uint160(uint256(keccak256("EIGEN_DA_USAGE_AUTHORIZATION_REGISTRY_QUORUM_RECIPIENT"))));

    function setUp() public virtual {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        token = new ERC20Mintable("Test Token", "TTK");
        usageAuthorizationRegistry = new UsageAuthorizationRegistryTestHarness(SCHEDULE_PERIOD);
        usageAuthorizationRegistry.initialize(OWNER);
        vm.prank(OWNER);
        usageAuthorizationRegistry.initializeQuorum(
            0,
            QUORUM_OWNER_0,
            UsageAuthorizationTypes.QuorumProtocolConfig({
                minNumSymbols: MIN_NUM_SYMBOLS,
                reservationAdvanceWindow: RESERVATION_ADVANCE_WINDOW,
                reservationRateLimitWindow: RESERVATION_RATE_LIMIT_WINDOW,
                onDemandRateLimitWindow: ON_DEMAND_RATE_LIMIT_WINDOW,
                onDemandEnabled: true
            })
        );
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.setQuorumConfig(
            0,
            UsageAuthorizationTypes.QuorumConfig({
                token: address(token),
                recipient: TEST_RECIPIENT,
                reservationSymbolsPerSecond: RESERVATION_SYMBOLS_PER_SECOND,
                onDemandSymbolsPerSecond: ON_DEMAND_SYMBOLS_PER_SECOND,
                onDemandPricePerSymbol: ON_DEMAND_PRICE_PER_SYMBOL
            })
        );
    }

    function boundTimestamps(uint256 timestampSeed)
        internal
        view
        returns (uint64 startTimestamp, uint64 endTimestamp)
    {
        uint256 currentPeriod = uint64(block.timestamp / SCHEDULE_PERIOD);
        uint256 startPeriod = bound(timestampSeed >> 64, START_PERIOD, currentPeriod + RESERVATION_ADVANCE_PERIODS - 1);
        uint256 endPeriod = bound(timestampSeed, startPeriod + 1, currentPeriod + RESERVATION_ADVANCE_PERIODS);
        startTimestamp = uint64(startPeriod * SCHEDULE_PERIOD);
        endTimestamp = uint64(endPeriod * SCHEDULE_PERIOD);
    }

    /// @notice Tests that we can add a reservation successfully.
    function test_AddReservation(address account, uint256 timestampSeed, uint64 symbolsPerSecond) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });
        vm.expectEmit(true, true, true, true);
        emit UsageAuthorizationLib.ReservationAdded(quorumId, account, reservation);
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.addReservation(quorumId, account, reservation);
        UsageAuthorizationTypes.Reservation memory storedReservation =
            usageAuthorizationRegistry.getReservation(quorumId, account);
        assertEq(storedReservation.startTimestamp, reservation.startTimestamp);
        assertEq(storedReservation.endTimestamp, reservation.endTimestamp);
        assertEq(storedReservation.symbolsPerSecond, reservation.symbolsPerSecond);
        for (uint64 i = startTimestamp / SCHEDULE_PERIOD; i < endTimestamp / SCHEDULE_PERIOD; i++) {
            assertEq(usageAuthorizationRegistry.getQuorumReservedSymbols(quorumId, i), symbolsPerSecond);
        }
    }

    /// @notice Tests that adding a reservation reverts if a reservation is currently active.
    function test_AddReservationRevertsIfReservationStillActive(
        address account,
        uint256 timestampSeed,
        uint64 symbolsPerSecond
    ) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));
        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.addReservation(quorumId, account, reservation);
        vm.expectRevert(
            abi.encodeWithSelector(
                IUsageAuthorizationRegistry.ReservationStillActive.selector, reservation.endTimestamp
            )
        );
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.addReservation(quorumId, account, reservation);
    }

    /// @notice Tests that we can successfully pass a reservation check.
    function test_CheckReservation(uint256 timestampSeed, uint64 symbolsPerSecond) public view {
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });

        usageAuthorizationRegistry.checkReservation(quorumId, reservation, SCHEDULE_PERIOD);
    }

    /// @notice Tests that a start timestamp not in the schedule period reverts.
    function test_CheckReservationRevertsIfStartTimestampNotInSchedulePeriod(
        uint256 timestampSeed,
        uint64 symbolsPerSecond
    ) public {
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        startTimestamp += uint64(bound(timestampSeed >> 128, 1, SCHEDULE_PERIOD - 1));
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });

        vm.expectRevert(
            abi.encodeWithSelector(
                IUsageAuthorizationRegistry.TimestampSchedulePeriodMismatch.selector, startTimestamp, SCHEDULE_PERIOD
            )
        );
        usageAuthorizationRegistry.checkReservation(quorumId, reservation, SCHEDULE_PERIOD);
    }

    /// @notice Tests that an end timestamp not in the schedule period reverts.
    function test_CheckReservationRevertsIfEndTimestampNotInSchedulePeriod(
        uint256 timestampSeed,
        uint64 symbolsPerSecond
    ) public {
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        endTimestamp -= uint64(bound(timestampSeed >> 128, 1, SCHEDULE_PERIOD - 1));
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });

        vm.expectRevert(
            abi.encodeWithSelector(
                IUsageAuthorizationRegistry.TimestampSchedulePeriodMismatch.selector, endTimestamp, SCHEDULE_PERIOD
            )
        );
        usageAuthorizationRegistry.checkReservation(quorumId, reservation, SCHEDULE_PERIOD);
    }

    /// @notice Tests that start timestamp must not be greater than the end timestamp.
    function test_CheckReservationRevertsIfStartTimestampGreaterThanEndTimestamp(
        uint256 timestampSeed,
        uint64 symbolsPerSecond
    ) public {
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: endTimestamp,
            endTimestamp: startTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });

        vm.expectRevert(
            abi.encodeWithSelector(
                IUsageAuthorizationRegistry.InvalidReservationPeriod.selector,
                reservation.startTimestamp,
                reservation.endTimestamp
            )
        );
        usageAuthorizationRegistry.checkReservation(quorumId, reservation, SCHEDULE_PERIOD);
    }

    /// @notice Tests that reservation length cannot exceed the quorum's reservation advance window
    function test_CheckReservationRevertsIfReservationTooLong(uint256 timestampSeed, uint64 symbolsPerSecond) public {
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));

        // Set the end timestamp to be too far in the future
        endTimestamp = uint64(block.timestamp) + RESERVATION_ADVANCE_WINDOW + SCHEDULE_PERIOD;

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });

        vm.expectRevert(
            abi.encodeWithSelector(
                IUsageAuthorizationRegistry.ReservationTooLong.selector,
                endTimestamp - uint64(block.timestamp / SCHEDULE_PERIOD * SCHEDULE_PERIOD),
                RESERVATION_ADVANCE_WINDOW
            )
        );
        usageAuthorizationRegistry.checkReservation(quorumId, reservation, SCHEDULE_PERIOD);
    }

    /// @notice Tests that increasing a reservation's reserved symbols successfully increases the quorum's reserved symbols
    function test_IncreaseReservedSymbols(uint256 timestampSeed, uint64 symbolsPerSecond) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 0, 100));

        usageAuthorizationRegistry.increaseReservedSymbols(
            quorumId, startTimestamp, endTimestamp, symbolsPerSecond, SCHEDULE_PERIOD
        );

        for (uint64 i = startTimestamp / SCHEDULE_PERIOD; i < endTimestamp / SCHEDULE_PERIOD; i++) {
            assertEq(usageAuthorizationRegistry.getQuorumReservedSymbols(quorumId, i), symbolsPerSecond);
        }
    }

    /// @notice Tests that increasing a reservation's reserved symbols reverts if not enough symbols are available
    function test_IncreaseReservedSymbolsRevertsIfNotEnoughSymbolsAvailable(
        uint256 timestampSeed,
        uint64 symbolsPerSecond
    ) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        vm.assume(symbolsPerSecond > RESERVATION_SYMBOLS_PER_SECOND);

        vm.expectRevert(
            abi.encodeWithSelector(
                IUsageAuthorizationRegistry.NotEnoughSymbolsAvailable.selector,
                startTimestamp,
                symbolsPerSecond,
                RESERVATION_SYMBOLS_PER_SECOND
            )
        );
        usageAuthorizationRegistry.increaseReservedSymbols(
            quorumId, startTimestamp, endTimestamp, symbolsPerSecond, SCHEDULE_PERIOD
        );
    }

    /// @notice Tests that decreasing a reservation's reserved symbols successfully decreases the quorum's reserved symbols
    function test_DecreaseReservedSymbols(uint256 timestampSeed, uint64 symbolsPerSecond) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 0, 100));

        usageAuthorizationRegistry.increaseReservedSymbols(
            quorumId, startTimestamp, endTimestamp, symbolsPerSecond, SCHEDULE_PERIOD
        );
        usageAuthorizationRegistry.decreaseReservedSymbols(
            quorumId, startTimestamp, endTimestamp, symbolsPerSecond, SCHEDULE_PERIOD
        );

        for (uint64 i = startTimestamp / SCHEDULE_PERIOD; i < endTimestamp / SCHEDULE_PERIOD; i++) {
            assertEq(usageAuthorizationRegistry.getQuorumReservedSymbols(quorumId, i), 0);
        }
    }

    /// @notice Tests that a reservation can be increased successfully.
    function test_IncreaseReservation(
        address account,
        uint256 timestampSeed,
        uint64 symbolsPerSecond,
        uint64 periodIncrease,
        uint64 symbolIncrease
    ) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 50));
        periodIncrease = uint64(bound(periodIncrease, 0, 10));
        symbolIncrease = uint64(bound(symbolIncrease, 0, 50));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.addReservation(quorumId, account, reservation);

        vm.warp(block.timestamp + SCHEDULE_PERIOD * periodIncrease);
        vm.expectEmit(true, true, true, true);
        reservation.endTimestamp += periodIncrease * SCHEDULE_PERIOD;
        reservation.symbolsPerSecond += symbolIncrease;
        emit UsageAuthorizationLib.ReservationIncreased(quorumId, account, reservation);
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.increaseReservation(quorumId, account, reservation);

        UsageAuthorizationTypes.Reservation memory storedReservation =
            usageAuthorizationRegistry.getReservation(quorumId, account);
        assertEq(storedReservation.startTimestamp, reservation.startTimestamp);
        assertEq(storedReservation.endTimestamp, reservation.endTimestamp);
        assertEq(storedReservation.symbolsPerSecond, reservation.symbolsPerSecond);
    }

    /// @notice Tests that increasing a reservation reverts if the start timestamp does not match the reservation's start timestamp.
    function test_IncreaseReservationRevertsIfStartTimestampDoesNotMatch(
        address account,
        uint256 timestampSeed,
        uint64 symbolsPerSecond
    ) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.addReservation(quorumId, account, reservation);
        vm.warp(block.timestamp + SCHEDULE_PERIOD);

        vm.expectRevert(
            abi.encodeWithSelector(
                IUsageAuthorizationRegistry.StartTimestampMustMatch.selector, reservation.startTimestamp
            )
        );
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.increaseReservation(
            quorumId,
            account,
            UsageAuthorizationTypes.Reservation({
                startTimestamp: startTimestamp + SCHEDULE_PERIOD,
                endTimestamp: endTimestamp + SCHEDULE_PERIOD,
                symbolsPerSecond: symbolsPerSecond
            })
        );
    }

    /// @notice Tests that increasing a reservation reverts if the reservation decreases.
    function test_IncreaseReservationRevertsIfReservationDecreases(
        address account,
        uint256 timestampSeed,
        uint64 symbolsPerSecond
    ) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.addReservation(quorumId, account, reservation);

        vm.expectRevert(
            abi.encodeWithSelector(
                IUsageAuthorizationRegistry.ReservationMustIncrease.selector,
                reservation.endTimestamp,
                reservation.symbolsPerSecond
            )
        );
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.increaseReservation(
            quorumId,
            account,
            UsageAuthorizationTypes.Reservation({
                startTimestamp: startTimestamp,
                endTimestamp: endTimestamp,
                symbolsPerSecond: symbolsPerSecond - 1
            })
        );
    }

    /// @notice Tests that a reservation can be decreased successfully.
    function test_DecreaseReservation(address account, uint256 timestampSeed, uint64 symbolsPerSecond) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 50));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.addReservation(quorumId, account, reservation);

        reservation.symbolsPerSecond -= 1; // Decrease the symbols per second by 1
        vm.expectEmit(true, true, true, true);
        emit UsageAuthorizationLib.ReservationDecreased(quorumId, account, reservation);
        vm.prank(account);
        usageAuthorizationRegistry.decreaseReservation(quorumId, reservation);

        UsageAuthorizationTypes.Reservation memory storedReservation =
            usageAuthorizationRegistry.getReservation(quorumId, account);
        assertEq(storedReservation.startTimestamp, reservation.startTimestamp);
        assertEq(storedReservation.endTimestamp, reservation.endTimestamp);
        assertEq(storedReservation.symbolsPerSecond, reservation.symbolsPerSecond);
    }

    /// @notice Tests that decreasing a reservation reverts if the start timestamp does not match the reservation's start timestamp.
    function test_DecreaseReservationRevertsIfStartTimestampDoesNotMatch(
        address account,
        uint256 timestampSeed,
        uint64 symbolsPerSecond
    ) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.addReservation(quorumId, account, reservation);

        vm.expectRevert(
            abi.encodeWithSelector(
                IUsageAuthorizationRegistry.StartTimestampMustMatch.selector, reservation.startTimestamp
            )
        );
        vm.prank(account);
        vm.warp(block.timestamp + SCHEDULE_PERIOD);
        usageAuthorizationRegistry.decreaseReservation(
            quorumId,
            UsageAuthorizationTypes.Reservation({
                startTimestamp: startTimestamp + SCHEDULE_PERIOD,
                endTimestamp: endTimestamp + SCHEDULE_PERIOD,
                symbolsPerSecond: symbolsPerSecond
            })
        );
    }

    /// @notice Tests that decreasing a reservation reverts if the reservation increases.
    function test_DecreaseReservationRevertsIfReservationIncreases(
        address account,
        uint256 timestampSeed,
        uint64 symbolsPerSecond
    ) public {
        vm.warp(SCHEDULE_PERIOD * START_PERIOD);
        uint64 quorumId = 0;
        (uint64 startTimestamp, uint64 endTimestamp) = boundTimestamps(timestampSeed);
        symbolsPerSecond = uint64(bound(symbolsPerSecond, 1, 100));

        UsageAuthorizationTypes.Reservation memory reservation = UsageAuthorizationTypes.Reservation({
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            symbolsPerSecond: symbolsPerSecond
        });
        vm.prank(QUORUM_OWNER_0);
        usageAuthorizationRegistry.addReservation(quorumId, account, reservation);

        vm.expectRevert(
            abi.encodeWithSelector(
                IUsageAuthorizationRegistry.ReservationMustDecrease.selector,
                reservation.endTimestamp,
                reservation.symbolsPerSecond
            )
        );

        vm.prank(account);
        usageAuthorizationRegistry.decreaseReservation(
            quorumId,
            UsageAuthorizationTypes.Reservation({
                startTimestamp: startTimestamp,
                endTimestamp: endTimestamp,
                symbolsPerSecond: symbolsPerSecond + 1
            })
        );
    }

    /// @notice Tests that a user can deposit on demand successfully.
    function test_DepositOnDemand(address account, address payer, uint256 amount) public {
        vm.assume(account != address(0));
        vm.assume(payer != address(0));
        uint64 quorumId = 0;
        amount = bound(amount, 1, type(uint80).max);
        token.mint(payer, amount);

        vm.prank(payer);
        token.approve(address(usageAuthorizationRegistry), amount);

        vm.expectEmit(true, true, true, true);
        emit UsageAuthorizationLib.DepositOnDemand(quorumId, account, amount, payer);
        vm.prank(payer);
        usageAuthorizationRegistry.depositOnDemand(quorumId, account, amount);
        uint256 onDemandDeposit = usageAuthorizationRegistry.getOnDemandDeposit(quorumId, account);
        assertEq(onDemandDeposit, amount);
    }

    /// @notice Tests that depositing on demand fails if on demand is disabled.
    function test_DepositOnDemandRevertsIfOnDemandDisabled(address payer, address account, uint256 amount) public {
        vm.assume(account != address(0));
        vm.assume(payer != address(0));
        uint64 quorumId = 0;
        amount = bound(amount, 1, type(uint80).max);
        token.mint(payer, amount);
        vm.prank(OWNER);
        usageAuthorizationRegistry.setQuorumProtocolConfig(
            quorumId,
            UsageAuthorizationTypes.QuorumProtocolConfig({
                minNumSymbols: MIN_NUM_SYMBOLS,
                reservationAdvanceWindow: RESERVATION_ADVANCE_WINDOW,
                reservationRateLimitWindow: RESERVATION_RATE_LIMIT_WINDOW,
                onDemandRateLimitWindow: ON_DEMAND_RATE_LIMIT_WINDOW,
                onDemandEnabled: false
            })
        );

        vm.prank(payer);
        token.approve(address(usageAuthorizationRegistry), amount);

        vm.expectRevert(abi.encodeWithSelector(IUsageAuthorizationRegistry.OnDemandDisabled.selector, quorumId));
        vm.prank(payer);
        usageAuthorizationRegistry.depositOnDemand(quorumId, account, amount);
    }

    /// @notice Tests that these functions are properly gated to the owner.
    function testOnlyOwnerFunctions(address account) public {
        vm.assume(account != OWNER);
        vm.startPrank(account);

        vm.expectRevert();
        usageAuthorizationRegistry.transferOwnership(account);

        vm.expectRevert();
        usageAuthorizationRegistry.initializeQuorum(
            0,
            account,
            UsageAuthorizationTypes.QuorumProtocolConfig({
                minNumSymbols: MIN_NUM_SYMBOLS,
                reservationAdvanceWindow: RESERVATION_ADVANCE_WINDOW,
                reservationRateLimitWindow: RESERVATION_RATE_LIMIT_WINDOW,
                onDemandRateLimitWindow: ON_DEMAND_RATE_LIMIT_WINDOW,
                onDemandEnabled: true
            })
        );

        vm.expectRevert();
        usageAuthorizationRegistry.setQuorumProtocolConfig(
            0,
            UsageAuthorizationTypes.QuorumProtocolConfig({
                minNumSymbols: MIN_NUM_SYMBOLS,
                reservationAdvanceWindow: RESERVATION_ADVANCE_WINDOW,
                reservationRateLimitWindow: RESERVATION_RATE_LIMIT_WINDOW,
                onDemandRateLimitWindow: ON_DEMAND_RATE_LIMIT_WINDOW,
                onDemandEnabled: true
            })
        );
    }

    /// @notice Tests that these functions are properly gated to the quorum owner.
    function testOnlyQuorumOwnerFunctions(address account) public {
        vm.assume(account != QUORUM_OWNER_0);
        vm.startPrank(account);

        vm.expectRevert();
        usageAuthorizationRegistry.addReservation(
            0, account, UsageAuthorizationTypes.Reservation({startTimestamp: 0, endTimestamp: 0, symbolsPerSecond: 0})
        );

        vm.expectRevert();
        usageAuthorizationRegistry.increaseReservation(
            0, account, UsageAuthorizationTypes.Reservation({startTimestamp: 0, endTimestamp: 0, symbolsPerSecond: 0})
        );

        vm.expectRevert();
        usageAuthorizationRegistry.setQuorumConfig(
            0,
            UsageAuthorizationTypes.QuorumConfig({
                token: address(token),
                recipient: TEST_RECIPIENT,
                reservationSymbolsPerSecond: RESERVATION_SYMBOLS_PER_SECOND,
                onDemandSymbolsPerSecond: ON_DEMAND_SYMBOLS_PER_SECOND,
                onDemandPricePerSymbol: ON_DEMAND_PRICE_PER_SYMBOL
            })
        );

        vm.expectRevert();
        usageAuthorizationRegistry.transferQuorumOwnership(0, account);
    }
}
