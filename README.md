# Dynatomic

[![Build Status](https://travis-ci.org/tylfin/dynatomic.svg?branch=master)](https://travis-ci.org/tylfin/dynatomic)
[![Build status](https://ci.appveyor.com/api/projects/status/g58owqmdpumqsmxq/branch/master?svg=true&passingText=Windows%20-%20OK)](https://ci.appveyor.com/project/tylfin/dynatomic/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/tylfin/dynatomic)](https://goreportcard.com/report/github.com/tylfin/dynatomic)
[![GoDoc](https://godoc.org/github.com/tylfin/dynatomic/pkg/dynatomic?status.svg)](https://godoc.org/github.com/tylfin/dynatomic/pkg/dynatomic)
[![Coverage Status](https://coveralls.io/repos/github/tylfin/dynatomic/badge.svg)](https://coveralls.io/github/tylfin/dynatomic)

Dynatomic is a library for using dynamodb as an atomic counter

- [Dynatomic](#dynatomic)
  - [Motivation](#motivation)
  - [Usage](#usage)
  - [Development](#development)
  - [Contributing](#contributing)

## Motivation

The dynatomic was written to use dynamodb as a quick and easy atomic counter.

The package tries to serve two unique use cases:

- Unique, fast real-time writes, e.g. user visits to a page or rate limiting
- Large number of asynchronous writes that need to be eventually consistent, e.g. API usage by a client for billing

## Usage

Basic usage:

```bash
// Initial dynatomic with a batch size of 100, a wait time of a second, your AWS config
// and a function that will notify the user of internal errors
d := New(100, time.Second, config, errHandler)
d.RowChan <- &types.Row{...}
d.RowChan <- &types.Row{...}
d.RowChan <- &types.Row{...}
...
d.Done()
```

Dynamo will update accordingly.

For example if you write the rows:

```bash
Table: MyTable, Key: A, Range: A, Incr: 5
Table: MyTable, Key: A, Range: A, Incr: 5
Table: MyTable, Key: A, Range: A, Incr: 5
Table: MyTable, Key: A, Range: A, Incr: 5
```

Then MyTable Key A, Range A will now show a value of 20

## Development

To copy the repository run:

```bash
go get github.com/tylfin/dynatomic
```

Then you can run the full test suite by doing:

```bash
docker-compose run dynatomic
```

## Contributing

1. Check for open issues or open a fresh issue to start a discussion around a feature idea or a bug
2. Fork the repository on GitHub to start making your changes to the master branch (or branch off of it
3. Write a test which shows that the bug was fixed or that the feature works as expected
4. Send a pull request and bug the maintainer until it gets merged and published
