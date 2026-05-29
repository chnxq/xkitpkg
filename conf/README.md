# Protocol Buffers (Protobuf) Configuration

This directory contains the Protocol Buffers (protobuf) definitions for the xkitpkg configuration system and the generated Go code.

## Directory Structure

- `protos/v1/` - Contains the `.proto` definition files for all configuration messages
- `v1/` - Contains the generated Go code from the protobuf definitions
- `buf.gen.yaml` - Configuration for the Buf code generator
- `buf.work.yaml` - Buf workspace configuration
- `buf.yaml` - Buf module configuration

## Installing Buf

Buf is a tool for working with Protocol Buffer schemas. To install Buf, run:

```bash
# Install buf using go install
go install github.com/bufbuild/buf/cmd/buf@latest

# Or install using curl on Linux/macOS
curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v$(curl -s https://api.github.com/repos/bufbuild/buf/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')/buf-$(uname -s)-$(uname -m).tar.gz" | \
    tar -xvzf - -C /tmp && \
    sudo mv /tmp/buf /usr/local/bin

# On Windows with PowerShell
$version = $(Invoke-RestMethod -Uri "https://api.github.com/repos/bufbuild/buf/releases/latest").tag_name
$downloadUrl = "https://github.com/bufbuild/buf/releases/download/$version/buf-$(Get-ComputerInfo).exe"
Invoke-WebRequest -Uri $downloadUrl -OutFile "$env:TEMP\buf.exe"
Move-Item "$env:TEMP\buf.exe" "C:\Program Files\buf.exe"
```

## Generate Go Code

To regenerate the Go code from the protobuf definitions, run the following command from the `conf` directory:

```bash
# Navigate to the conf directory
cd D:\GoProjects\chnxq\xkitpkg\conf

# Generate Go code from protobuf definitions
buf generate
```

Alternatively, you can run the generation from anywhere within the project:

```bash
# Generate code using the buf configuration
buf generate --template buf.gen.yaml
```

## Adding New Protobuf Definitions

When adding new `.proto` files:

1. Add your `.proto` file to the `protos/v1/` directory
2. Make sure to include the proper Go package option in your proto file:
   ```protobuf
   option go_package = "github.com/chnxq/xkitpkg/conf/v1;conf";
   ```
3. Run `buf generate` to generate the corresponding Go code
4. Import and use the generated types in your Go code

## Buf Commands

Common buf commands used in this project:

```bash
# Validate protobuf files
buf lint

# Check for breaking changes
buf breaking --against '.git/main'

# Generate code
buf generate

# Format protobuf files
buf format -w
```

## Dependencies

This project uses external protobuf definitions from:

- Google APIs (`buf.build/googleapis/googleapis`)
- Protobuf Validation (`buf.build/envoyproxy/protoc-gen-validate`)
- xkit APIs (`buf.build/xkit/apis`)
- Gnostic (`buf.build/gnostic/gnostic`)
- Gogo Protobuf (`buf.build/gogo/protobuf`)

These dependencies are managed through the `buf.yaml` file in the `protos` directory.