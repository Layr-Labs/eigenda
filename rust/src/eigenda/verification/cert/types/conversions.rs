//! Type conversion utilities between EigenDA and arkworks representations
//!
//! This module provides conversion implementations for seamlessly converting
//! between EigenDA's Solidity-compatible types and arkworks' cryptographic types used
//! for elliptic curve operations.
//!
//! ## Key Conversions
//!
//! - **G1Point ↔ G1Affine**: Converts between EigenDA's 256-bit coordinate representation
//!   and arkworks' native BN254 G1 point format using standard `From`/`Into` traits
//! - **G2Point ↔ G2Affine**: Handles the more complex G2 field extension elements
//! - **Identity/Zero handling**: Properly maps between different representations of
//!   the point at infinity
//!
//! ## Design Principles
//!
//! - **Standard Rust traits**: Uses `From`/`Into` for type conversions, following Rust conventions
//! - **Bidirectional conversions**: All conversions are implemented in both directions
//! - **Field element ordering**: Correctly handles the [imaginary, real] vs [real, imaginary]
//!   difference between EigenDA and arkworks G2 representations
//!
//! ## Usage
//!
//! ```rust,ignore
//! use ark_bn254::G1Affine;
//! use sov_eigenda_adapter::eigenda::cert::G1Point;
//!
//! // Convert arkworks point to EigenDA format
//! let arkworks_point = G1Affine::generator();
//! let eigenda_point: G1Point = arkworks_point.into();
//!
//! // Convert back to arkworks format  
//! let back_to_arkworks: G1Affine = eigenda_point.into();
//! ```

use alloy_primitives::Uint;
use ark_bn254::{Fq, Fq2, G1Affine, G2Affine};
use ark_ec::AffineRepr;
use ark_ff::PrimeField;

use crate::eigenda::cert::{G1Point, G2Point};
use crate::eigenda::verification::cert::convert;

impl Default for G2Point {
    /// Create a default G2Point representing the point at infinity.
    ///
    /// Returns a G2Point with all coordinates set to zero, which represents
    /// the identity element (point at infinity) in EigenDA's representation.
    /// This is equivalent to the identity point in arkworks G2Affine.
    fn default() -> Self {
        Self {
            x: vec![Uint::ZERO, Uint::ZERO],
            y: vec![Uint::ZERO, Uint::ZERO],
        }
    }
}

impl From<G1Affine> for G1Point {
    /// Convert an arkworks G1Affine point to EigenDA's G1Point representation.
    ///
    /// Handles the identity/infinity point by returning a zero representation
    /// when the arkworks point is at infinity.
    fn from(affine: G1Affine) -> Self {
        match affine.xy() {
            Some((x, y)) => G1Point {
                x: convert::fq_to_uint(x),
                y: convert::fq_to_uint(y),
            },
            None => G1Point::default(),
        }
    }
}

impl From<G2Affine> for G2Point {
    /// Convert an arkworks G2Affine point to EigenDA's G2Point representation.
    ///
    /// **Important field element ordering difference:**
    /// - EigenDA points are represented as [imaginary, real]
    /// - arkworks points are represented as [real, imaginary]
    ///
    /// This conversion correctly maps between the two representations and
    /// handles the identity/infinity point by returning zeros.
    fn from(affine: G2Affine) -> Self {
        match affine.xy() {
            Some((x, y)) => G2Point {
                x: vec![convert::fq_to_uint(x.c1), convert::fq_to_uint(x.c0)],
                y: vec![convert::fq_to_uint(y.c1), convert::fq_to_uint(y.c0)],
            },
            None => G2Point::default(),
        }
    }
}

impl From<G1Point> for G1Affine {
    /// Convert EigenDA's G1Point representation to arkworks G1Affine.
    ///
    /// Detects the zero point (both coordinates zero) and maps it to
    /// arkworks' identity representation. Otherwise converts the 256-bit
    /// coordinates to field elements using big-endian byte order.
    ///
    /// Uses `new_unchecked` since we trust the input coordinates represent
    /// a valid curve point from EigenDA's verified data.
    fn from(point: G1Point) -> G1Affine {
        if point.x.is_zero() && point.y.is_zero() {
            return G1Affine::identity();
        }

        let x_bytes: [u8; 32] = point.x.to_be_bytes();
        let y_bytes: [u8; 32] = point.y.to_be_bytes();

        let x = Fq::from_be_bytes_mod_order(&x_bytes);
        let y = Fq::from_be_bytes_mod_order(&y_bytes);

        G1Affine::new_unchecked(x, y)
    }
}

