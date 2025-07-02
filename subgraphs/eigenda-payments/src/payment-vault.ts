import {
  GlobalRatePeriodIntervalUpdated as GlobalRatePeriodIntervalUpdatedEvent,
  GlobalSymbolsPerPeriodUpdated as GlobalSymbolsPerPeriodUpdatedEvent,
  Initialized as InitializedEvent,
  OnDemandPaymentUpdated as OnDemandPaymentUpdatedEvent,
  OwnershipTransferred as OwnershipTransferredEvent,
  PriceParamsUpdated as PriceParamsUpdatedEvent,
  ReservationPeriodIntervalUpdated as ReservationPeriodIntervalUpdatedEvent,
  ReservationUpdated as ReservationUpdatedEvent
} from "../generated/PaymentVault/PaymentVault"
import {
  GlobalRatePeriodIntervalUpdated,
  GlobalSymbolsPerPeriodUpdated,
  Initialized,
  OnDemandPaymentUpdated,
  OwnershipTransferred,
  PriceParamsUpdated,
  ReservationPeriodIntervalUpdated,
  ReservationUpdated,
  ActiveReservation
} from "../generated/schema"

export function handleGlobalRatePeriodIntervalUpdated(
  event: GlobalRatePeriodIntervalUpdatedEvent
): void {
  let entity = new GlobalRatePeriodIntervalUpdated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.previousValue = event.params.previousValue
  entity.newValue = event.params.newValue

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleGlobalSymbolsPerPeriodUpdated(
  event: GlobalSymbolsPerPeriodUpdatedEvent
): void {
  let entity = new GlobalSymbolsPerPeriodUpdated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.previousValue = event.params.previousValue
  entity.newValue = event.params.newValue

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleInitialized(event: InitializedEvent): void {
  let entity = new Initialized(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.version = event.params.version

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOnDemandPaymentUpdated(
  event: OnDemandPaymentUpdatedEvent
): void {
  let entity = new OnDemandPaymentUpdated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.account = event.params.account
  entity.onDemandPayment = event.params.onDemandPayment
  entity.totalDeposit = event.params.totalDeposit

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOwnershipTransferred(
  event: OwnershipTransferredEvent
): void {
  let entity = new OwnershipTransferred(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.previousOwner = event.params.previousOwner
  entity.newOwner = event.params.newOwner

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handlePriceParamsUpdated(event: PriceParamsUpdatedEvent): void {
  let entity = new PriceParamsUpdated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.previousMinNumSymbols = event.params.previousMinNumSymbols
  entity.newMinNumSymbols = event.params.newMinNumSymbols
  entity.previousPricePerSymbol = event.params.previousPricePerSymbol
  entity.newPricePerSymbol = event.params.newPricePerSymbol
  entity.previousPriceUpdateCooldown = event.params.previousPriceUpdateCooldown
  entity.newPriceUpdateCooldown = event.params.newPriceUpdateCooldown

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleReservationPeriodIntervalUpdated(
  event: ReservationPeriodIntervalUpdatedEvent
): void {
  let entity = new ReservationPeriodIntervalUpdated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.previousValue = event.params.previousValue
  entity.newValue = event.params.newValue

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleReservationUpdated(event: ReservationUpdatedEvent): void {
  let entity = new ReservationUpdated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.account = event.params.account
  entity.reservation_symbolsPerSecond =
    event.params.reservation.symbolsPerSecond
  entity.reservation_startTimestamp = event.params.reservation.startTimestamp
  entity.reservation_endTimestamp = event.params.reservation.endTimestamp
  entity.reservation_quorumNumbers = event.params.reservation.quorumNumbers
  entity.reservation_quorumSplits = event.params.reservation.quorumSplits

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()

  // Create or update the ActiveReservation entity for this account
  let activeReservation = ActiveReservation.load(event.params.account)
  if (activeReservation == null) {
    activeReservation = new ActiveReservation(event.params.account)
  }
  
  activeReservation.account = event.params.account
  activeReservation.symbolsPerSecond = event.params.reservation.symbolsPerSecond
  activeReservation.startTimestamp = event.params.reservation.startTimestamp
  activeReservation.endTimestamp = event.params.reservation.endTimestamp
  activeReservation.quorumNumbers = event.params.reservation.quorumNumbers
  activeReservation.quorumSplits = event.params.reservation.quorumSplits
  activeReservation.lastUpdatedBlock = event.block.number
  activeReservation.lastUpdatedTimestamp = event.block.timestamp
  activeReservation.lastUpdatedTransactionHash = event.transaction.hash
  
  activeReservation.save()
}
