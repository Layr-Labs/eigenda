import { newMockEvent } from "matchstick-as"
import { ethereum, BigInt, Address } from "@graphprotocol/graph-ts"
import {
  GlobalRatePeriodIntervalUpdated,
  GlobalSymbolsPerPeriodUpdated,
  Initialized,
  OnDemandPaymentUpdated,
  OwnershipTransferred,
  PriceParamsUpdated,
  ReservationPeriodIntervalUpdated,
  ReservationUpdated
} from "../generated/PaymentVault/PaymentVault"

export function createGlobalRatePeriodIntervalUpdatedEvent(
  previousValue: BigInt,
  newValue: BigInt
): GlobalRatePeriodIntervalUpdated {
  let globalRatePeriodIntervalUpdatedEvent = changetype<
    GlobalRatePeriodIntervalUpdated
  >(newMockEvent())

  globalRatePeriodIntervalUpdatedEvent.parameters = new Array()

  globalRatePeriodIntervalUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "previousValue",
      ethereum.Value.fromUnsignedBigInt(previousValue)
    )
  )
  globalRatePeriodIntervalUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "newValue",
      ethereum.Value.fromUnsignedBigInt(newValue)
    )
  )

  return globalRatePeriodIntervalUpdatedEvent
}

export function createGlobalSymbolsPerPeriodUpdatedEvent(
  previousValue: BigInt,
  newValue: BigInt
): GlobalSymbolsPerPeriodUpdated {
  let globalSymbolsPerPeriodUpdatedEvent = changetype<
    GlobalSymbolsPerPeriodUpdated
  >(newMockEvent())

  globalSymbolsPerPeriodUpdatedEvent.parameters = new Array()

  globalSymbolsPerPeriodUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "previousValue",
      ethereum.Value.fromUnsignedBigInt(previousValue)
    )
  )
  globalSymbolsPerPeriodUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "newValue",
      ethereum.Value.fromUnsignedBigInt(newValue)
    )
  )

  return globalSymbolsPerPeriodUpdatedEvent
}

export function createInitializedEvent(version: i32): Initialized {
  let initializedEvent = changetype<Initialized>(newMockEvent())

  initializedEvent.parameters = new Array()

  initializedEvent.parameters.push(
    new ethereum.EventParam(
      "version",
      ethereum.Value.fromUnsignedBigInt(BigInt.fromI32(version))
    )
  )

  return initializedEvent
}

export function createOnDemandPaymentUpdatedEvent(
  account: Address,
  onDemandPayment: BigInt,
  totalDeposit: BigInt
): OnDemandPaymentUpdated {
  let onDemandPaymentUpdatedEvent = changetype<OnDemandPaymentUpdated>(
    newMockEvent()
  )

  onDemandPaymentUpdatedEvent.parameters = new Array()

  onDemandPaymentUpdatedEvent.parameters.push(
    new ethereum.EventParam("account", ethereum.Value.fromAddress(account))
  )
  onDemandPaymentUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "onDemandPayment",
      ethereum.Value.fromUnsignedBigInt(onDemandPayment)
    )
  )
  onDemandPaymentUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "totalDeposit",
      ethereum.Value.fromUnsignedBigInt(totalDeposit)
    )
  )

  return onDemandPaymentUpdatedEvent
}

export function createOwnershipTransferredEvent(
  previousOwner: Address,
  newOwner: Address
): OwnershipTransferred {
  let ownershipTransferredEvent = changetype<OwnershipTransferred>(
    newMockEvent()
  )

  ownershipTransferredEvent.parameters = new Array()

  ownershipTransferredEvent.parameters.push(
    new ethereum.EventParam(
      "previousOwner",
      ethereum.Value.fromAddress(previousOwner)
    )
  )
  ownershipTransferredEvent.parameters.push(
    new ethereum.EventParam("newOwner", ethereum.Value.fromAddress(newOwner))
  )

  return ownershipTransferredEvent
}

export function createPriceParamsUpdatedEvent(
  previousMinNumSymbols: BigInt,
  newMinNumSymbols: BigInt,
  previousPricePerSymbol: BigInt,
  newPricePerSymbol: BigInt,
  previousPriceUpdateCooldown: BigInt,
  newPriceUpdateCooldown: BigInt
): PriceParamsUpdated {
  let priceParamsUpdatedEvent = changetype<PriceParamsUpdated>(newMockEvent())

  priceParamsUpdatedEvent.parameters = new Array()

  priceParamsUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "previousMinNumSymbols",
      ethereum.Value.fromUnsignedBigInt(previousMinNumSymbols)
    )
  )
  priceParamsUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "newMinNumSymbols",
      ethereum.Value.fromUnsignedBigInt(newMinNumSymbols)
    )
  )
  priceParamsUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "previousPricePerSymbol",
      ethereum.Value.fromUnsignedBigInt(previousPricePerSymbol)
    )
  )
  priceParamsUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "newPricePerSymbol",
      ethereum.Value.fromUnsignedBigInt(newPricePerSymbol)
    )
  )
  priceParamsUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "previousPriceUpdateCooldown",
      ethereum.Value.fromUnsignedBigInt(previousPriceUpdateCooldown)
    )
  )
  priceParamsUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "newPriceUpdateCooldown",
      ethereum.Value.fromUnsignedBigInt(newPriceUpdateCooldown)
    )
  )

  return priceParamsUpdatedEvent
}

export function createReservationPeriodIntervalUpdatedEvent(
  previousValue: BigInt,
  newValue: BigInt
): ReservationPeriodIntervalUpdated {
  let reservationPeriodIntervalUpdatedEvent = changetype<
    ReservationPeriodIntervalUpdated
  >(newMockEvent())

  reservationPeriodIntervalUpdatedEvent.parameters = new Array()

  reservationPeriodIntervalUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "previousValue",
      ethereum.Value.fromUnsignedBigInt(previousValue)
    )
  )
  reservationPeriodIntervalUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "newValue",
      ethereum.Value.fromUnsignedBigInt(newValue)
    )
  )

  return reservationPeriodIntervalUpdatedEvent
}

export function createReservationUpdatedEvent(
  account: Address,
  reservation: ethereum.Tuple
): ReservationUpdated {
  let reservationUpdatedEvent = changetype<ReservationUpdated>(newMockEvent())

  reservationUpdatedEvent.parameters = new Array()

  reservationUpdatedEvent.parameters.push(
    new ethereum.EventParam("account", ethereum.Value.fromAddress(account))
  )
  reservationUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "reservation",
      ethereum.Value.fromTuple(reservation)
    )
  )

  return reservationUpdatedEvent
}
