FROM gcc:15-trixie

COPY . /src/

WORKDIR /src/

RUN gfortran -static -lgfortran -O3 -static-libgfortran -static-libgcc -o primes primes.f08
