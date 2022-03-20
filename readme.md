# WEBL 3

This project is a protean project that changes whenever I want to try some new way to do a blog project.

I’m in the process of adding a Desktop Mac App to use as a blog composer and site admin utility. When it’s done, I’ll delete the React based web app that does the same thing.

- [console](./console) -- A mac app for creating and editing posts, and setting the site tile and description.
- [server](./server) -- A Go server (postgres, templates, GraphQL).
- [database](./database) -- Database scripts, etc, for postgres.

(Go to the [server](./server) to see the old readme with its original motivations.)

I’ve just refactored the source repo so you can’t really build this without already knowing how to put it together. If you see a Makefile in the top level directory, it means I’ve refactored the build system.

Basically, you build the `webl` Go app (binary), and then point it at the resources in the `WeblServer` directory and at wherever you build the `WeblAdmin` app and things should work.

Don’t do it, though. It’s all running on a server, but I want this project to be in flux until the Mac app is done.

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
