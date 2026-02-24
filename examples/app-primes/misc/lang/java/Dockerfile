FROM openjdk:11-bullseye

COPY primes Primes.java /src/

WORKDIR /src

RUN javac Primes.java
RUN jar cfe primes.jar Primes *.class
