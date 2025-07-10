// little-endian is more efficient but Ethereum relies on little-endian
// little-endian is more natural to reason about in the context of bitmaps
use crate::error::SignaturesVerificationError;
use bitvec::{bitarr, order::Lsb0};
use primitive_types::U256;

const MAX_BIT_INDICES_LENGTH: usize = 256;

pub fn bit_indices_to_bitmap(
    // implied little-endian
    bit_indices: &[u8],
    // implied little-endian
    upper_bound_bit_index: u8,
) -> Result<U256, SignaturesVerificationError<'_>> {
    use SignaturesVerificationError::*;

    // abort early here even though other checks (sorted + unique) would catch it
    let bit_indices_len = bit_indices.len();
    if bit_indices_len > MAX_BIT_INDICES_LENGTH {
        return Err(BitIndicesGreaterThanMaxLength {
            bit_indices_len,
            max_bit_indices_len: MAX_BIT_INDICES_LENGTH,
        });
    }

    if bit_indices_len == 0 {
        return Ok(U256::default());
    }

    let mut prev_bit_index = None;
    let mut bitmap = bitarr![u8, Lsb0; 0; 256];
    for bit_index in bit_indices {
        if Some(bit_index) < prev_bit_index {
            return Err(BitIndicesNotSorted { bit_indices });
        }

        if Some(bit_index) == prev_bit_index {
            return Err(BitIndicesNotUnique { bit_indices });
        }

        prev_bit_index = Some(&bit_index);

        bitmap.set(*bit_index as usize, true);
    }

    // safe to unwrap since empty bit_indices has already been checked
    if bit_indices.last().unwrap() >= &upper_bound_bit_index {
        return Err(BitIndexGreaterThanOrEqualToUpperBound {
            bit_indices,
            upper_bound_bit_index,
        });
    }

    let slice = bitmap.as_raw_slice();
    let bitmap = U256::from_little_endian(slice);
    Ok(bitmap)
}

#[cfg(test)]
mod tests {
    use super::*;
    use SignaturesVerificationError::*;
    use alloc::vec;

    #[test]
    fn test_bit_indices_to_bitmap_succeeds_given_empty_input() {
        let bit_indices = vec![];
        let upper_bound_bit_index = u8::MAX;
        let result = bit_indices_to_bitmap(&bit_indices, upper_bound_bit_index);
        assert_eq!(result.unwrap(), U256::default());
    }

    #[test]
    fn test_bit_indices_to_bitmap_fails_when_it_exceeds_max_len() {
        let bit_indices = vec![0u8; 257];
        let upper_bound_bit_index = u8::MAX;
        let result = bit_indices_to_bitmap(&bit_indices, upper_bound_bit_index);
        assert_eq!(
            result.unwrap_err(),
            BitIndicesGreaterThanMaxLength {
                bit_indices_len: bit_indices.len(),
                max_bit_indices_len: MAX_BIT_INDICES_LENGTH
            }
        );
    }

    #[test]
    fn test_bit_indices_to_bitmap_fails_if_not_sorted() {
        let bit_indices = vec![42u8, 41u8, 43u8];
        let upper_bound_bit_index = u8::MAX;
        let result = bit_indices_to_bitmap(&bit_indices, upper_bound_bit_index);
        assert_eq!(
            result.unwrap_err(),
            BitIndicesNotSorted {
                bit_indices: &bit_indices
            }
        );
    }

    #[test]
    fn test_bit_indices_to_bitmap_fails_if_greater_than_upper_bound() {
        let bit_indices = vec![40u8, 41u8, 43u8];
        let upper_bound_bit_index = 42u8;
        let result = bit_indices_to_bitmap(&bit_indices, upper_bound_bit_index);
        assert_eq!(
            result.unwrap_err(),
            BitIndexGreaterThanOrEqualToUpperBound {
                bit_indices: &bit_indices,
                upper_bound_bit_index,
            }
        );
    }

    #[test]
    fn test_bit_indices_to_bitmap_fails_if_equal_to_upper_bound() {
        let bit_indices = vec![40u8, 41u8, 42u8];
        let upper_bound_bit_index = 42u8;
        let result = bit_indices_to_bitmap(&bit_indices, upper_bound_bit_index);
        assert_eq!(
            result.unwrap_err(),
            BitIndexGreaterThanOrEqualToUpperBound {
                bit_indices: &bit_indices,
                upper_bound_bit_index,
            }
        );
    }

    #[test]
    fn test_bit_indices_to_bitmap_succeeds_when_setting_the_0th_bit() {
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        // index: | 255 | 254 | 253 | 252 |...|  3  |  2  | *1* | *0* |
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        // bits:  |  0  |  0  |  0  |  0  |...|  0  |  0  |  1  |  1  |
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        let bit_indices = vec![0u8, 1u8];
        let upper_bound_bit_index = u8::MAX;
        let result = bit_indices_to_bitmap(&bit_indices, upper_bound_bit_index);
        let actual = result.unwrap();
        let expected = U256::from(3u8);
        assert_eq!(actual, expected);
    }

    #[test]
    fn test_bit_indices_to_bitmap_succeeds_when_targeting_decimal_8_as_bitmap() {
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        // index: | 255 | 254 | 253 | 252 |...| *3* |  2  |  1  |  0  |
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        // bits:  |  0  |  0  |  0  |  0  |...|  1  |  0  |  0  |  0  |
        //        +-----+-----+-----+-----+...+-----+-----+-----+-----+
        let bit_indices = vec![3u8];
        let upper_bound_bit_index = u8::MAX;
        let result = bit_indices_to_bitmap(&bit_indices, upper_bound_bit_index);
        let actual = result.unwrap();

        let expected = U256::from(8u8);

        assert_eq!(actual, expected);
    }
}
