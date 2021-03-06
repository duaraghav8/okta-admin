# okta-admin
This is a Commandline application to perform administrative tasks in [Okta](https://www.okta.com).

It does not allow you to make requests to arbitrary endpoints of the Okta API. Rather, it is designed to speed up frequently performed administrative tasks like creating users, adding them to groups, listing resources and so forth.

## Installation
Download the pre-compiled binary for your platform from the [Releases](https://github.com/duaraghav8/okta-admin/releases) page or use `make` to compile the code yourself (see [Developing](#developing)). Move the binary to an appropriate directory on your system to make it available on the `PATH`.

## Usage

Use `-help` to display the list of commands available and to get information about specific commands.

### Examples
1. Create a new user in your Okta Organization
```bash
okta-admin create-user \
    -org-url https://hogwarts.okta.co.uk \
    -api-token xxxxx \
    -email harry.potter@hogwarts.co.uk \
    -fname Harry -lname Potter \
    -team Seekers
```
The above command creates a new user in an organization and assigns them a Team. It also highlights how you can specify the organization URL and Okta API Token via commandline arguments, although **this is not the recommended way to supply credentials**.

2. Add a member to groups in the organization
```bash
export OKTA_ORG_URL="https://hogwarts.okta.com/"
export OKTA_API_TOKEN="xxxxx"

okta-admin assign-groups -email albus.dumbledore@hogwarts.co.uk -groups TheOrder

okta-admin assign-groups \
    -email draco.malfoy@hogwarts.co.uk \
    -groups "Slytherin, pure-blood, rich_kids"

okta-admin assign-groups \
    -email newt.scamander@hogwarts.co.uk \
    -groups hogwarts-alumni,MinistryOfMagic
```
These commands demonstrate the different ways in which you can specify `groups` to assign to a member. Any option capable of accepting multiple values can be given a comma-separated list of them. Notice how the organization credentials this time are passed via environment variables. This is the recommended way to work with Okta Admin, especially when running the tool in automation.

3. List Groups present in the organization
```bash
# Load credentials from an environment file
source ~/.okta/creds.env

# List names of all groups
okta-admin list-groups

# Get info about select groups
okta-admin list-groups -groups azkaban,durmstrang -detailed
```

## Developing
This project uses [Go Modules](https://blog.golang.org/using-go-modules) for dependency management. You must have at least Go version 1.11 installed on your system to develop this project.

### Test
```bash
make fmt
make test
```

### Dependencies
Before sending a PR or building with your local changes, use `mod` to clean up dependencies.
```bash
# Install required dependencies & remove irrelevant ones
go mod tidy
# Add dependencies to vendor/
go mod vendor
```

### Build
To build the binary for your platform(s), run the following make tasks:
```
# Linux
make linux/amd64

# Darwin
make darwin/amd64

# Windows
make windows/amd64
```

## License
This code is licensed under the MPLv2 license.
