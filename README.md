# tmpl
[![Build Status](https://travis-ci.com/nwtgck/tmpl.svg?token=TuxNpqznwwyy7hyJwBVm&branch=develop)](https://travis-ci.com/nwtgck/tmpl)

Template Engine for Repository

## Install

```bash
go get -u github.com/nwtgck/tmpl
```


## Usage

Here is an example to create a project from [docker.tmpl](https://github.com/nwtgck/docker.tmpl).

```bash
tmpl new https://github.com/nwtgck/docker.tmpl.git mydocker
```

Then, you have `mydocker/` directory in pwd.

## Use existing .tmpl

```bash
cd some_existing.tmpl
tmpl fill
```

OR

```bash
tmpl fill some_existing.tmpl
```

## Binaries

Available binaries at [Releases](https://github.com/nwtgck/tmpl/releases).
