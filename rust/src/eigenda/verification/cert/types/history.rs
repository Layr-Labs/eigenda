use std::fmt::Display;

use hashbrown::HashMap;
use thiserror::Error;

use crate::eigenda::verification::cert::types::BlockNumber;

#[derive(Debug, Error, PartialEq)]
pub enum HistoryError {
    #[error("Element ({0}) not in interval {0}")]
    ElementNotInInterval(String, String),

    #[error("Degenerate interval {0}")]
    DegenerateInterval(String),

    #[error("Missing history entry {0}")]
    MissingHistoryEntry(u32),
}

#[derive(Default, Debug, Clone)]
pub struct History<T: Copy + std::fmt::Debug>(pub HashMap<u32, Update<T>>);

impl<T: Copy + std::fmt::Debug> History<T> {
    pub(crate) fn try_get_at(&self, index: u32) -> Result<Update<T>, HistoryError> {
        use HistoryError::*;

        self.0
            .get(&index)
            .copied()
            .ok_or(MissingHistoryEntry(index))
    }
}

#[derive(Default, Debug, Copy, Clone)]
pub struct Update<T: Copy + std::fmt::Debug> {
    interval: Interval<BlockNumber>,
    value: T,
}

impl<T: Copy + std::fmt::Debug> Update<T> {
    pub fn new(
        update_block: BlockNumber,
        next_update_block: BlockNumber,
        value: T,
    ) -> Result<Self, HistoryError> {
        let interval = Interval::new(update_block, next_update_block)?;
        let update = Self { interval, value };
        Ok(update)
    }

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

#[derive(Default, Debug, Clone, Copy)]
pub(crate) struct Interval<T: PartialOrd + Display> {
    left_inclusive: T,
    right_exclusive: T,
}

impl<T: PartialOrd + Display> Display for Interval<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "[{}, {})", self.left_inclusive, self.right_exclusive)
    }
}

impl<T: PartialOrd + Display> Interval<T> {
    pub fn new(left_inclusive: T, right_exclusive: T) -> Result<Self, HistoryError> {
        use HistoryError::*;

        let is_valid = left_inclusive < right_exclusive;
        let interval = Self {
            left_inclusive,
            right_exclusive,
        };
        match is_valid {
            true => Ok(interval),
            false => Err(DegenerateInterval(interval.to_string())),
        }
    }

    pub fn contains(&self, element: T) -> bool {
        element >= self.left_inclusive && element < self.right_exclusive
    }
}

#[cfg(test)]
mod tests {
    use hashbrown::HashMap;

    use crate::eigenda::verification::cert::types::{
        BlockNumber,
        history::{History, HistoryError::*, Interval, Update},
    };

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
