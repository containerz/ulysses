# Changelog

#### Version 0.3.5.2

Fixes:

- Fixed a bug causing Docker containers to fail with empty launch command

#### Version 0.3.5.1

Improvements:

- Application launch command is no more required. This can be useful for Docker containers with default CMDs and ENTRYPOINTs.

#### Version 0.3.5.0

Improvements:

- Added `args` and `env` to application, allowing to pass arbitrary arguments (non key-value pairs) and environmental variables.

#### Version 0.3.4.0

**This version introduces new required --storage flag, see readme for explanation**

Improvements:

- Persist bootstrapped Cassandra contact points to persistent storage to handle stack-deploy failovers

#### Version 0.3.3.1

Improvements/Fixes:

- Reduced Marathon healthchecks interval for faster deployment    
- Output stack context after deployment is finished

#### Version 0.3.3.0

New Features:

- Docker support for Marathon tasks

#### Version 0.3.2.0

Improvements:

- Run-once tasks now respect constraints

#### Version 0.3.1.1

Fixes:

- Fixed broken list command in developer mode

#### Version 0.3.1.0

New Features:

- Developer mode

#### Version 0.3.0.0

New Features:

- Added run-once tasks

#### Version 0.2.1.0

New Features:

- Added possibility to pass arbitrary variables to run command
- Added possibility to skip applications while running a stack
- Added possibility to pass global variables during server start

Improvements:

- Various bug fixes
- Exposed more stack variables for datastax-enterprise-mesos
- Structure context with variable precendence (global, arbitrary, stack)
- All Mesos and Marathon tasks run are now labeled with stack and zone names

#### Version 0.2.0.1

New Features:

- Changed Application.Tasks to be ordered map

Improvements:

- Test coverage
- Readme for task runners
- Fixed some data races

#### Version 0.1.0.0

Basic functionality