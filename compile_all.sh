# This script compiles Peek for release on all platforms supported platforms.
# Author: fwuffyboi
# Last modified: 2024-03-20
# Supported platforms: Linux (x86_64)/amd64, Linux (aarch64)/arm64

# if in directory "peek", move to src directory
if [ -d "src" ]; then
    cd src
fi

# Compile for Linux (x86_64)/amd64
echo "Compiling for Linux (x86_64)/amd64..."
GOOS=linux GOARCH=amd64 go build -o peek_linux_amd64

# Compile for Linux (aarch64)/arm64
echo "Compiling for Linux (aarch64)/arm64..."
GOOS=linux GOARCH=arm64 go build -o peek_linux_arm64

# move compiled binaries to root directory
echo "Moving compiled binaries to root directory..."
mv peek_linux_amd64 ../peek_linux_amd64
mv peek_linux_arm64 ../peek_linux_arm64

# return to root directory
cd ..

# say we're done
echo "Done compiling Peek for all supported platforms!"