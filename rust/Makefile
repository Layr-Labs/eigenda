lint:
# checks for errors by compiling all targets
	cargo check --all-targets --all-features
	cargo fmt --all -- --check
# actual linting
	cargo clippy --all-targets --all-features -- -D warnings -D missing-docs
# checks for unused deps
	cargo machete