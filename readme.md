# Web Log

Learning golang.

## Goal

A single binary (except the Postgres database) that can render blog posts via cacheable server-side templates (to allow for search engines) as well as create and edit them via an interactive single-page client app. Normally, you'd serve both client and API via a web-proxy, etc, etc, but I want to see how close I can get to the JVM world's `uberjar` concept.

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
  }
}
```

By "sparse" I mean that if you just want to change the web port, you can create a file with only that setting in it:

```javascript
{ "web": { "port" : "3001" } }
```

And the app will merge that into the defaults. No need to copy the defaults and tweak. Just change the ones you need to change and omit the rest.

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

    $ createuser blogsvc -P
    $ createdb blogdb -O blogsvc

Use `wanheda` as the password (if you want to go with the application defaults). The app itself will take care of populating all the tables.

There's also a script (for Docker):

    ./script/pg

for starting, stopping, cleaning (etc) a Docker postgres image, if that works better.

**Delete database**

Drop the database, then the user:

    $ dropdb blogdb
    $ dropuser blogsvc

If something is holding open a connection, just restart the database itself and try again:

    $ brew services restart postgres

Homebrew's `brew services` stuff is actually kind handy.

    $ brew install postgres
    $ brew services start postgres
    $ brew services list

This is especially useful if you don't want to or can't use Docker on your Mac. (I have an old Macbook Air, for instance, that won't run it.)
