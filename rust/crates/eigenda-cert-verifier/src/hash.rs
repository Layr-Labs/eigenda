use tiny_keccak::{Hasher, Keccak};

pub type BeHash = [u8; 32];

pub fn keccak256(inputs: &[&[u8]]) -> [u8; 32] {
    let mut hasher = Keccak::v256();
    for input in inputs {
        hasher.update(input);
    }
    let mut output = [0u8; 32];
    hasher.finalize(&mut output);
    output
}
