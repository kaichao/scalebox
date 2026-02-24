IsPrime <- function(n) {
    if (n == 1) {
        return(FALSE)
    }
    if (n == 2) {
        return(TRUE)
    }
    for(i in 2:(n-1)){
        if(n %% i == 0){
            return(FALSE)
        }
    }
    return(TRUE)
}

GetNumPrimes <- function(n0,len){
    ret <- 0
    for(i in n0:(n0+len-1)){
        if(IsPrime(i)){
            ret <- ret+1
        }
    }
    return(ret)
}

args <- commandArgs(trailingOnly = TRUE)

start <- as.integer(args[1])
length <- as.integer(args[2])

print (start)
print (length)
numPrimes <- GetNumPrimes(start,length)

print(numPrimes)
