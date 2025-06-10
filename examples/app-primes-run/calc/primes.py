#!/usr/bin/python

import sys

def isPrime(n):
    if n<2:
        return False
    for i in range(2,n-1):
        if n%i==0:
            return False
    return True

def getNumPrimes(start, length):
	ret = 0
	for k in range(0, length):
		if isPrime(start + k):
			ret = ret + 1
	return ret

start=int(sys.argv[1])
length=int(sys.argv[2])

print(getNumPrimes(start,length))
