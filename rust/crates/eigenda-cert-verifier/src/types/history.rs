use hashbrown::HashMap;

use crate::{
    error::CertVerificationError::{self, *},
    types::BlockNumber,
};

#[derive(Default, Debug, Clone)]
pub struct History<T: Copy>(pub HashMap<u32, Update<T>>);

impl<T: Copy> History<T> {
    pub(crate) fn try_get_at(&self, index: u32) -> Result<Update<T>, CertVerificationError> {
        self.0.get(&index).copied().ok_or(MissingHistoryEntry)
    }
}

#[derive(Default, Debug, Copy, Clone)]
pub struct Update<T: Copy> {
    interval: Interval<BlockNumber>,
    value: T,
}

impl<T: Copy> Update<T> {
    pub fn new(
        update_block: BlockNumber,
        next_update_block: BlockNumber,
        value: T,
    ) -> Result<Self, CertVerificationError> {
        let interval = Interval::new(update_block, next_update_block)?;
        let update = Self { interval, value };
        Ok(update)
    }

    pub(crate) fn try_get_against(
        &self,
        reference_block: BlockNumber,
    ) -> Result<T, CertVerificationError> {
        self.interval
            .contains(reference_block)
            .then_some(self.value)
            .ok_or(ElementNotInInterval)
    }
}

#[derive(Default, Debug, Clone, Copy)]
pub(crate) struct Interval<T: PartialOrd> {
    left_inclusive: T,
    right_exclusive: T,
}

impl<T: PartialOrd> Interval<T> {
    pub fn new(left_inclusive: T, right_exclusive: T) -> Result<Self, CertVerificationError> {
        match left_inclusive < right_exclusive {
            true => Ok(Self {
                left_inclusive,
                right_exclusive,
            }),
            false => Err(DegenerateInterval),
        }
    }

    pub fn contains(&self, element: T) -> bool {
        element >= self.left_inclusive && element < self.right_exclusive
    }
}

#[cfg(test)]
mod tests {
    use hashbrown::HashMap;

    use crate::error::CertVerificationError::*;
    use crate::types::BlockNumber;
    use crate::types::history::{History, Interval};

    use super::Update;

    #[test]
    fn element_before_left_is_not_in_interval() {
        let interval = Interval::new(42, 52).unwrap();
        assert_eq!(interval.contains(41), false);
    }

    #[test]
    fn element_at_left_is_in_interval() {
        let interval = Interval::new(42, 52).unwrap();
        assert_eq!(interval.contains(42), true);
    }

    #[test]
    fn element_in_interval() {
        let interval = Interval::new(42, 52).unwrap();
        assert_eq!(interval.contains(43), true);
    }

    #[test]
    fn element_at_right_is_not_in_interval() {
        let interval = Interval::new(42, 52).unwrap();
        assert_eq!(interval.contains(52), false);
    }

    #[test]
    fn element_after_right_is_not_in_interval() {
        let interval = Interval::new(42, 52).unwrap();
        assert_eq!(interval.contains(53), false);
    }

    #[test]
    fn degenerate_interval_where_left_equals_right() {
        let err = Interval::new(42, 42).unwrap_err();
        assert_eq!(err, DegenerateInterval);
    }

    #[test]
    fn degenerate_interval_where_left_greater_than_right() {
        let err = Interval::new(52, 42).unwrap_err();
        assert_eq!(err, DegenerateInterval);
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
