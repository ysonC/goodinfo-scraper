#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 <date>"
  echo "Example: $0 2025-11-14"
  exit 1
}

if [[ $# -ne 1 ]]; then
  usage
fi

zipDate="$1"
rawZip="raw-${zipDate}.zip"
combinedZip="combined-${zipDate}.zip"

if [[ ! -d "downloaded_stock" ]]; then
  echo "downloaded_stock directory not found"
  exit 1
fi

if [[ ! -d "final_output" ]]; then
  echo "final_output directory not found"
  exit 1
fi

echo "Creating ${rawZip} from downloaded_stock..."
zip -r -q "${rawZip}" downloaded_stock
echo "Done."

echo "Creating ${combinedZip} from final_output..."
zip -r -q "${combinedZip}" final_output
echo "Done."
