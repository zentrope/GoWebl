# Web Log

Learning Golang.

## Goal

The goal of this learning project is to produce a single binary (except the Postgres database) that can render blog posts via cacheable server-side templates (to allow for search engines) as well as create and edit them via an interactive single-page client app. Normally, you'd serve both client and API via a web-proxy, etc, etc, but I want to see how close I can get to the JVM world's `uberjar` concept.

## Quick Start

Assuming you've got a database going using the defaults (see below):

    $ make build
    $ ./webl

And if you want to point to a configuration file:

    $ ./webl -c /path/to/config.json

The config file should be a sparse version of what you can see in `./resources/config.json`. The default looks like:

```javascript
{
  "storage": {
    "user" : "blogsvc",
    "password": "wanheda",
    "database": "blogdb",
    "host": "localhost",
    "port": "5432"
  },
  "web": {
    "port": "8080"
    "title": "Web Log",
    "base-url": "http://localhost:8080/"
  }
}
```

By "sparse" I mean that if you just want to change the web port, you can create a file with only that setting in it:

```javascript
{ "web": { "port" : "3001" } }
```

And the app will merge that into the defaults. No need to copy the defaults and tweak. Just change the ones you need to change and omit the rest.

## Development

I work on the project using two terminals:

In one, I start up the `yarn` process for working with the React-based admin app.

    $ cd webl/admin
    $ yarn start

In another terminal, I start the server process itself:

    $ make init
    $ go run main.go

And that's it. The admin app will reload when you make changes. You'll have to `^C` the server process, then up arrow and return to cycle it.

## Docker

Running, building, docker-composering: I've just not gotten to it yet. If you want to try this application out, I'd recommend going with stock Homebrew (or whatever) installs of Golang, Postgres and Yarn and leave it at that.

However, I do plan to create some scripts that allow for building the app without having to install any of the dependencies (build or otherwise).

## Database Notes

**Default (dev) database params**

These are the default connection parameters for a Postgres instance:

    database: blogdb
    user:     blogsvc
    pass:     wanheda
    host:     localhost
    port:     5432

This app is set up for a config file to change all these, but I'll document that later, when there's something worth worrying about.

**Create Database**

Create a user, then a database (owned by the user):

    $ make db-init

This will also set the user's password. The app itself will take care of populating all the tables.

**Delete database**

Drop the database, then the user:

    $ make db-clean

If something is holding open a connection, just restart the database itself and try again:

    $ brew services restart postgres

Homebrew's `brew services` stuff is actually kind handy.

    $ brew install postgres
    $ brew services start postgres
    $ brew services list

This is especially useful if you don't want to or can't use Docker on your Mac. (I have an old Macbook Air, for instance, that won't run it.)

## License

Copyright (c) 2017 Keith Irwin

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published
by the Free Software Foundation, either version 3 of the License,
or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see
[http://www.gnu.org/licenses/](http://www.gnu.org/licenses/).
