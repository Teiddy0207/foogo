param(
    [string]$ProtoRoot = "api/proto",
    [string]$GoOut = "gen/go",
    [string]$PythonOut = "Detect-service"
)

$ErrorActionPreference = "Stop"

if (!(Test-Path $GoOut)) {
    New-Item -ItemType Directory -Path $GoOut -Force | Out-Null
}

protoc -I $ProtoRoot `
  --go_out=$GoOut --go_opt=paths=source_relative `
  --go-grpc_out=$GoOut --go-grpc_opt=paths=source_relative `
  "$ProtoRoot/detect/v1/detect.proto"

python -m grpc_tools.protoc -I $ProtoRoot `
  --python_out=$PythonOut --grpc_python_out=$PythonOut `
  "$ProtoRoot/detect/v1/detect.proto"

Write-Host "Proto generated successfully for Go and Python."
