#!/bin/bash

# Check if the directory path is provided as an argument
if [ -z "$1" ]; then
  echo "Usage: $0 <directory_path>"
  exit 1
fi

# Assign the provided directory path to a variable
directory="$1"

# Check if the provided path is a valid directory
if [ ! -d "$directory" ]; then
  echo "Error: $directory is not a valid directory."
  exit 1
fi

# Remove all existing .json files in the specified directory
rm -f "$directory"/*.json

# Loop through all .txt files in the specified directory
for file in "$directory"/*.txt; do
  # Check if the file exists (in case there are no .txt files)
  if [ -f "$file" ]; then
    # Read the contents of the file, remove leading zeroes, and convert to JSON array
    json_array=$(awk '{print $1 + 0}' "$file" | jq -R -s -c 'split("\n")[:-1]')

    # Generate the new filename by removing leading zeroes, extra words, and changing the extension to .json
    new_filename=$(basename "$file" | sed -E 's/^0+//' | sed 's/Chapter//g' | sed 's/\.txt$/.json/')

    # Write the JSON array to the new file
    echo "$json_array" > "$directory/$new_filename"

    # Generate the new .txt filename by removing leading zeroes and extra words
    new_txt_filename=$(basename "$file" | sed -E 's/^0+//' | sed 's/Chapter//g')

    # Rename the original .txt file without asking for permission
    mv -f "$file" "$directory/$new_txt_filename"
  fi
done
