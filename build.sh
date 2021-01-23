#!/bin/zsh

# Clean up
rm -r bin/api/

# Build
for d in api/*/ ; do
  buildPath="./""$d""main"
  outPath="./bin/""${d%?}"

  echo "Building $buildPath to $outPath"

  env GOOS=linux go build -o "$outPath" "$buildPath"
done

# Package
for d in bin/api/* ; do
  zip "$d.zip" "$d"
done