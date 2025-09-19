use std::fmt::Display;

use hashbrown::HashMap;
use thiserror::Error;

use crate::eigenda::verification::cert::types::BlockNumber;

/// Errors that can occur when working with historical data structures.
///
/// These errors typically arise when trying to access historical operator state
/// data at specific block heights, such as when block ranges are invalid or
/// data is missing from the historical record.
#[derive(Debug, Error, PartialEq)]
pub enum HistoryError {
    /// The requested block number is not within the valid interval for this historical data
    #[error("Element ({0}) not in interval {1}")]
    ElementNotInInterval(String, String),

    /// The historical interval is invalid (e.g., start block >= end block)
    #[error("Degenerate interval {0}")]
    DegenerateInterval(String),

    /// No historical entry exists at the requested index
    #[error("Missing history entry {0}")]
    MissingHistoryEntry(u32),

    /// Invalid block order (update_block >= next_update_block when next_update_block != 0)
    #[error(
        "Invalid block order: update block {update_block} >= next update block {next_update_block}"
    )]
    InvalidBlockOrder {
        update_block: u32,
        next_update_block: u32,
    },
}

/// Historical data structure that tracks values over block ranges.
///
/// Stores a mapping of indices to `Update` objects, each containing a value
/// and the block range during which it was valid. This is used to track
/// historical operator states, stakes, and other time-dependent data in
/// EigenDA's on-chain contracts.
///
/// The generic parameter `T` represents the type of value being tracked
/// (e.g., stake amounts, truncated hashes, quorum bitmaps).
#[derive(Default, Debug, Clone)]
pub struct History<T: Copy + std::fmt::Debug>(pub HashMap<u32, Update<T>>);

impl<T: Copy + std::fmt::Debug> History<T> {
    /// Retrieve a historical update entry at the specified index.
    ///
    /// # Arguments
    /// * `index` - Index of the historical update to retrieve
    ///
    /// # Returns
    /// The `Update<T>` at the specified index
    ///
    /// # Errors
    /// Returns `HistoryError::MissingHistoryEntry` if no entry exists at the given index
    pub(crate) fn try_get_at(&self, index: u32) -> Result<Update<T>, HistoryError> {
        use HistoryError::*;

        self.0
            .get(&index)
            .copied()
            .ok_or(MissingHistoryEntry(index))
    }
}

/// A single update entry in historical data with an associated validity interval.
///
/// Contains a value and the block number range during which this value was active.
/// The interval is left-inclusive and right-exclusive: [start_block, end_block).
/// A `right_exclusive` value of 0 indicates the update is still current.
#[derive(Default, Debug, Copy, Clone)]
pub struct Update<T: Copy + std::fmt::Debug> {
    interval: Interval,
    value: T,
}

impl<T: Copy + std::fmt::Debug> Update<T> {
    /// Create a new update with the specified block range and value.
    ///
    /// # Arguments
    /// * `update_block` - Block number when this update became active
    /// * `next_update_block` - Block number when this update was superseded (0 means never)
    /// * `value` - The value associated with this update
    ///
    /// # Returns
    /// A new `Update` instance if the block range is valid
    ///
    /// # Errors
    /// Returns `HistoryError::InvalidBlockOrder` if `update_block >= next_update_block`
    /// (unless next_update_block is 0, which indicates the update is still current)
    pub fn new(
        update_block: BlockNumber,
        next_update_block: BlockNumber,
        value: T,
    ) -> Result<Self, HistoryError> {
        if next_update_block != 0 && update_block >= next_update_block {
            return Err(HistoryError::InvalidBlockOrder {
                update_block,
                next_update_block,
            });
        }

        let interval = Interval::new(update_block, next_update_block)?;
        let update = Self { interval, value };
        Ok(update)
    }

    pub fn update_block_number(&self) -> BlockNumber {
        self.interval.left_inclusive
    }

    pub fn next_update_block_number(&self) -> BlockNumber {
        self.interval.right_exclusive
    }

    pub fn value(&self) -> &T {
        &self.value
    }

    /// Retrieve the value from this update if it was valid at the given block number.
    ///
    /// Checks if the reference block number falls within this update's validity interval
    /// and returns the associated value if so.
    ///
    /// # Arguments
    /// * `reference_block` - Block number to check against this update's interval
    ///
    /// # Returns
    /// The value `T` if the reference block is within the validity interval
    ///
    /// # Errors
    /// Returns `HistoryError::ElementNotInInterval` if the reference block is outside
    /// the validity interval for this update
    pub(crate) fn try_get_against(&self, reference_block: BlockNumber) -> Result<T, HistoryError> {
        use HistoryError::*;

        self.interval
            .contains(reference_block)
            .then_some(self.value)
            .ok_or(ElementNotInInterval(
                reference_block.to_string(),
                self.interval.to_string(),
            ))
    }
}

