# Web Log

Learning golang.

## Notes

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

I believe started and stopped stuff survices a reboot, if you care.
