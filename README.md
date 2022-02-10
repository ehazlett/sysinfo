# sysinfo
This is a simple HTTP application that returns system info.

# Trace Support
There is also simple OpenTelemetry tracing support via the `-t` flag.
Configure that to send to any OpenTelemetry collector and you should
see traces for both the `/` and `/version` endpoints.
