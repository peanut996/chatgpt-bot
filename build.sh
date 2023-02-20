echo "Building bot..."
go build -ldflags '-w -s' -o bot .
echo "Done!"