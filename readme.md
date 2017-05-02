# Web Log

Learning golang.

## Notes

**Create Database**

```shell
createuser blogsvc -P
createdb blogdb -O blogsvc
```

The app itself will take care of populating all the tables. I'm not sure if "createdb" and "superuser" are actually needed, here.

There's also a script:

    ./script/pg

for starting, stopping, cleaning (etc) a docker postgres image, if that works better.

### Delete database

```shell
dropdb blogdb
dropuser blogsvc
```

If something is holding open a connection, just restart the database itself and try again:

```shell
brew services restart postgres
```

Homebrew's `brew services` stuff is actually kind handy.
