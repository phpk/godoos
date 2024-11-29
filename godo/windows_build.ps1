param (
    [string]$GOOS = $(go env GOOS),
    [string]$GOARCH = $(go env GOARCH)
)

$BINARY_NAME = "godo"
$DEPS_DIR = "deps"
$PLATFORM_DIR = Join-Path -Path $DEPS_DIR -ChildPath $GOOS
$ZIP_FILE = Join-Path -Path $DEPS_DIR -ChildPath "${GOOS}.zip"

if (Test-Path -Path $PLATFORM_DIR) {
    Write-Output "Compressing $PLATFORM_DIR..."
    Compress-Archive -Path $PLATFORM_DIR -DestinationPath $ZIP_FILE
    Write-Output "Compression completed."
} else {
    Write-Output "Directory $PLATFORM_DIR does not exist. Skipping compression."
}

Write-Output "Building $BINARY_NAME for $GOOS/$GOARCH..."
go build -ldflags="-s -w" -o $BINARY_NAME

Write-Output "Build completed."