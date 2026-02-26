# 多语言质数数量计算镜像

## 镜像构建

- 构建全部镜像
```sh
make clean-all
```

- 按语言构建镜像
```sh

cd app-primes/misc/lang
TAG=python make build
TAG=bash make build
TAG=php make build
TAG=nodejs make build
TAG=r make build
TAG=octave make build
TAG=julia make build
TAG=scala make build

TAG=java make build
TAG=rust make build
TAG=golang make build
TAG=fortran make build
TAG=c make build

TAG=postgres make build

```

## 镜像清除

```sh
make clean-all
```