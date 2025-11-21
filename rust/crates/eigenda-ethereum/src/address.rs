use std::str::FromStr;

use alloy_primitives::Address;
use alloy_primitives::AddressError;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

/// Ethereum address wrapper to implement jsonSchema trait.
/// This is needed to comply with the sovereign sdk.
/// See [crate::provider::EigenDaProviderConfig] for more details.
#[derive(Debug, derive_more::Display, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize)]
pub struct EthereumAddress(Address);

impl JsonSchema for EthereumAddress {
    fn schema_name() -> String {
        "EthereumAddress".to_string()
    }

    fn json_schema(_generator: &mut schemars::r#gen::SchemaGenerator) -> schemars::schema::Schema {
        serde_json::from_value(serde_json::json!({
            "type": "string",
            "pattern": "^0x[a-fA-F0-9]{40}$",
            "description": "An Ethereum address",
        }))
        .expect("valid schema")
    }
}

impl From<Address> for EthereumAddress {
    fn from(value: Address) -> Self {
        Self(value)
    }
}

impl From<EthereumAddress> for Address {
    fn from(value: EthereumAddress) -> Self {
        value.0
    }
}

impl FromStr for EthereumAddress {
    type Err = AddressError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(EthereumAddress(Address::parse_checksummed(s, None)?))
    }
}

#[cfg(test)]
mod tests {
    use std::str::FromStr;

    use super::EthereumAddress;

    const ADDR_1: &str = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";

    #[test]
    fn test_address_debug_from_string() {
        let raw_address_str = ADDR_1;
        let address = EthereumAddress::from_str(raw_address_str).unwrap();
        let output = format!("{address}");
        assert_eq!(raw_address_str, output);
    }

    #[test]
    fn test_address_conversion() {
        let raw_address_str = ADDR_1;
        let address = EthereumAddress::from_str(raw_address_str).unwrap();
        let eth_address: alloy_primitives::Address = address.into();
        let address_back: EthereumAddress = eth_address.into();
        assert_eq!(address, address_back);
    }
}
