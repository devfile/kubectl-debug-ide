# kubectl debug-ide

A `kubectl` plugin to debug Pods with an IDE rather than the CLI.

:warning: This plugin is in its alpha stage and is missing some important feature compared to `kubectl debug`. Running
the IDE in an ephemeral debugging container (rather than copying the target Pod) is not supported. Workloads other than
`Pods` aren't supported either.

## How does `kubectl debug-ide` work?

Running `kubectl debug-ide` is similar to running `kubectl debug` (same flags) but an IDE is started and a git
repository is cloned within the debugging container.

![kubectl debug-ide in action](img/demo.gif)

The command returns the URL of the remote IDE (local IDE support is not implemented yet). Opening the IDE URL in a 
browser allow debugging a process in the target container. Note that the IDE runs in the debugging container, but
"local" debugging is possible as the debugging container and the target container share the same process namespace.

![debugging using an IDE](img/outyet-debug.gif)

## Examples

The following command creates a copy of the Pod `$TARGET_POD` with an extra container using image
`$DEBUGGING_CONTAINER_IMG` where `$GIT_REPO` is cloned and an IDE is started. 

```sh
TARGET_POD="outyet"
DEBUGGING_CONTAINER_IMG="ghcr.io/l0rd/outyet-dev:latest"
TARGET_POD_COPY="outyet-debug"
GIT_REPO="https://github.com/l0rd/outyet.git"

kubectl debug-ide $TARGET_POD \
  --image $DEBUGGING_CONTAINER_IMG \
  --copy-to $TARGET_POD_COPY \
  --share-processes \
  --git-repository $GIT_REPO
```

:mega: The containers in the copy of target Pod share the PID namespace. This is helpful to attach the IDE debugger to
the target process as they run in separate containers.

Delete the `DevWorkspace` Custom resource to stop the debugging session and cleanup the Kubernetes resources created by 
`kubectl debug-id`:

```sh
TARGET_POD="outyet"
DEBUGGING_CONTAINER_IMG="ghcr.io/l0rd/outyet-dev:latest"
TARGET_POD_COPY="outyet-debug"
GIT_REPO="https://github.com/l0rd/outyet.git"

kubectl debug-ide $TARGET_POD \
  --image $DEBUGGING_CONTAINER_IMG \
  --copy-to $TARGET_POD_COPY \
  --share-processes \
  --git-repository $GIT_REPO
```

:mega: `kubectl delete pod` won't work in this case as the DevWorkspace Operator will restart the Pod.

## Requirements

Running `kubectl debug-ide` requires the [DevWorkspace Operator](https://github.com/devfile/devworkspace-operator/tree/main).
`kubectl debug-ide` creates Custom Resources of type `DevWorkspace`.

Building (and currently installing too) requires [Go](https://go.dev/dl/).

## Installation

To install `kubectl debug-ide` clone this git repository and run `go install`:  

```sh
$ go install cmd/kubectl-debug_ide.go
```

## Shell completion

Copy the file `./kubectl_complete-ide` somewhere on `$PATH` and give it executable permissions to enable shell
completion for the plugin.

:mega: kubectl v1.26 or higher is required for shell completion to work for plugins.

## Cleanup

You can "uninstall" this plugin from kubectl by simply removing it from your PATH:

    $ rm  $GOBIN/kubectl-debug_ide

## Acknowledgments

The awesome [sample-cli-plugin](https://github.com/kubernetes/sample-cli-plugin) was used to kick off this plugin.

`kubectl debug-ide` is just a generator of DevWorkspace Custom Resources. The heavy lifting, the Cloud Development
Environment provisioning, is done by the [DevWorkspace Operator](https://github.com/devfile/devworkspace-operator).
