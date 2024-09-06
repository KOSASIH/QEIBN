#!/bin/bash

# Set the deployment directory
DEPLOY_DIR="deploy"

# Create the deployment directory if it doesn't exist
if [ ! -d "$DEPLOY_DIR" ]; then
  mkdir "$DEPLOY_DIR"
fi

# Extract the build tarball
tar -xzf "build.tar.gz" -C "$DEPLOY_DIR"

# Set the environment variables
export ENTANGLEMENT_SIMULATOR_PORT=8080
export AI_MODEL_TRAINER_PORT=8081
export IBC_NODE_PORT=8082

# Start the entanglement simulator
cd "$DEPLOY_DIR/build"
./entanglement-simulator &

# Start the AI model trainer
./ai-model-trainer &

# Start the IBC node
./ibc-node &
