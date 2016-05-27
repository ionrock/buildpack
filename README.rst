===========
 Buildpack
===========

Buildpack is a command line tool for managing and implementing simple
build processes. It is similar to `Heroku
buildpacks<https://devcenter.heroku.com/articles/buildpacks>`_ with
the some differences being:

 - programming language environment is not detected
 - make is used to support complex pipelines
 - the api is based make targets rather than scripts

The scope of buidpack is not to support building a specific
deliverable for a specific system based on a wide array of code
respository. Instead, buildpack is meant to make it easy to run tests,
build artifacts or run services for a limited set of known
repositories.


Repository API
==============

By default, buildpack uses the Makefile in a repository for its
actions. The essential targets are:

bootstrap:
  This will build the environment in a workspace however you intend in
  order to run tests or processes.

test:
  The bootstrap is run and then the test target run to execute any
  tests.

run:
  The bootstrap is run and then the run target is called to run your
  code.

If your Makefile has these targets in your repository, buildpack will
happily call these commands to perform a build.


Custom Buildpacks
=================

When a project doesn't contain a Makefile a buildpack can be
created. The buildpack must contain a Makefile and can have any other
set of files that will be copied into the checked out repo as though
it is part of it. You can implement the required targets however you
see fit.
