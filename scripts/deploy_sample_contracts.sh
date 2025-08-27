#!/bin/bash

# Script to deploy Sui sample contracts to a local network
# Assumes local network is already running

set -e  # Exit immediately if a command exits with a non-zero status

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Define the Sui RPC URL (allow overriding via environment variable)
SUI_RPC_URL=${SUI_RPC_URL:-"http://localhost:9000"}
echo -e "${YELLOW}Using Sui RPC URL: ${SUI_RPC_URL}${NC}"

# NOTE: This can be useful if Sui CLI is used instead of the docker image
# echo -e "${YELLOW}Setting Sui environment to local...${NC}"
# sui client switch --env localnet

echo -e "${YELLOW}Starting deployment of Sui sample contracts...${NC}"

# Change to the contracts directory
cd "$(dirname "$0")/../contracts/test"

# Build the contracts
echo -e "${YELLOW}Building contracts...${NC}"
sui move build -d

# Get the build package path
# PACKAGE_PATH=$(realpath ./build/TestContract/package.json)
# echo -e "${YELLOW}Package path: ${PACKAGE_PATH}${NC}"

# Deploy the contracts
echo -e "${YELLOW}Deploying contracts...${NC}"
RESULT=$(sui client publish --gas-budget 20000000 -d)

echo -e "${YELLOW}Result: ${RESULT}${NC}"

# Extract the package ID from the result using grep and sed
PACKAGE_ID=$(echo "$RESULT" | grep "PackageID:" | sed -E 's/.*PackageID: (0x[0-9a-f]+).*/\1/')

if [ -z "$PACKAGE_ID" ]; then
    echo "Failed to extract package ID from deployment result"
    exit 1
fi

echo -e "${GREEN}Contracts deployed successfully!${NC}"
echo -e "${GREEN}Package ID: ${PACKAGE_ID}${NC}"

# Create a counter object
echo -e "${YELLOW}Creating a counter object...${NC}"
sui client call --package $PACKAGE_ID --module counter --function initialize --gas-budget 20000000

echo -e "${GREEN}Deployment complete!${NC}"
echo -e "${YELLOW}You can now interact with the deployed contracts.${NC}"
echo -e "${YELLOW}Example commands:${NC}"
echo -e "  sui client --url $SUI_RPC_URL call --package $PACKAGE_ID --module counter --function increment --args \$COUNTER_ID --gas-budget 10000000"
echo -e "  sui client --url $SUI_RPC_URL call --package $PACKAGE_ID --module counter --function increment_mult --args \$COUNTER_ID 5 10 --gas-budget 10000000"