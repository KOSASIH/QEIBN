#!/bin/bash

# Set the build directory
BUILD_DIR="build"

# Create the build directory if it doesn't exist
if [ ! -d "$BUILD_DIR" ]; then
  mkdir "$BUILD_DIR"
fi

# Build the entanglement simulator
cd tools/entanglement-simulator
go build -o "$BUILD_DIR/entanglement-simulator" entanglement-simulator.go
cd -

# Build the AI model trainer
cd ai-model-trainer
cargo build --release
cd -
mv ai-model-trainer/target/release/ai-model-trainer "$BUILD_DIR/ai-model-trainer"

# Build the IBC protocol
cd node-architecture
go build -o "$BUILD_DIR/ibc-node" node.go
cd -

# Create a tarball of the build directory
tar -czf "$BUILD_DIR.tar.gz" "$BUILD_DIR"
