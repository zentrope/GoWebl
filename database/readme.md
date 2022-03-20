# Webl Data

These are the old database notes sectioned off into a separate project.


## The Makefile

This project uses a `Makefile` as a script runner, not really as a software builder. It’s easy enough to read, but you can also figure out its commands with:

    $ make

or

    $ make help

which both produce:

    $ make
    db-clean                  Delete the local dev database
    db-init                   Create a local dev database with default creds
    db-schema                 Try to use the webl server to load the schema
    help                      Display this help listing

If postgres is running and it’s your first time with this code:

    $ make db-init
    $ make db-schema

And you’re ready to go with ... whatever it is you need to do.

## Install Database

Homebrew's `brew services` stuff is handy:

    $ brew install postgres
    $ brew services start postgres
    $ brew services list

This is especially useful if you don't want to or can't use Docker on your Mac.

## Default (dev) database params

These are the default connection parameters for a Postgres instance:

    database: webl_db
    user:     webl_user
    pass:     wanheda
    host:     localhost
    port:     5432

The app [WeblServer](../WeblServer) takes a config file to change all these.

## Create Database

Create a user, then a database (owned by the user):

    $ make db-init

This will also set the user's password. The app itself will take care of populating all the tables.

## Create the database schema

Using the `Makefile`:

    $ make db-schema

This works well enough using `zsh` on my Mac, but may not work for you. Wait, who are you?

To avoid using the `Makefile`:

    $ cd ../WeblServer
    $ go run . -resources resources
    $ ^D  # Don't leave the server running, unless you want to.

Each time you start the server, it’ll inspect a migration table and apply database update scripts accordingly. The database has an initial post, and an initial author, login: `root@example.com`, password is my first.last name.

## Delete database

Drop the database, then the user:

    $ make db-clean

If something is holding open a connection, just restart the database itself and try again:

    $ brew services restart postgres

## Database is Ancient

If your local databases are incompatible because you've updated your server multiple times without ever starting it or reading the fine print about upgrading, you can do the following:

    $ brew services stop postgresql
    $ mv /usr/local/var/postgres ~/.Trash
    $ initdb /usr/local/var/postgres
    $ brew services start postgresql
    $ brew services ls

Make sure the process is stopped (and it should be already given the problem), remove the database files, run the `initdb` comment, start up postgresql and then verify that it's running.

After that, run:

    $ make db-init

And hopefully, you're back in business.

## License

Copyright (c) 2017-2022 Keith Irwin

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
