# NVIDIA Device Plugin Checker

A simple utility to check the status of NVIDIA GPUs through the Kubernetes device plugin socket.

## Prerequisites

- Kubernetes cluster with NVIDIA GPUs
- NVIDIA device plugin installed and running
- Access to `/var/lib/kubelet/device-plugins/nvidia-gpu.sock`

## Usage

1. Choose the appropriate binary for your system from the `bin` directory:
   - Linux AMD64: `nvdp-checker-linux-amd64`
   - Linux ARM64: `nvdp-checker-linux-arm64`
   - macOS AMD64: `nvdp-checker-darwin-amd64`
   - macOS ARM64: `nvdp-checker-darwin-arm64`

2. Copy the binary to the target machine:
   ```bash
   scp bin/nvdp-checker-linux-amd64 user@remote-machine:~/nvdp-checker
   ```

3. Make the binary executable:
   ```bash
   chmod +x nvdp-checker
   ```

4. Run the checker:
   ```bash
   sudo ./nvdp-checker
   ```

   Note: sudo is required to access the device plugin socket.

## Output

The program will display:
- Total number of GPUs
- Number of healthy GPUs
- GPU distribution across NUMA nodes

Press Ctrl+C to exit the program.

## Troubleshooting

1. "Socket not found":
   - Verify that the NVIDIA device plugin is running
   - Check if the socket exists at `/var/lib/kubelet/device-plugins/nvidia-gpu.sock`
   - Ensure you have the necessary permissions

2. Connection errors:
   - Verify that the NVIDIA device plugin is running
   - Check Kubernetes logs for any device plugin issues 