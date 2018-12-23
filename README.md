# Random Hashing Tool

The tool generates hashes for input number and list of random numbers based on input number

[![Go Report Card](https://goreportcard.com/badge/github.com/ilyakaznacheev/rand-hashing)](https://goreportcard.com/report/github.com/Masterminds/glide) [![Build Status](https://travis-ci.org/ilyakaznacheev/rand-hashing.svg?branch=master)](https://travis-ci.org/ilyakaznacheev/rand-hashing) [![Coverage Status](https://coveralls.io/repos/github/ilyakaznacheev/rand-hashing/badge.svg?branch=master)](https://coveralls.io/github/ilyakaznacheev/rand-hashing?branch=master)

## Installation
```bash
go get -v github.com/ilyakaznacheev/rand-hashing/...
```

## Usage

Run client, that will listen and print all generated hashes

```bash
go run cmd/client/client.go
```

After that, you can run any number of hash generators

```bash
./gethashes number iterations
```

- number - base number for hashing. Should be 6 or more digits length
- iterations - a number of hashes, generated for random numbers, based on base number

Note that file `gethashes` has to have execution permissions.

Generated keys are also stored into Redis list with key "randhash".

## Configuration

Edit config file `configs/config.yml` to setup Redis connection settings
