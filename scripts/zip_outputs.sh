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
dataDir="data"
rawDir="${dataDir}/downloaded_stock"
finalDir="${dataDir}/final_output"
archiveDir="${dataDir}/archives/${zipDate}"

mkdir -p "${archiveDir}"

if [[ ! -d "${rawDir}" ]]; then
  echo "${rawDir} directory not found"
  exit 1
fi

if [[ ! -d "${finalDir}" ]]; then
  echo "${finalDir} directory not found"
  exit 1
fi

echo "Creating ${rawZip} from ${rawDir}..."
zip -r -q "${rawZip}" "${rawDir}"
echo "Done."

echo "Creating ${combinedZip} from ${finalDir}..."
zip -r -q "${combinedZip}" "${finalDir}"
echo "Done."

echo "Moving files to archive folder..."
mv "${rawZip}" "${combinedZip}" "${archiveDir}"
echo "Done."
