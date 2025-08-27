use alloy_consensus::EthereumTxEnvelope;

/// Copied from:
/// * https://docs.rs/alloy-consensus/1.0.27/src/alloy_consensus/transaction/envelope.rs.html#264-272
#[cfg(feature = "native")]
pub fn map_eip4844(
    tx: EthereumTxEnvelope<alloy_consensus::TxEip4844Variant>,
) -> EthereumTxEnvelope<alloy_consensus::TxEip4844> {
    match tx {
        EthereumTxEnvelope::Legacy(tx) => EthereumTxEnvelope::Legacy(tx),
        EthereumTxEnvelope::Eip2930(tx) => EthereumTxEnvelope::Eip2930(tx),
        EthereumTxEnvelope::Eip1559(tx) => EthereumTxEnvelope::Eip1559(tx),
        EthereumTxEnvelope::Eip4844(tx) => {
            EthereumTxEnvelope::Eip4844(tx.map(alloy_consensus::TxEip4844::from))
        }
        EthereumTxEnvelope::Eip7702(tx) => EthereumTxEnvelope::Eip7702(tx),
    }
}

/// Copied from:
/// * https://docs.rs/alloy-consensus/1.0.27/src/alloy_consensus/transaction/typed.rs.html#697-710
/// * https://docs.rs/alloy-consensus/1.0.27/src/alloy_consensus/transaction/envelope.rs.html#532-538
pub mod serde_bincode_compat {
    extern crate alloc;

    use alloc::borrow::Cow;
    use alloy_consensus::Signed;
    use alloy_primitives::PrimitiveSignature;
    use serde::{Deserialize, Deserializer, Serialize, Serializer};
    use serde_with::{DeserializeAs, SerializeAs};

