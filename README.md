# peh3

peh (PNDLM Environment Helper) is a tool by [PNDLM](https://pndlm.com) for building simple, reliable, reusable software development environments based on Docker Stack.

Philosophical goals—
* Minimize pre-requisite installation requirements on dev machines
* Limit custom tooling and abstraction; directly leverage great off-the-shelf tooling wherever possible
* Remain dev/deploy platform agnostic

Components are—
* A Golang CLI library (https://github.com/pndlm/peh/peh3) with common various docker/npm/dev-ops functions, commands and shortcuts
* A `project-template` used to kick off new software project git repositories that are based on peh
	* Project-specific CLI utility— the template contains a utility app in `/peh/src` that ingests the peh3 library and is customizable to whatever the project needs
* A global helper utility that allows one to type `peh` at the command-line within any peh project and run its specific CLI utility

## Setup

Install the base pre-requisites—
* Docker Desktop — https://docs.docker.com/desktop/install
	* For Linux, strongly recommend following the [Digital Ocean](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-22-04) instructions for your distribution
* Go — https://go.dev/dl

Install the global `peh` helper on your development machine—

```bash
GOPROXY=direct go install github.com/pndlm/peh/helper/peh@latest
```

After a shell restart or reboot, you should now be able to run `peh` at your command-line to see a list of available commands.

## Adding peh3 to an existing peh2 project

* Copy a fresh `peh` directory into the project
```
cd projectdir
npx degit pndlm/peh/project-template/peh peh --force
```
* Add `use ./peh/src` to `go.work` file as necessary
* Grep and replace all instances of `pndlm.com/example-project` with a customized name e.g. `client.com/appname/peh`
* Install dependencies
```
cd peh/src
go get github.com/pndlm/peh/peh3
go work sync
```

## Developing on peh3 library

To develop and test changes to peh3 library alongside of a consuming project, check out this repository next to your active project and review/uncomment the `replace` directive at the bottom of your project's `/peh/src/go.mod` file.

## TODO

* New project creation instructions
```bash
# don't forget to install the latest version of the library....
GOPROXY=direct go get github.com/pndlm/peh/peh3
```
* Add docker-compose.yaml and README to `project-template`
* Test fully new installation
* Make creating a new installation a command on `peh` helper itself ?
