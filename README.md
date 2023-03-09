# platform-app

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/morning-night-guild/platform-app?style=plastic)

## document

- [api](https://github.com/morning-night-guild/platform-app/tree/gh-pages/api)
- [database](https://github.com/morning-night-guild/platform-app/tree/gh-pages/database)
- [proto](https://github.com/morning-night-guild/platform-app/tree/gh-pages/proto)

## directory structure

four layered architecture

```shell
.
├── domain
│   ├── model           // domain model include id
│   ├── value           // value object
│   ├── repository      // interface domain model persistence
│   └── rpc             // interface domain model remote procedure call
├── usecase
│   ├── interactor      // implements port
│   ├── port            // usecase interface
│   └── mock            // for test
├── adapter
│   ├── controller      // core adapter
│   ├── gateway         // core adapter (implements repository)
│   ├── handler         // api adapter
│   ├── external        // api adapter (implements rpc)
│   └── mock            // for test
└── driver
    ├── config
    ├── connect
    ├── cors
    ├── database
    ├── env
    ├── http
    ├── interceptor
    ├── middleware
    ├── newrelic
    ├── router
    └── server
```

## commit message prefix

Create an issue and include the number in the PREFIX when implementing.

| PREFIX           | meaning                                                                                                |
| ---------------- | ------------------------------------------------------------------------------------------------------ |
| **feat(#x)**     | A new feature                                                                                          |
| **fix(#x)**      | A bug fix                                                                                              |
| **docs(#x)**     | Documentation only changes                                                                             |
| **style(#x)**    | Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc) |
| **refactor(#x)** | A code change that neither fixes a bug nor adds a feature                                              |
| **perf(#x)**     | A code change that improves performance                                                                |
| **test(#x)**     | Adding missing or correcting existing tests                                                            |
| **chore(#x)**    | Changes to the build process or auxiliary tools and libraries such as documentation generation         |

[reference](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#commits)

## local development setup

`.env` is required to Makefile.  
`.env` file is prepared to switch the port number in each developer's environment.

```shell
touch .env
make env
```

```shell
make dev
```
