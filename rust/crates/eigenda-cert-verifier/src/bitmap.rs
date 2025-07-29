use alloy_primitives::Bytes;
// little-endian is more efficient but Ethereum relies on big-endian
// little-endian is more natural to reason about in the context of bitmaps
use bitvec::array::BitArray;

use crate::error::CertVerificationError;

const MAX_BIT_INDICES_LENGTH: usize = 256;

pub type Bitmap = BitArray<[u64; 4]>;

pub fn bit_indices_to_bitmap(
    bit_indices: &Bytes,
    upper_bound_bit_index: Option<u8>,
) -> Result<Bitmap, CertVerificationError> {
    use CertVerificationError::*;
    use core::cmp::Ordering::*;

    let upper_bound_bit_index = upper_bound_bit_index.unwrap_or(u8::MAX);

    match bit_indices.len() {
        0 => return Ok(Bitmap::default()),
        // abort early here even though other checks (sorted + unique) would catch it
        len if len > MAX_BIT_INDICES_LENGTH => return Err(BitIndicesGreaterThanMaxLength),
        _ => {
            // safe to unwrap since we're in a branch where bit_indices is non-empty
            if *bit_indices.last().unwrap() >= upper_bound_bit_index {
                return Err(BitIndexNotLessThanUpperBound);
            }

            let mut prev_bit_index = None;
            let mut bitmap = Bitmap::default();
            for bit_index in bit_indices {
                match Some(bit_index).cmp(&prev_bit_index) {
                    Less => return Err(BitIndicesNotSorted),
                    Equal => return Err(BitIndicesNotUnique),
                    Greater => {
                        prev_bit_index = Some(bit_index);
                        bitmap.set(*bit_index as usize, true);
                    }
                }
            }

            Ok(bitmap)
        }
    }
}

#[cfg(test)]
mod tests {
    use alloc::vec;

    use crate::bitmap::{Bitmap, bit_indices_to_bitmap};
    use crate::error::CertVerificationError::*;

    #[test]
    fn bit_indices_to_bitmap_succeeds_given_empty_input() {
        let bit_indices = vec![];
        let upper_bound_bit_index = None;
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap(), Bitmap::default());
    }

    #[test]
    fn bit_indices_to_bitmap_succeeds_when_setting_the_0th_bit() {
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        // index: | 255 | 254 | 253 | 252 |...|  3  |  2  | *1* | *0* |
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        // bits:  |  0  |  0  |  0  |  0  |...|  0  |  0  |  1  |  1  |
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        let bit_indices = vec![0u8, 1u8];
        let upper_bound_bit_index = None;
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        let actual = result.unwrap();
        let mut expected = Bitmap::default();
        expected.set(0, true);
        expected.set(1, true);
        assert_eq!(actual, expected);
    }

    #[test]
    fn bit_indices_to_bitmap_succeeds_when_targeting_decimal_8_as_bitmap() {
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        // index: | 255 | 254 | 253 | 252 |...| *3* |  2  |  1  |  0  |
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        // bits:  |  0  |  0  |  0  |  0  |...|  1  |  0  |  0  |  0  |
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        let bit_indices = vec![3u8];
        let upper_bound_bit_index = None;
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        let actual = result.unwrap();

        let mut expected = Bitmap::default();
        expected.set(3, true);

        assert_eq!(actual, expected);
    }

    #[test]
    fn bit_indices_to_bitmap_fails_when_it_exceeds_max_len() {
        let bit_indices = vec![0u8; 257];
        let upper_bound_bit_index = None;
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap_err(), BitIndicesGreaterThanMaxLength,);
    }

    #[test]
    fn bit_indices_to_bitmap_fails_if_not_sorted() {
        let bit_indices = vec![42u8, 41u8, 43u8];
        let upper_bound_bit_index = None;
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap_err(), BitIndicesNotSorted,);
    }

    #[test]
    fn bit_indices_to_bitmap_fails_if_greater_than_upper_bound() {
        let bit_indices = vec![40u8, 41u8, 43u8];
        let upper_bound_bit_index = Some(42);
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap_err(), BitIndexNotLessThanUpperBound,);
    }

    #[test]
    fn bit_indices_to_bitmap_fails_if_equal_to_upper_bound() {
        let bit_indices = vec![40u8, 41u8, 42u8];
        let upper_bound_bit_index = Some(42);
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap_err(), BitIndexNotLessThanUpperBound);
    }

    #[test]
    fn bit_indices_to_bitmap_fails_with_duplicate_bit_indices() {
        let bit_indices = vec![42u8, 42u8];
        let upper_bound_bit_index = Some(43);
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap_err(), BitIndicesNotUnique);
    }

    #[test]
    fn bit_indices_to_bitmap_succeeds_with_empty_input_and_zero_upper_bound() {
        let bit_indices = vec![];
        let upper_bound_bit_index = Some(0);
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap(), Bitmap::default());
    }
}