/// A block number interval representing the validity period of a historical update.
///
/// Uses a half-open interval [left_inclusive, right_exclusive) where:
/// - `left_inclusive`: The first block where the update became valid
/// - `right_exclusive`: The first block where the update was superseded (exclusive)
///
/// A special case allows `right_exclusive = 0` to indicate the update is still current.
#[derive(Default, Debug, Clone, Copy)]
pub(crate) struct Interval {
    left_inclusive: BlockNumber,
    right_exclusive: BlockNumber,
}

impl Display for Interval {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "[{}, {})", self.left_inclusive, self.right_exclusive)
    }
}

impl Interval {
    /// Create a new interval with the specified block range.
    ///
    /// # Arguments
    /// * `left_inclusive` - First block where the interval is valid (inclusive)
    /// * `right_exclusive` - First block where the interval ends (exclusive, 0 means current)
    ///
    /// # Returns
    /// A valid `Interval` if the parameters are valid
    ///
    /// # Errors
    /// Returns `HistoryError::DegenerateInterval` if `left_inclusive >= right_exclusive`
    /// (unless `right_exclusive = 0`, which indicates the interval is still current)
    pub fn new(
        left_inclusive: BlockNumber,
        right_exclusive: BlockNumber,
    ) -> Result<Self, HistoryError> {
        use HistoryError::*;

        // special case `right_exclusive == 0` is allowed
        let is_valid = (left_inclusive < right_exclusive) || right_exclusive == 0;
        let interval = Self {
            left_inclusive,
            right_exclusive,
        };
        match is_valid {
            true => Ok(interval),
            false => Err(DegenerateInterval(interval.to_string())),
        }
    }

    /// Check if a block number falls within this interval.
    ///
    /// # Arguments
    /// * `element` - Block number to test for inclusion
    ///
    /// # Returns
    /// `true` if the block number is within [left_inclusive, right_exclusive),
    /// where `right_exclusive = 0` is treated as "no upper bound"
    pub fn contains(&self, element: BlockNumber) -> bool {
        element >= self.left_inclusive
            && (self.right_exclusive == 0 || element < self.right_exclusive)
    }
}

#[cfg(test)]
mod tests {
    use hashbrown::HashMap;

    use crate::eigenda::verification::cert::types::BlockNumber;
    use crate::eigenda::verification::cert::types::history::HistoryError::*;
    use crate::eigenda::verification::cert::types::history::{History, Interval, Update};

    #[test]
    fn element_before_left_is_not_in_interval() {
        let interval = Interval::new(42, 52).unwrap();
        assert!(!interval.contains(41));
    }

    #[test]
    fn element_at_left_is_in_interval() {
        let interval = Interval::new(42, 52).unwrap();
        assert!(interval.contains(42));
    }

    #[test]
    fn element_in_interval() {
        let interval = Interval::new(42, 52).unwrap();
        assert!(interval.contains(43));
    }

    #[test]
    fn element_at_right_is_not_in_interval() {
        let interval = Interval::new(42, 52).unwrap();
        assert!(!interval.contains(52));
    }

    #[test]
    fn element_after_right_is_not_in_interval() {
        let interval = Interval::new(42, 52).unwrap();
        assert!(!interval.contains(53));
    }

    #[test]
    fn degenerate_interval_where_left_equals_right() {
        let err = Interval::new(42, 42).unwrap_err();
        assert_eq!(err, DegenerateInterval("[42, 42)".into()));
    }

    #[test]
    fn degenerate_interval_where_left_greater_than_right() {
        let err = Interval::new(52, 42).unwrap_err();
        assert_eq!(err, DegenerateInterval("[52, 42)".into()));
    }

    #[test]
    fn new_update_with_invalid_inputs() {
        let result = Update::new(52, 42, 3);
        assert!(result.is_err());
    }

    #[test]
    fn try_get_update_against_valid_reference_block() {
        let value = 3;
        let update = Update::new(42, 52, value).unwrap();
        assert_eq!(update.try_get_against(43), Ok(value));
    }

    #[test]
    fn try_get_update_against_invalid_reference_block() {
        let value = 3;
        let update = Update::new(42, 52, value).unwrap();
        assert!(update.try_get_against(41).is_err());
    }

    #[test]
    fn try_get_history_entry_at_existing_index() {
        let history = HashMap::from([(42, Default::default())]);
        let history = History::<BlockNumber>(history);
        assert!(history.try_get_at(42).is_ok());
    }

    #[test]
    fn try_get_history_entry_at_missing_index() {
        let history = HashMap::from([(42, Default::default())]);
        let history = History::<BlockNumber>(history);
        assert!(history.try_get_at(52).is_err());
    }
}
