![build](https://github.com/uhppoted/uhppoted-api/workflows/build/badge.svg)

# uhppoted-api

*THIS INCREASINGLY INACCURATELY NAMED REPOSITORY HAS BEEN ARCHIVED AS OF v0.7.0*

_The code has been migrated to [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) and all future development
will take place in the repository. It will remain in place for the next round of releases until superseded 
throughout by _uhppoted-lib._

Higher level Go API for access control systems based on the *UHPPOTE UT0311-L0x* TCP/IP Wiegand controller boards. 

This module:
- Abstracts the device level functionality provided by *uhppote-core* to provide the functionality 
common to *uhppote-cli*, *uhppoted-rest* and *uhppoted-mqtt*
- Provides wrapper functions that support for invoking functionality across multiple devices.
- Implements an ACL (access control list) API to unify access control across multiple controllers.

## Releases

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
| v0.7.0    | Added support for time profiles from the extended API                                     |
| v0.6.12   | Additional validation of bind, broadcast and listen ports when loading configuration      |
| v0.6.10   | Adds configuration options for initial release of `uhppoted-app-wild-apricot`             |
| v0.6.8    | Maintenance release for version compatibility with `uhppote-core` `v0.6.8`                |
| v0.6.7    | Implements `record-special-events` for enabling/disabling door events                     |
| v0.6.5    | Maintenance release for version compatibility with NodeRED module                         |
| v0.6.4    | Added support for uhppoted-app-sheets                                                     |
| v0.6.3    | Added support for `uhppoted-mqtt` ACL API                                                 |
| v0.6.2    | Added support for `uhppoted-rest` ACL API                                                 |
| v0.6.1    | Added support for `uhppote-cli` ACL functions                                             |
| v0.6.0    | Added support for `uhppoted-acl-s3` ACL functions                                         |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules* |

### Building from source

Assuming you have `Go` and `make` installed:

```
git clone https://github.com/uhppoted/uhppoted-api.git
cd uhppoted-api
make build
```

If you prefer not to use `make`:
```
git clone https://github.com/uhppoted/uhppoted-api.git
cd uhppoted-api
mkdir bin
go build -o bin ./...
```

#### Dependencies

| *Dependency*                                             | *Description*                                          |
| -------------------------------------------------------- | ------------------------------------------------------ |
| [uhppote-core](https://github.com/uhppoted/uhppote-core) | Device level API implementation                        |
| golang.org/x/sys/windows                                 | Support for Windows services                           |
| golang.org/x/lint/golint                                 | Additional *lint* check for release builds             |

## uhppote-api