    /// Bincode-compatible [`super::EthereumTypedTransaction`] serde implementation.
    ///
    /// Intended to use with the [`serde_with::serde_as`] macro in the following way:
    /// ```rust
    /// use alloy_consensus::{serde_bincode_compat, EthereumTypedTransaction};
    /// use serde::{de::DeserializeOwned, Deserialize, Serialize};
    /// use serde_with::serde_as;
    ///
    /// #[serde_as]
    /// #[derive(Serialize, Deserialize)]
    /// struct Data<T: Serialize + DeserializeOwned + Clone + 'static> {
    ///     #[serde_as(as = "serde_bincode_compat::EthereumTypedTransaction<'_, T>")]
    ///     receipt: EthereumTypedTransaction<T>,
    /// }
    /// ```
    #[derive(Debug, Serialize, Deserialize)]
    pub enum EthereumTypedTransaction<'a, Eip4844: Clone = alloy_consensus::transaction::TxEip4844> {
        /// Legacy transaction
        Legacy(alloy_consensus::serde_bincode_compat::TxLegacy<'a>),
        /// EIP-2930 transaction
        Eip2930(alloy_consensus::serde_bincode_compat::TxEip2930<'a>),
        /// EIP-1559 transaction
        Eip1559(alloy_consensus::serde_bincode_compat::TxEip1559<'a>),
        /// EIP-4844 transaction
        /// Note: assumes EIP4844 is bincode compatible, which it is because no flatten or skipped
        /// fields.
        Eip4844(Cow<'a, Eip4844>),
        /// EIP-7702 transaction
        Eip7702(alloy_consensus::serde_bincode_compat::TxEip7702<'a>),
    }

    impl<'a, T: Clone> From<&'a alloy_consensus::EthereumTypedTransaction<T>>
        for EthereumTypedTransaction<'a, T>
    {
        fn from(value: &'a alloy_consensus::EthereumTypedTransaction<T>) -> Self {
            match value {
                alloy_consensus::EthereumTypedTransaction::Legacy(tx) => Self::Legacy(tx.into()),
                alloy_consensus::EthereumTypedTransaction::Eip2930(tx) => Self::Eip2930(tx.into()),
                alloy_consensus::EthereumTypedTransaction::Eip1559(tx) => Self::Eip1559(tx.into()),
                alloy_consensus::EthereumTypedTransaction::Eip4844(tx) => {
                    Self::Eip4844(Cow::Borrowed(tx))
                }
                alloy_consensus::EthereumTypedTransaction::Eip7702(tx) => Self::Eip7702(tx.into()),
            }
        }
    }

    impl<'a, T: Clone> From<EthereumTypedTransaction<'a, T>>
        for alloy_consensus::EthereumTypedTransaction<T>
    {
        fn from(value: EthereumTypedTransaction<'a, T>) -> Self {
            match value {
                EthereumTypedTransaction::Legacy(tx) => Self::Legacy(tx.into()),
                EthereumTypedTransaction::Eip2930(tx) => Self::Eip2930(tx.into()),
                EthereumTypedTransaction::Eip1559(tx) => Self::Eip1559(tx.into()),
                EthereumTypedTransaction::Eip4844(tx) => Self::Eip4844(tx.into_owned()),
                EthereumTypedTransaction::Eip7702(tx) => Self::Eip7702(tx.into()),
            }
        }
    }

    impl<T: Serialize + Clone> SerializeAs<alloy_consensus::EthereumTypedTransaction<T>>
        for EthereumTypedTransaction<'_, T>
    {
        fn serialize_as<S>(
            source: &alloy_consensus::EthereumTypedTransaction<T>,
            serializer: S,
        ) -> Result<S::Ok, S::Error>
        where
            S: Serializer,
        {
            EthereumTypedTransaction::<'_, T>::from(source).serialize(serializer)
        }
    }

    impl<'de, T: Deserialize<'de> + Clone>
        DeserializeAs<'de, alloy_consensus::EthereumTypedTransaction<T>>
        for EthereumTypedTransaction<'de, T>
    {
        fn deserialize_as<D>(
            deserializer: D,
        ) -> Result<alloy_consensus::EthereumTypedTransaction<T>, D::Error>
        where
            D: Deserializer<'de>,
        {
            EthereumTypedTransaction::<'_, T>::deserialize(deserializer).map(Into::into)
        }
    }

    /// Bincode-compatible [`super::EthereumTxEnvelope`] serde implementation.
    ///
    /// Intended to use with the [`serde_with::serde_as`] macro in the following way:
    /// ```rust
    /// use alloy_consensus::{serde_bincode_compat, EthereumTxEnvelope};
    /// use serde::{de::DeserializeOwned, Deserialize, Serialize};
    /// use serde_with::serde_as;
    ///
    /// #[serde_as]
    /// #[derive(Serialize, Deserialize)]
    /// struct Data<T: Serialize + DeserializeOwned + Clone + 'static> {
    ///     #[serde_as(as = "serde_bincode_compat::EthereumTxEnvelope<'_, T>")]
    ///     receipt: EthereumTxEnvelope<T>,
    /// }
    /// ```
    #[derive(Debug, Serialize, Deserialize)]
    pub struct EthereumTxEnvelope<'a, Eip4844: Clone = alloy_consensus::transaction::TxEip4844> {
        /// Transaction signature
        signature: PrimitiveSignature,
        /// bincode compatible transaction
        transaction: EthereumTypedTransaction<'a, Eip4844>,
    }

    impl<'a, T: Clone> From<&'a super::EthereumTxEnvelope<T>> for EthereumTxEnvelope<'a, T> {
        fn from(value: &'a super::EthereumTxEnvelope<T>) -> Self {
            match value {
                super::EthereumTxEnvelope::Legacy(tx) => Self {
                    signature: *tx.signature(),
                    transaction: EthereumTypedTransaction::Legacy(tx.tx().into()),
                },
                super::EthereumTxEnvelope::Eip2930(tx) => Self {
                    signature: *tx.signature(),
                    transaction: EthereumTypedTransaction::Eip2930(tx.tx().into()),
                },
                super::EthereumTxEnvelope::Eip1559(tx) => Self {
                    signature: *tx.signature(),
                    transaction: EthereumTypedTransaction::Eip1559(tx.tx().into()),
                },
                super::EthereumTxEnvelope::Eip4844(tx) => Self {
                    signature: *tx.signature(),
                    transaction: EthereumTypedTransaction::Eip4844(Cow::Borrowed(tx.tx())),
                },
                super::EthereumTxEnvelope::Eip7702(tx) => Self {
                    signature: *tx.signature(),
                    transaction: EthereumTypedTransaction::Eip7702(tx.tx().into()),
                },
            }
        }
    }

    impl<'a, T: Clone> From<EthereumTxEnvelope<'a, T>> for super::EthereumTxEnvelope<T> {
        fn from(value: EthereumTxEnvelope<'a, T>) -> Self {
            let EthereumTxEnvelope {
                signature,
                transaction,
            } = value;
            let transaction: alloy_consensus::transaction::EthereumTypedTransaction<T> =
                transaction.into();
            match transaction {
                alloy_consensus::transaction::EthereumTypedTransaction::Legacy(tx) => {
                    Signed::new_unhashed(tx, signature).into()
                }
                alloy_consensus::transaction::EthereumTypedTransaction::Eip2930(tx) => {
                    Signed::new_unhashed(tx, signature).into()
                }
                alloy_consensus::transaction::EthereumTypedTransaction::Eip1559(tx) => {
                    Signed::new_unhashed(tx, signature).into()
                }
                alloy_consensus::transaction::EthereumTypedTransaction::Eip4844(tx) => {
                    Self::Eip4844(Signed::new_unhashed(tx, signature))
                }
                alloy_consensus::transaction::EthereumTypedTransaction::Eip7702(tx) => {
                    Signed::new_unhashed(tx, signature).into()
                }
            }
        }
    }

    impl<T: Serialize + Clone> SerializeAs<super::EthereumTxEnvelope<T>> for EthereumTxEnvelope<'_, T> {
        fn serialize_as<S>(
            source: &super::EthereumTxEnvelope<T>,
            serializer: S,
        ) -> Result<S::Ok, S::Error>
        where
            S: Serializer,
        {
            EthereumTxEnvelope::<'_, T>::from(source).serialize(serializer)
        }
    }

    impl<'de, T: Deserialize<'de> + Clone> DeserializeAs<'de, super::EthereumTxEnvelope<T>>
        for EthereumTxEnvelope<'de, T>
    {
        fn deserialize_as<D>(deserializer: D) -> Result<super::EthereumTxEnvelope<T>, D::Error>
        where
            D: Deserializer<'de>,
        {
            EthereumTxEnvelope::<'_, T>::deserialize(deserializer).map(Into::into)
        }
    }

    #[cfg(test)]
    mod tests {
        use crate::ethereum::tx::serde_bincode_compat;
        use alloy_consensus::{
            EthereumTxEnvelope, EthereumTypedTransaction, transaction::TxEip4844,
        };
        use arbitrary::Arbitrary;
        use bincode::config;
        use rand::Rng;
        use serde::{Deserialize, Serialize};
        use serde_with::serde_as;

        #[test]
        fn test_typed_tx_envelope_bincode_roundtrip() {
            #[serde_as]
            #[derive(Debug, PartialEq, Eq, Serialize, Deserialize)]
            struct Data {
                #[serde_as(as = "serde_bincode_compat::EthereumTxEnvelope<'_>")]
                transaction: EthereumTxEnvelope<TxEip4844>,
            }

            let mut bytes = [0u8; 1024];
            rand::thread_rng().fill(bytes.as_mut_slice());
            let data = Data {
                transaction: EthereumTxEnvelope::arbitrary(&mut arbitrary::Unstructured::new(
                    &bytes,
                ))
                .unwrap(),
            };

            let encoded = bincode::serde::encode_to_vec(&data, config::legacy()).unwrap();
            let (decoded, _) =
                bincode::serde::decode_from_slice::<Data, _>(&encoded, config::legacy()).unwrap();
            assert_eq!(decoded, data);
        }

        #[test]
        fn test_typed_tx_bincode_roundtrip() {
            #[serde_as]
            #[derive(Debug, PartialEq, Eq, Serialize, Deserialize)]
            struct Data {
                #[serde_as(as = "serde_bincode_compat::EthereumTypedTransaction<'_>")]
                transaction: EthereumTypedTransaction<TxEip4844>,
            }

            let mut bytes = [0u8; 1024];
            rand::thread_rng().fill(bytes.as_mut_slice());
            let data = Data {
                transaction: EthereumTypedTransaction::arbitrary(
                    &mut arbitrary::Unstructured::new(&bytes),
                )
                .unwrap(),
            };

            let encoded = bincode::serde::encode_to_vec(&data, config::legacy()).unwrap();
            let (decoded, _) =
                bincode::serde::decode_from_slice::<Data, _>(&encoded, config::legacy()).unwrap();
            assert_eq!(decoded, data);
        }
    }
}
