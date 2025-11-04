# coordinator-go

**coordinator-go** is a lightweight Saga Orchestrator built on top of the
[cmd-stream-go](https://github.com/cmd-stream/cmd-stream-go) library. It
provides a simple, pattern-based way to coordinate distributed operations
without introducing complex infrastructure.

> ⚠️ Experimental Project  
> This project is experimental and not ready for production use.

## Why coordinator-go?

1. **High performance, low footprint.**
   Designed for fast, reliable workflow execution with minimal resource usage.
2. **Simple services.**
   With `coordinator-go`, services remain simple — they behave like ordinary
   services interacting with end-users. No workflow-specific infrastructure or
   APIs are required.
3. **SaaS/Managed option.**
   A potential SaaS version could fully abstract distributed-system complexity,
   offering a ready-to-use solution for Saga orchestration.

These features aim to significantly reduce development and operational
overhead.

## Key Concepts & Benefits

1. **No new DSL or framework to learn.**
   Workflows are implemented using standard `cmd-stream-go` Commands.
2. **Backed by a distributed log, not a database.**
   Logs are efficient and fast for workflow progress tracking.
3. **Flexible communication layer.**
   Works with cmd-stream, HTTP, gRPC, and other transports.
4. **Persistent Command outcomes.**
   During execution, a Command stores outcomes — each representing the result
   of a successful step in the workflow.
5. **No single point of failure.**
   Run one or more `coordinator-go` instances (optionally behind a load
   balancer), each handling workflows for a distinct set of services.
6. **Fault-tolerant & highly available.**
   Platforms like Kubernetes can restart and scale instances automatically.
7. **Idempotency required.**
   Services must be idempotent — workflow steps may be repeated after restarts
   or failures.

## How It Works

```
Command  -> Workflow

Receiver -> Service Gateway
```

In `coordinator-go`, a Command represents a workflow. The Receiver acts as a
service gateway responsible for sending requests to services and receiving
their responses.

### Command Persistence

Before executing a Command is persisted to the distributed storage. This
guarantees it won’t be lost in the case of an unexpected shutdown.

### On Startup

When started, `coordinator-go`:

1. Loads and executes all uncompleted Commands from persistent storage.
2. Enforces a concurrency limit (`Options.MaxCmds`) — new Commands exceeding
   the limit are ignored until capacity frees up.
3. Begins accepting new Commands from clients.

### Command Errors

`Cmd.Exec()` may return:

- `ErrCmdDelayed` — If one of the services is temporarily unavailable.
- `ErrCmdBlocked` — If one of the services is unavailable for an extended
  period (for example, when a circuit breaker is open).
- Other errors — Will cause the coordinator to close the client connection
  (standard `cmd-stream-go` behavior).

### Delayed Commands

If a service becomes unresponsive, the Command should return `ErrCmdDelayed`,
because the outcome is unknown — the operation may have succeeded or failed.

Delayed Commands stay in memory and are retried according to
`Options.SlowRetryInterval`.

### Suspending the Coordinator

If a service remains unavailable long enough to trigger a circuit breaker,
`Cmd.Exec()` should return `ErrCmdBlocked`, causing `coordinator-go` to stop
accepting new Commands.

While suspended:

- Existing Commands continue retrying.
- Retries use `Options.SlowRetryInterval`.
- After the first successful retry, interval switches to
  `Options.FastRetryInterval` (default: `0`).
- Both intervals use jitter to avoid retry storms.

### Intermediate Outcomes

Intermediate Command outcomes serve two purposes:

- Resume execution from the last successful step during retries.
- Rebuild service state by replaying outcomes if the service loses all of its
  data.
