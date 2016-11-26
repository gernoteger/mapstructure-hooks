# About mapstructure-hooks

This is a small extension to Mitshell Hashimoto's excellent library [mapstructure](https://github.com/mitchellh/mapstructure).
It allows one to fill arrays of interfaces whose concrete type is determined by their content. A typical usecase is configuring 
log handlers from a yaml file and processing the imported map through mapstructure.
