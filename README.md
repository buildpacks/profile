# `docker.io/buildpacksio/profile`

The Profile Buildpack is a Cloud Native Buildpack that implements the behavior of `.profile` scripts. A `.profile` script is executed prior to the start of an application.

## Behavior

This buildpack will participate all the following conditions are met

* `bash` is available in the container
* a `.profile` script file exists at the root of the application

The buildpack will do the following:

* contribute an `exec.d` script which will at runtime source the `.profile` script and ensure environment variables set by the script are available to your application

## License

This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
