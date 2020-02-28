# uhppoted-api

High level Go API for access control systems based on the *UHPPOTE UT0311-L0x* TCP/IP Wiegand controller boards. This module
abstracts the device level functionality provided by *uhppote-core* to provide the functionality common to *uhppote-cli*, 
*uhppoted-rest* and *uhppoted-mqtt*.

## Releases

### Building from source

#### Dependencies

| *Dependency*                        | *Description*                                          |
| ----------------------------------- | ------------------------------------------------------ |
| com.github/uhppoted/uhppote-core    | Device level API implementation                        |
| golang.org/x/sys/windows            | Support for Windows services                           |
| golang.org/x/lint/golint            | Additional *lint* check for release builds             |

## uhppote-api



