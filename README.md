# uhppoted-api

High level Go API for access control systems based on the *UHPPOTE UT0311-L0x* TCP/IP Wiegand controller boards. 

This module
abstracts the device level functionality provided by *uhppote-core* to provide the functionality common to *uhppote-cli*, 
*uhppoted-rest* and *uhppoted-mqtt*.

## Releases

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
| v0.6.3    | Added support for `uhppoted-mqtt` ACL API                                                 |
| v0.6.2    | Added support for `uhppoted-rest` ACL API                                                 |
| v0.6.1    | Added support for `uhppote-cli` ACL functions                                             |
| v0.6.0    | Added support for `uhppoted-acl-s3` ACL functions                                         |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules* |

### Building from source

#### Dependencies

| *Dependency*                                             | *Description*                                          |
| -------------------------------------------------------- | ------------------------------------------------------ |
| [uhppote-core][https://github.com/uhppoted/uhppote-core] | Device level API implementation                        |
| golang.org/x/sys/windows                                 | Support for Windows services                           |
| golang.org/x/lint/golint                                 | Additional *lint* check for release builds             |

## uhppote-api


