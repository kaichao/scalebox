FROM openjdk:11.0.16-slim

COPY primes Primes.java /src/

WORKDIR /src

RUN javac Primes.java
RUN jar cfe primes.jar Primes *.class
