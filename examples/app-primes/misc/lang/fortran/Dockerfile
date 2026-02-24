FROM gcc:15-trixie

# FROM gcc:11-bullseye

COPY . /src/

WORKDIR /src/

RUN gfortran -static -lgfortran -O3 -static-libgfortran -static-libgcc -o primes primes.f08
