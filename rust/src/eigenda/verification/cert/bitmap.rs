//! Bitmap operations for quorum membership tracking
//!
//! This module provides utilities for working with bitmaps that represent
//! quorum membership in the EigenDA protocol. Bitmaps are used to efficiently
//! track which quorums an operator belongs to or which quorums are involved
//! in a particular operation.
//!
//! ## Usage
//!
//! Bitmaps are typically created from a list of bit indices:
//! ```rust,ignore
//! # // This test is ignored because the eigenda module is private
//! use sov_eigenda_adapter::eigenda::verification::cert::bitmap::bit_indices_to_bitmap;
//! use alloy_primitives::Bytes;
//!
//! # fn example() -> Result<(), Box<dyn std::error::Error>> {
//! let bit_indices: Bytes = vec![0, 2, 5].into(); // quorums 0, 2, and 5
//! let bitmap = bit_indices_to_bitmap(&bit_indices, Some(8))?;
//! # Ok(())
//! # }
//! ```

use alloy_primitives::Bytes;
use bitvec::array::BitArray;
use thiserror::Error;

/// Maximum number of bit indices supported (256 quorums max)
pub(crate) const MAX_BIT_INDICES_LENGTH: usize = 256;

/// Efficient bitmap representation using 4 usize values (256 bits total)
///
/// This allows tracking up to 256 quorums, which is sufficient for EigenDA's
/// current design. Uses little-endian bit ordering for efficiency.
pub type Bitmap = BitArray<[usize; 4]>;

/// Errors that can occur during bitmap operations
#[derive(Debug, Error, PartialEq)]
pub enum BitmapError {
    /// Too many bit indices provided (max 256)
    #[error("Bit indices length ({len}) exceeds max byte slice length ({max_len})")]
    IndicesGreaterThanMaxLength { len: usize, max_len: usize },

    /// Bit indices contain duplicate values
    #[error("Bit indices not unique")]
    IndicesNotUnique,

    /// Bit indices are not in ascending order
    #[error("Bit indices not ordered")]
    IndicesNotSorted,

    /// One or more bit indices exceed the specified upper bound
    #[error("One or more bit indices are greater than or equal to the provided upper bound")]
    IndexThanOrEqualToUpperBound,
}

/// Convert a list of bit indices to a bitmap representation.
///
/// Creates a bitmap where each bit index in the input list sets the corresponding
/// bit in the output bitmap. The input must be sorted in ascending order with
/// no duplicates.
///
/// # Arguments
/// * `bit_indices` - Sorted list of bit indices to set (0-255)
/// * `upper_bound_bit_index` - Maximum allowed bit index (optional, defaults to 255)
///
/// # Returns
/// A bitmap with the specified bits set
///
/// # Errors
/// Returns [`BitmapError`] if:
/// - Too many indices are provided (> 256)
/// - Indices are not sorted in ascending order
/// - Indices contain duplicates
/// - Any index exceeds the upper bound
///
/// # Examples
/// ```rust,ignore
/// # // This test is ignored because the eigenda module is private
/// use sov_eigenda_adapter::eigenda::verification::cert::bitmap::bit_indices_to_bitmap;
/// use alloy_primitives::Bytes;
///
/// # fn example() -> Result<(), Box<dyn std::error::Error>> {
/// let indices: Bytes = vec![0, 2, 5].into();
/// let bitmap = bit_indices_to_bitmap(&indices, Some(8))?;
/// // Results in bitmap with bits 0, 2, and 5 set
/// # Ok(())
/// # }
/// ```
pub fn bit_indices_to_bitmap(
    bit_indices: &Bytes,
    upper_bound_bit_index: Option<u8>,
) -> Result<Bitmap, BitmapError> {
    use BitmapError::*;
    use core::cmp::Ordering::*;

    let upper_bound_bit_index = upper_bound_bit_index.unwrap_or(u8::MAX);

    match bit_indices.len() {
        0 => Ok(Bitmap::default()),
        // abort early here even though other checks (sorted + unique) would catch it
        len if len > MAX_BIT_INDICES_LENGTH => Err(IndicesGreaterThanMaxLength {
            len,
            max_len: MAX_BIT_INDICES_LENGTH,
        }),
        _ => {
            // safe to unwrap since we're in a branch where bit_indices is non-empty
            if *bit_indices.last().unwrap() >= upper_bound_bit_index {
                return Err(IndexThanOrEqualToUpperBound);
            }

            let mut prev_bit_index = None;
            let mut bitmap = Bitmap::default();
            for bit_index in bit_indices {
                match Some(bit_index).cmp(&prev_bit_index) {
                    Less => return Err(IndicesNotSorted),
                    Equal => return Err(IndicesNotUnique),
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
    use crate::eigenda::verification::cert::bitmap::{
        Bitmap, BitmapError::*, MAX_BIT_INDICES_LENGTH, bit_indices_to_bitmap,
    };

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
        assert_eq!(
            result.unwrap_err(),
            IndicesGreaterThanMaxLength {
                len: 257,
                max_len: MAX_BIT_INDICES_LENGTH
            }
        );
    }

    #[test]
    fn bit_indices_to_bitmap_fails_if_not_sorted() {
        let bit_indices = vec![42u8, 41u8, 43u8];
        let upper_bound_bit_index = None;
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap_err(), IndicesNotSorted,);
    }

    #[test]
    fn bit_indices_to_bitmap_fails_if_greater_than_upper_bound() {
        let bit_indices = vec![40u8, 41u8, 43u8];
        let upper_bound_bit_index = Some(42);
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap_err(), IndexThanOrEqualToUpperBound,);
    }

    #[test]
    fn bit_indices_to_bitmap_fails_if_equal_to_upper_bound() {
        let bit_indices = vec![40u8, 41u8, 42u8];
        let upper_bound_bit_index = Some(42);
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap_err(), IndexThanOrEqualToUpperBound);
    }

    #[test]
    fn bit_indices_to_bitmap_fails_with_duplicate_bit_indices() {
        let bit_indices = vec![42u8, 42u8];
        let upper_bound_bit_index = Some(43);
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap_err(), IndicesNotUnique);
    }

    #[test]
    fn bit_indices_to_bitmap_succeeds_with_empty_input_and_zero_upper_bound() {
        let bit_indices = vec![];
        let upper_bound_bit_index = Some(0);
        let result = bit_indices_to_bitmap(&bit_indices.into(), upper_bound_bit_index);
        assert_eq!(result.unwrap(), Bitmap::default());
    }
}
