# Service-Template

This repository can be used as a template for creating go microservices. It
defines the project structure that has to be followed, provides some
configuration files, and explains how to get started.

**CONTENTS**
1. [Init a new project](#init-a-new-project)
2. [Microservices structure](#microservice-structure)
3. [Create and run a new microservice](#create-and-run-a-new-microservice)
4. [CI/CD](#cicd)
5. [Style guide](#style-guide)


## Init a new project

To create a new microservice project we will use `gonew`. First install the tool
using:
```bash
go install golang.org/x/tools/cmd/gonew@latest
```

To create a new project we need to provide two arguments:
* the path to the template that we are using
* the module name of the project that we are creating

```bash
gonew github.com/eventscompass/service-template github.com/eventscompass/<service-name>
```

Note that the `.circleci/config.yml` and the `Dockerfile` files are already
configured, but you still need to replace `<service-name>` with the real service
name. Run `go mod vendor` to download the dependencies and then push to `main`:

```bash
git add .
git commit -m "first commit"
git push -u origin main
```

## Microservice structure

For our go microservices we use a very simple folder structure conforming to the
following rules:
* all source code lives inside the `src/` folder
* source code is separated into business logic (inside `src/internal/`) and api
  logic (inside `src/`)
* clients for the api are inside the `pkg/client/` folder
* dependencies live inside the `vendor/` folder and are committed together with
  the src code

```bash
.circleci/
├─ config.yml
pkg/
├─ client/        # clients for the service apis
src/
├─ internal/
├─ config.go
├─ events.go      # event subscriptions
├─ grpc.go        # grpc api of the service
├─ main.go
├─ rest.go        # rest api of the service
vendor/
.golangci.yml
Dockerfile
README.md
```

The code for defining and starting the microservice should be inside
`src/main.go`. The api endpoints and handlers live inside `src/rest.go` in case
the service exposes a rest api, and inside `src/grpc.go` in case the service
exposes a grpc api. It is possible for a service to expose both a rest and a
grpc api, although this is not advisable. If the service is subscribed for
events, then the event handler logic should be placed inside `src/events.go`.
Integration tests have to be provided for all api endpoints as well as for the
event subscriptions. Client code for the api(s) of the service has to be
provided inside the `pkg/client/` folder.

In addition every service must have a `Dockerfile`, a `.golangci.yml`
configuration file for running linter checks, and a `.circleci/config.yml`
configuration file for running the CI/CD pipelines.

Every microservice must also have a `README.md` file that explains the purpose
of the service, documents the endpoints that it exposes, and lists the
environment variables that are used for configuring it.


## Create and run a new microservice

Every microservice has to implement the
[service.CloudService](https://github.com/eventscompass/service-framework/blob/bea214ff9294fac13989ed8611affebfb77996b4/service/cloudservice.go#L12)
interface.

```go
type ServiceName struct {
    // ...
}

var _ service.CloudService = (*ServiceName)(nil)
```

To start the service simply call
[service.Start](https://github.com/eventscompass/service-framework/blob/bea214ff9294fac13989ed8611affebfb77996b4/service/start.go#L27).

```go
func main() {
    service.Start(&ServiceName{})
}
```

To build a docker image of the service use the existing `Dockerfile`. Replace
`<service-name>` with the real service name and run:

```bash
docker build -t <service-name> .
```

The rest api should be available on port `:8080` and the grpc api should be
available on port `:8081`.


## CI/CD

For CI/CD we use CircleCI. The configuration for running the pipelines is
already defined inside `.circleci/config.yml`, but you still need to replace
`<service-name>` with the real service name.

In addition, in order for the pipelines to run, a new project has to be created
on the
[CircleCI dashboard](https://app.circleci.com/projects/organizations/circleci%2FRGExmu1KKSYDZZz3vWH7qN/).
Click on `Create New Project`, and generate a new SSH key pair with

```bash
ssh-keygen -t ed25519 -f ~/.ssh/circleci-<service-name> -C email@example.com
```

Then copy the private key to the CircleCI dashboard and the public key to Github
(search for ***Settings &rarr; Deploy Keys***).

Once the project is created you can create embed code for a status badge. Go to
***Project Settings &rarr; Status Badges*** and click on `Add API Token`. This
will create a status token which can be used to query CircleCI for the status of
a job. Then copy the generated embed code and paste it in the `README.md` file.


## Style Guide

When writing go code some style conventions have to be followed in order to keep
the code base idiomatic, consistent and manageable.

First and foremost, make sure you run `gofmt` using the default arguments. Code
that is not properly formatted will not pass the build checks. While `gofmt`
fixes the majority of the mechanical style issues, there are a lot of
non-mechanical style points which should be followed and are addressed in the
following style guides:
* [Effective Go](https://go.dev/doc/effective_go)
* Google's [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
* Google's [Go Style guide](https://google.github.io/styleguide/go/)
* Uber's [Go Style guide](https://github.com/uber-go/guide/blob/master/style.md)

In addition, there are a few other specific rules:
* Try to limit line length to 80 characters. This makes it easier to open
  multiple windows on a single monitor.
* Exported functions are allowed to return only errors defined in
  [service-framework/service/errors.go](https://github.com/eventscompass/service-framework/blob/main/service/errors.go).
  Every exported function must document the errors it could possibly return.
* Try to structure the definitions inside a file as follows:
  * Put exported constants and variables at the top of the file.
  * Put exported types and functions after that. Constructors should be placed
    directly after the type.
  * Put unexported types and functions after that.
  * Put unexported constants and variables at the end of the file.
* In most cases interfaces should be defined in the package that uses them. In
  other languages an interface says "this is the functionality that I am
  providing". In Go, an interface says "this is the functionality that I need".
  If you want a contract for the provided functionality, then write tests.