impl From<G2Point> for G2Affine {
    /// Convert EigenDA's G2Point representation to arkworks G2Affine.
    ///
    /// **Important field element ordering difference:**
    /// - EigenDA points are represented as [imaginary, real]
    /// - arkworks points are represented as [real, imaginary]
    ///
    /// This conversion correctly maps between the two representations,
    /// detects zero points, and creates valid G2 field extension elements.
    ///
    /// Uses `new_unchecked` since we trust the input represents a valid
    /// curve point from EigenDA's verified data.
    fn from(point: G2Point) -> Self {
        if point.x[0].is_zero()
            && point.y[0].is_zero()
            && point.x[1].is_zero()
            && point.y[1].is_zero()
        {
            return G2Affine::identity();
        }

        let x_c0_bytes: [u8; 32] = point.x[1].to_be_bytes();
        let x_c1_bytes: [u8; 32] = point.x[0].to_be_bytes();
        let y_c0_bytes: [u8; 32] = point.y[1].to_be_bytes();
        let y_c1_bytes: [u8; 32] = point.y[0].to_be_bytes();

        let x_c0 = Fq::from_be_bytes_mod_order(&x_c0_bytes);
        let x_c1 = Fq::from_be_bytes_mod_order(&x_c1_bytes);
        let y_c0 = Fq::from_be_bytes_mod_order(&y_c0_bytes);
        let y_c1 = Fq::from_be_bytes_mod_order(&y_c1_bytes);

        let x = Fq2::new(x_c0, x_c1);
        let y = Fq2::new(y_c0, y_c1);

        G2Affine::new_unchecked(x, y)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_point_to_affine() {
        // Use readable hex string instead of uint! macro
        let point = G1Point {
            x: "0x00000000000000000000000000000000000000000000000000000000075bcd15"
                .parse()
                .unwrap(),
            y: "0x000000000000000000000000000000000000000000000000000000003ade68b1"
                .parse()
                .unwrap(),
        };

        let affine: G1Affine = point.into();
        assert!(!affine.is_zero());
    }

    #[test]
    fn test_affine_to_point() {
        // Use hex string for better readability
        let x_hex = "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20";
        let y_hex = "0x2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f40";

        let x_bytes = hex::decode(&x_hex[2..]).unwrap();
        let y_bytes = hex::decode(&y_hex[2..]).unwrap();

        let mut x_array = [0u8; 32];
        let mut y_array = [0u8; 32];
        x_array.copy_from_slice(&x_bytes);
        y_array.copy_from_slice(&y_bytes);

        let x = Fq::from_be_bytes_mod_order(&x_array);
        let y = Fq::from_be_bytes_mod_order(&y_array);
        let point = G1Affine::new_unchecked(x, y);

        let converted: G1Point = point.into();
        let back_converted: G1Affine = converted.into();

        assert_eq!(point, back_converted);
    }

    #[test]
    fn test_affine_to_point_identity() {
        let affine = G1Affine::identity();
        let point: G1Point = affine.into();

        assert_eq!(point.x, Uint::ZERO);
        assert_eq!(point.y, Uint::ZERO);
    }

    #[test]
    fn test_point_to_affine_zero() {
        let point = G1Point {
            x: Uint::ZERO,
            y: Uint::ZERO,
        };

        let affine: G1Affine = point.into();
        assert_eq!(affine, G1Affine::identity());
    }

    #[test]
    fn test_point_to_affine_g2() {
        // Use readable hex strings for G2 coordinates
        let point = G2Point {
            x: vec![
                "0x00000000000000000000000000000000000000000000000000000000075bcd15"
                    .parse()
                    .unwrap(),
                "0x000000000000000000000000000000000000000000000000000000006a24222"
                    .parse()
                    .unwrap(),
            ],
            y: vec![
                "0x000000000000000000000000000000000000000000000000000000003ade68b1"
                    .parse()
                    .unwrap(),
                "0x000000000000000000000000000000000000000000000000000000001a7dd93a"
                    .parse()
                    .unwrap(),
            ],
        };

        let affine: G2Affine = point.into();
        assert!(!affine.is_zero());
    }

    #[test]
    fn test_affine_to_point_g2() {
        // Use hex strings for better readability
        let x_c0_hex = "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20";
        let x_c1_hex = "0x2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f40";
        let y_c0_hex = "0x4142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f60";
        let y_c1_hex = "0x6162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f80";

        let x_c0_bytes = hex::decode(&x_c0_hex[2..]).unwrap();
        let x_c1_bytes = hex::decode(&x_c1_hex[2..]).unwrap();
        let y_c0_bytes = hex::decode(&y_c0_hex[2..]).unwrap();
        let y_c1_bytes = hex::decode(&y_c1_hex[2..]).unwrap();

        let mut x_c0_array = [0u8; 32];
        let mut x_c1_array = [0u8; 32];
        let mut y_c0_array = [0u8; 32];
        let mut y_c1_array = [0u8; 32];

        x_c0_array.copy_from_slice(&x_c0_bytes);
        x_c1_array.copy_from_slice(&x_c1_bytes);
        y_c0_array.copy_from_slice(&y_c0_bytes);
        y_c1_array.copy_from_slice(&y_c1_bytes);

        let x_c0 = Fq::from_be_bytes_mod_order(&x_c0_array);
        let x_c1 = Fq::from_be_bytes_mod_order(&x_c1_array);
        let y_c0 = Fq::from_be_bytes_mod_order(&y_c0_array);
        let y_c1 = Fq::from_be_bytes_mod_order(&y_c1_array);

        let x = Fq2::new(x_c0, x_c1);
        let y = Fq2::new(y_c0, y_c1);
        let affine = G2Affine::new_unchecked(x, y);

        let converted: G2Point = affine.into();
        let back_converted: G2Affine = converted.into();

        assert_eq!(affine, back_converted);
    }

    #[test]
    fn test_affine_to_point_identity_g2() {
        let affine = G2Affine::identity();
        let point: G2Point = affine.into();

        assert_eq!(point.x[0], Uint::ZERO);
        assert_eq!(point.x[1], Uint::ZERO);
        assert_eq!(point.y[0], Uint::ZERO);
        assert_eq!(point.y[1], Uint::ZERO);
    }

    #[test]
    fn test_point_to_affine_zero_g2() {
        let point = G2Point {
            x: vec![Uint::ZERO, Uint::ZERO],
            y: vec![Uint::ZERO, Uint::ZERO],
        };

        let affine: G2Affine = point.into();
        assert_eq!(affine, G2Affine::identity());
    }
}
