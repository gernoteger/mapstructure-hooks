[![GoDoc](https://godoc.org/github.com/gernoteger/mapstructure-hooks?status.svg)](https://godoc.org/github.com/gernoteger/mapstructure-hooks)
[![Go Report Card](https://goreportcard.com/badge/gernoteger/mapstructure-hooks)](https://goreportcard.com/report/gernoteger/mapstructure-hooks)

# About

This is a small extension to Mitshell Hashimoto's excellent library [mapstructure](https://github.com/mitchellh/mapstructure).
It allows one to fill arrays of interfaces whose concrete type is determined by their content. A typical usecase is configuring 
log handlers from a yaml file and processing the imported map through mapstructure.

The design was heavily inspired by [logrus_mate](https://github.com/gogap/logrus_mate)

# Usage

For a detailled description look at the [examples](https://godoc.org/github.com/gernoteger/mapstructure-hooks#example-Decode) 
in the [godocs](https://godoc.org/github.com/gernoteger/mapstructure-hooks).