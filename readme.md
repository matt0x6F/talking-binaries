# Talking binaries

Creating a viable plugin system in Go is challenging. Some avenues (and architectural examples) I considered are:

- [go-plugin](https://github.com/hashicorp/go-plugin)
- [Go plugins](https://pkg.go.dev/plugin)
- [KubeCtl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)

This repository is my effort to learn what it takes to make an effective plugin system for Go applications. I intend to
take this experience and either craft a better system or pick a strategy more similar to the aforementioned projects.

## Example

```shell
$: ./dist/program
[main] 2022/01/24 01:57:30 Got config: {Plugins:[{Ref:echo Name:dist/plugin-echo Path:} {Ref:shell Name:dist/plugin-shell Path:}] Steps:[{Name:say_hi Plugin:echo Config:map[say:hello world] Parallel:false Next:<nil> Prev:<nil> Status:0} {Name:say_hi_via_bash Plugin:shell Config:map[commands:[echo -n "hello world" echo -n "hello world 2"]] Parallel:false Next:<nil> Prev:<nil> Status:0}]}
[main] 2022/01/24 01:57:30 Found plugin at location: dist/plugin-echo
[main] 2022/01/24 01:57:30 Found plugin at location: dist/plugin-shell
[main] 2022/01/24 01:57:30 Built plugin store: map[echo:/home/matt/Projects/matt/talking-binaries/dist/plugin-echo shell:/home/matt/Projects/matt/talking-binaries/dist/plugin-shell]
[main] 2022/01/24 01:57:30 Executing step: say_hi
[main] 2022/01/24 01:57:30 Sending data: {"Config":{"say":"hello world"},"Parallel":false}
[main] 2022/01/24 01:57:30 Binary execution time took 3.24336ms
[main] 2022/01/24 01:57:30 Error encountered while unmarshaling command output: invalid character 'C' looking for beginning of value
User Output:
None
[main] 2022/01/24 01:57:30 Step execution time took 3.429525ms
[main] 2022/01/24 01:57:30 Executing step: say_hi_via_bash
[main] 2022/01/24 01:57:30 Sending data: {"Config":{"commands":["echo -n \"hello world\"","echo -n \"hello world 2\""]},"Parallel":false}
[main] 2022/01/24 01:57:30 Binary execution time took 8.258791ms
[plugin-shell] 2022/01/24 01:57:30 Beginning execution
[plugin-shell] 2022/01/24 01:57:30 Config received: {Config:{Commands:[echo -n "hello world" echo -n "hello world 2"] Parallel:false}}
[plugin-shell] 2022/01/24 01:57:30 Preparing to execute 2 commands
[plugin-shell] 2022/01/24 01:57:30 Executing command: /usr/bin/echo -n "hello world"
[plugin-shell] 2022/01/24 01:57:30 Executing command: /usr/bin/echo -n "hello world 2"
User Output:
"hello world"
"hello world 2"
[main] 2022/01/24 01:57:30 Step execution time took 8.683066ms
```

## Analysis

I recorded two metrics: binary execution time and step execution time, where binary execution time is a subset of step
execution time. The echo plugin is the most lightweight.

| name  | binary execution time | step execution time |
|:------|:---------------------:|:-------------------:|
| echo  |       3.24336ms       |     3.429525ms      |
| shell |      8.258791ms       |     8.683066ms      |

## Features

Features implemented so far.

- Recursive search for configuration
- Ordered execution of steps
- JSON-based protocol for transmitting state
- Plugin discovery - Plugins can be discovered in two ways.
  - Relative locations (eg: `dist/plugin-shell`)
  - `PATH`-based location (eg: `plugin-shell` in `/usr/local/bin`)
- Log mirroring - Plugins use their own stdin/stdout/stderr, however, it is synced back to the main process.

## Plans

- The step queue should be able to smartly run some steps in parallel
- Plugins should be made aware that they're running in parallel (two instances of the same binary)
- Global variable storage
  - Receive updates/additions from plugins
  - Variable substitution in the `Config` block

## Shortcomings

Some shortcomings of the current architecture.

- **[Critical]** If a plugin terminates early it's logs are never synced back to the main process and thus never shown to the user
- **[Minor]** Limited to request/response style transactions