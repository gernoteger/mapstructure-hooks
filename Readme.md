[![GoDoc](https://godoc.org/github.com/gernoteger/mapstructure-hooks?status.svg)](https://godoc.org/github.com/gernoteger/mapstructure-hooks)

# About

This is a small extension to Mitshell Hashimoto's excellent library [mapstructure](https://github.com/mitchellh/mapstructure).
It allows one to fill arrays of interfaces whose concrete type is determined by their content. A typical usecase is configuring 
log handlers from a yaml file and processing the imported map through mapstructure.

# Usage

For a detailled description look at the [examples](https://godoc.org/github.com/gernoteger/mapstructure-hooks#example-Decode) 
in the [godocs](https://godoc.org/github.com/gernoteger/mapstructure-hooks).