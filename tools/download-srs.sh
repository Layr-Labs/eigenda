#!/bin/bash

# Script to download and verify SRS files for EigenDA
# Usage: ./download-srs.sh <size_in_bytes>
# Example: ./download-srs.sh 16777216  # for 16MB

set -euo pipefail

# Maximum allowed size (set to 16GB)
MAX_SIZE_BYTES=$((16 * 1024 * 1024 * 1024))

# Check for required dependencies
function check_dependencies() {
    local missing_deps=0
    for cmd in curl sha256sum awk grep tr; do
        if ! command -v "$cmd" &> /dev/null; then
            echo "Error: Required command '$cmd' not found"
            missing_deps=1
        fi
    done
    
    if [ $missing_deps -eq 1 ]; then
        echo "Please install missing dependencies and try again."
        return 1
    fi
}

# Perform a curl request - fail immediately if download fails
function download_file() {
    local url="$1"
    local output="$2"
    local range="$3"
    local file_description="$4"
    
    echo "Downloading ${file_description}..."
    
    if [ -z "$range" ]; then
        if ! curl -f -S -o "$output" --progress-bar "$url"; then
            echo "Error: Download failed for ${file_description}"
            return 1
        fi
    else
        if ! curl -f -S -o "$output" -r "$range" --progress-bar "$url"; then
            echo "Error: Download failed for ${file_description}"
            return 1
        fi
    fi
    
    return 0
}

# Main function to download SRS files
function download_srs_files() {
    local size_bytes="$1"
    local g2_size_bytes=$((size_bytes * 2))
    local base_url="https://srs-mainnet.s3.amazonaws.com/kzg"
    local output_dir="srs-files"
    local hash_file="${output_dir}/srs-files-${size_bytes}.sha256"
    
    # Create output directory if it doesn't exist
    mkdir -p "${output_dir}"
    
    echo "Downloading SRS files to support blob size of ${size_bytes} bytes..."
    
    # Validate that the base URL is accessible and get file size information
    echo "Checking server availability and file sizes..."
    local g2_total_size
    local g1_total_size
    
    # Check g2.point size
    local curl_result
    if ! curl_result=$(curl -f -S -I "${base_url}/g2.point"); then
        echo "Error: Cannot access ${base_url}/g2.point. The URL may be incorrect or the server is unavailable."
        return 1
    fi
    
    g2_total_size=$(echo "$curl_result" | grep -i "Content-Length" | awk '{print $2}' | tr -d '\r')
    
    if [ -z "$g2_total_size" ]; then
        echo "Error: Could not determine g2.point file size"
        return 1
    fi
    
    # Check g1.point size
    if ! curl_result=$(curl -f -S -I "${base_url}/g1.point"); then
        echo "Error: Cannot access ${base_url}/g1.point. The URL may be incorrect or the server is unavailable."
        return 1
    fi
    
    g1_total_size=$(echo "$curl_result" | grep -i "Content-Length" | awk '{print $2}' | tr -d '\r')
    
    if [ -z "$g1_total_size" ]; then
        echo "Error: Could not determine g1.point file size"
        return 1
    fi
    
    echo "Total g1.point size: ${g1_total_size} bytes"
    echo "Total g2.point size: ${g2_total_size} bytes"
    
    # Validate that our request sizes are reasonable
    if [ "$size_bytes" -gt "$g1_total_size" ]; then
        echo "Error: Requested g1 size ($size_bytes bytes) is larger than the source g1.point file ($g1_total_size bytes)"
        return 1
    fi
    
    if [ "$g2_size_bytes" -gt "$g2_total_size" ]; then
        echo "Error: Requested g2 size ($g2_size_bytes bytes) is larger than the source g2.point file ($g2_total_size bytes)"
        return 1
    fi
    
    # Calculate the start position for g2.trailing.point
    # Since we want the last G2_SIZE_BYTES of the file
    local g2_trailing_start=$((g2_total_size - g2_size_bytes))
    
    # g2_trailing_start should never be negative due to the validation above,
    # but if it is somehow negative, we should fail rather than download incorrect data
    if [ "$g2_trailing_start" -lt 0 ]; then
        echo "Error: Cannot download g2.trailing.point of size ${g2_size_bytes} bytes - insufficient source data"
        return 1
    fi
    
    # Download g1.point file
    if ! download_file "${base_url}/g1.point" "${output_dir}/g1.point" "0-$((size_bytes - 1))" "g1.point (${size_bytes} bytes)"; then
        return 1
    fi
    
    # Download g2.point file
    if ! download_file "${base_url}/g2.point" "${output_dir}/g2.point" "0-$((g2_size_bytes - 1))" "g2.point (${g2_size_bytes} bytes)"; then
        return 1
    fi
    
    # Download g2.trailing.point file
    if ! download_file "${base_url}/g2.point" "${output_dir}/g2.trailing.point" "${g2_trailing_start}-$((g2_total_size - 1))" "g2.trailing.point (${g2_size_bytes} bytes from the end of g2.point)"; then
        return 1
    fi
    
    # Generate hashes for verification
    echo "Generating verification hashes..."
    
    # Check if all files exist before generating hashes
    for file in "${output_dir}/g1.point" "${output_dir}/g2.point" "${output_dir}/g2.trailing.point"; do
        if [ ! -f "$file" ]; then
            echo "Error: File $file not found. Cannot generate hashes."
            return 1
        fi
    done
    
    # Generate the hash file
    {
        echo "# SRS files hashes for blob size ${size_bytes} bytes"
        echo "# Generated on $(date)"
        echo "# Format: SHA256 (filename)"
        echo ""
        sha256sum "${output_dir}/g1.point" | sed "s|${output_dir}/||"
        sha256sum "${output_dir}/g2.point" | sed "s|${output_dir}/||"
        sha256sum "${output_dir}/g2.trailing.point" | sed "s|${output_dir}/||"
    } > "${hash_file}"
    
    echo "Download complete. Files saved to ${output_dir}/"
    echo "Verification hashes saved to ${hash_file}"
    echo ""
    echo "Files downloaded:"
    ls -lh "${output_dir}"
    echo ""
    echo "Verification hashes:"
    cat "${hash_file}"
    
    return 0
}

# Validate command line arguments
function validate_arguments() {
    if [ $# -ne 1 ]; then
        echo "Usage: $0 <size_in_bytes>"
        echo "Example: $0 16777216  # for 16MB"
        return 1
    fi
    
    local size_bytes=$1
    
    # Check if size is a number
    if ! [[ "$size_bytes" =~ ^[0-9]+$ ]]; then
        echo "Error: Size must be a positive integer"
        return 1
    fi
    
    # Check if size is reasonable
    if [ "$size_bytes" -lt 32 ]; then
        echo "Error: Size must be at least 32 bytes"
        return 1
    fi
    
    if [ "$size_bytes" -gt "$MAX_SIZE_BYTES" ]; then
        echo "Error: Size must be less than $MAX_SIZE_BYTES bytes (16GB)"
        return 1
    fi
    
    return 0
}

# Main script execution
function main() {
    # Check for dependencies first
    if ! check_dependencies; then
        return 1
    fi
    
    # Validate arguments
    if ! validate_arguments "$@"; then
        return 1
    fi
    
    # Download SRS files
    if ! download_srs_files "$1"; then
        return 1
    fi
    
    echo ""
    echo "SRS files download and verification completed successfully!"
    return 0
}

# Run the script
main "$@"
exit $?