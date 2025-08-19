use alloy_primitives::Uint;
use ark_bn254::{Fq, Fq2, G1Affine, G2Affine};
use ark_ec::AffineRepr;
use ark_ff::PrimeField;
use eigenda_cert::{BlobCommitment, G1Point, G2Point};

use crate::eigenda::verification::cert::convert;

pub trait DefaultExt: Sized {
    fn default_ext() -> Self;
}

impl DefaultExt for BlobCommitment {
    fn default_ext() -> Self {
        Self {
            commitment: G1Point::default_ext(),
            length_commitment: G2Point::default_ext(),
            length_proof: G2Point::default_ext(),
            length: u32::default(),
        }
    }
}

impl DefaultExt for G1Point {
    fn default_ext() -> Self {
        Self {
            x: Uint::ZERO,
            y: Uint::ZERO,
        }
    }
}

impl DefaultExt for G2Point {
    fn default_ext() -> Self {
        Self {
            x: vec![Uint::ZERO, Uint::ZERO],
            y: vec![Uint::ZERO, Uint::ZERO],
        }
    }
}

pub trait FromExt<T>: Sized {
    fn from_ext(value: T) -> Self;
}

pub trait IntoExt<T>: Sized {
    fn into_ext(self) -> T;
}

impl<T, U> IntoExt<U> for T
where
    U: FromExt<T>,
{
    fn into_ext(self) -> U {
        FromExt::from_ext(self)
    }
}

impl FromExt<G1Affine> for G1Point {
    fn from_ext(affine: G1Affine) -> Self {
        match affine.xy() {
            Some((x, y)) => G1Point {
                x: convert::fq_to_uint(x),
                y: convert::fq_to_uint(y),
            },
            None => G1Point::default_ext(),
        }
    }
}

impl FromExt<G2Affine> for G2Point {
    fn from_ext(affine: G2Affine) -> Self {
        match affine.xy() {
            Some((x, y)) => G2Point {
                x: vec![convert::fq_to_uint(x.c0), convert::fq_to_uint(x.c1)],
                y: vec![convert::fq_to_uint(y.c0), convert::fq_to_uint(y.c1)],
            },
            None => G2Point::default_ext(),
        }
    }
}

impl FromExt<G1Point> for G1Affine {
    fn from_ext(point: G1Point) -> G1Affine {
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

impl FromExt<G2Point> for G2Affine {
    fn from_ext(point: G2Point) -> Self {
        if point.x[0].is_zero()
            && point.y[0].is_zero()
            && point.x[1].is_zero()
            && point.y[1].is_zero()
        {
            return G2Affine::identity();
        }

        let x_c0_bytes: [u8; 32] = point.x[0].to_be_bytes();
        let x_c1_bytes: [u8; 32] = point.x[1].to_be_bytes();
        let y_c0_bytes: [u8; 32] = point.y[0].to_be_bytes();
        let y_c1_bytes: [u8; 32] = point.y[1].to_be_bytes();

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
    use alloy_primitives::uint;

    #[test]
    fn test_point_to_affine() {
        let point = G1Point {
            x: uint!(123456789_U256),
            y: uint!(987654321_U256),
        };

        let affine: G1Affine = point.into_ext();
        assert!(!affine.is_zero());
    }

    #[test]
    fn test_affine_to_point() {
        let x = Fq::from_be_bytes_mod_order(&[
            0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10, 0x11,
            0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
            0x20,
        ]);
        let y = Fq::from_be_bytes_mod_order(&[
            0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e,
            0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c,
            0x3d, 0x3e, 0x3f, 0x40,
        ]);
        let point = G1Affine::new_unchecked(x, y);

        let converted: G1Point = point.into_ext();
        let back_converted: G1Affine = converted.into_ext();

        assert_eq!(point, back_converted);
    }

    #[test]
    fn test_affine_to_point_identity() {
        let affine = G1Affine::identity();
        let point: G1Point = affine.into_ext();

        assert_eq!(point.x, Uint::ZERO);
        assert_eq!(point.y, Uint::ZERO);
    }

    #[test]
    fn test_point_to_affine_zero() {
        let point = G1Point {
            x: Uint::ZERO,
            y: Uint::ZERO,
        };

        let affine: G1Affine = point.into_ext();
        assert_eq!(affine, G1Affine::identity());
    }

    #[test]
    fn test_point_to_affine_g2() {
        let point = G2Point {
            x: vec![uint!(123456789_U256), uint!(111222333_U256)],
            y: vec![uint!(987654321_U256), uint!(444555666_U256)],
        };

        let affine: G2Affine = point.into_ext();
        assert!(!affine.is_zero());
    }

    #[test]
    fn test_affine_to_point_g2() {
        let x_c0 = Fq::from_be_bytes_mod_order(&[
            0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10, 0x11,
            0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
            0x20,
        ]);
        let x_c1 = Fq::from_be_bytes_mod_order(&[
            0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e,
            0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c,
            0x3d, 0x3e, 0x3f, 0x40,
        ]);
        let y_c0 = Fq::from_be_bytes_mod_order(&[
            0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e,
            0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c,
            0x5d, 0x5e, 0x5f, 0x60,
        ]);
        let y_c1 = Fq::from_be_bytes_mod_order(&[
            0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e,
            0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c,
            0x7d, 0x7e, 0x7f, 0x80,
        ]);
        let x = Fq2::new(x_c0, x_c1);
        let y = Fq2::new(y_c0, y_c1);
        let affine = G2Affine::new_unchecked(x, y);

        let converted: G2Point = affine.into_ext();
        let back_converted: G2Affine = converted.into_ext();

        assert_eq!(affine, back_converted);
    }

    #[test]
    fn test_affine_to_point_identity_g2() {
        let affine = G2Affine::identity();
        let point: G2Point = affine.into_ext();

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

        let affine: G2Affine = point.into_ext();
        assert_eq!(affine, G2Affine::identity());
    }
}
