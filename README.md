# Service Check Application 🚀

**Note: This project is currently under development.🔧**

## Motivation 📌

The motivation for this project arose from the common challenge that many applications do not automatically start or require additional configuration after a PC restart, such as `podman`. The objective of this project is to provide a service that ensures processes remain active, regardless of the provider or application.

## Overview 📋

The Service Check application is a Golang-based tool that provides command-line interface (CLI) utilities and a remote procedure call (RPC) server for managing and monitoring service availability, alongside maintaining detailed logs for control and traceability purposes.
## Components 🛠️

### Command-Line Interface (CLI) 💻

The CLI utilities are located in the `cmd` directory. These tools facilitate various operations concerning service management. The following commands are available:

- **Create**: Initializes a service setup based on a configuration provided in `example/individual-service.json`.
- **List**: Retrieves a list of services with their current status.
- **Logs**: Fetches logs related to service operations, aiding in auditing and debugging processes.

### RPC Server 🖥️

The RPC server resides in the `internal/rpc` directory and is implemented using Go's standard library. Its primary role is to check the availability of services in real-time. It offers the following functionalities:

- **Service Monitoring**: Continuously checks and reports the status of configured services.
- **Logging**: Captures detailed logs of all operations, providing insight into service status changes over time.
- **Control Operations**: Allows for remote management of service states through standardized RPC methodologies.

## Future Developments 📈

As the project progresses, additional features and improvements will be incorporated, including more robust error handling, enhanced logging capabilities, and expanded service control options.

## Getting Started 🚀

### Build ⚙️

Use the `Makefile` for building the application:
```bash
make build       # Builds the main service check application
make run         # Builds and runs the test server
```

### Test 🧪

To execute tests, run:
```bash
make test
```

## Contributing 🤝

As this project is still under active development, contributions are encouraged. Please follow standard Go formatting and conventions when contributing code.

## License 📜

This project is provided under an open-source license. Please review the LICENSE file for more details.