## v0.7.x

### IN PROGRESS

- [ ] Rename SetTimeProfile to PutTimeProfile for consistency
- [ ] Rearchitecture UHPPOTED as an interface+implementation

- [x] Update GetCard(s) for time profiles
- [x] Update PutCard(s) for time profiles
- [x] Update ACL for time profiles
- [x] Implement set/get/clear-time-profile

## TODO

1. Rework healthcheck to remove need for IUHPPOTE::DeviceList
2. Rework healthcheck to remove need for IUHPPOTE::ListenAddr
3. GetDevices: rename DeviceSummary.Address to IpAddress and use Address for IP+Port


### uhppoted-api

- [ ] websocket + GraphQL (?)
- [ ] IFTTT
- [ ] Braid (?)
- [ ] MacOS launchd socket handoff
- [ ] Linux systemd socket handoff
- [ ] conf file decoder: JSON
- [ ] Rework plist encoder
- [ ] move ACL and events to separate API's
- [ ] Make events consistent across everything
- [ ] Rework uhppoted-xxx Run, etc to use [method expressions](https://talks.golang.org/2012/10things.slide#9)
- [ ] system API (for health-check, watchdog, configuration, etc)
- [ ] Parallel-ize health-check 

### Documentation

- [ ] godoc
- [ ] build documentation

### Other

1. github project page
2. Integration tests
3. EventLogger 
    - MacOS: use [system logging](https://developer.apple.com/documentation/os/logging)
    - Windows: event logging
4. TLA+/Alloy models:
    - watchdog/health-check
    - concurrent connections
    - HOTP counter update
    - key-value stores
    - event buffer logic
5. Update file watchers to fsnotify when that is merged into the standard library (1.4 ?)
    - https://github.com/golang/go/issues/4068
6. go-fuzz
