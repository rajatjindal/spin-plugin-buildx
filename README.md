## spin buildx

A spin plugin to help with toolchains support for building spin apps. It is powered by containers and orchestrated by [dagger](https://dagger.io.

### Install:

```
spin plugin install -u https://github.com/rajatjindal/spin-plugin-buildx/releases/download/canary/buildx.json
```

### Build and export
Create a file `.toolchains` in your Spin App directory

e.g.

```
spin=3.0.0
golang=1.23.2
tinygo=0.34.0
```

and then run `spin buildx`

