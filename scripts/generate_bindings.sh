
echo "Generating bindings for Move contracts"

echo "Starting up the Sui Node"
# Sui node is needed generate contract bytecodes in base64
sui start --force-regenesis &
SUIPID=$!

echo "Waiting for Node to be available..."
# Wait for the Sui node to initialize
sleep 5

echo "Generating bindings..."

# Build the bindings (add the path to contracts you want to generate bindings for)
go run bindgen/main.go --moveConfig ./contracts/test/ --input ./contracts/test/sources/counter.move --output ./bindings/generated/counter
go run bindgen/main.go --moveConfig ./contracts/test/ --input ./contracts/test/sources/complex.move --output ./bindings/generated/complex

# Stop the Sui node process once finished
kill $SUIPID
