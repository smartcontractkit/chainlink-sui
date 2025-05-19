
echo "Generating bindings for Move Sui contracts..."

# Build the bindings (add the path to contracts you want to generate bindings for)

# Test Package
go run bindgen/main.go --moveConfig ./contracts/test/ --input ./contracts/test/sources/counter.move --output ./bindings/generated/test/counter
go run bindgen/main.go --moveConfig ./contracts/test/ --input ./contracts/test/sources/complex.move --output ./bindings/generated/test/complex

