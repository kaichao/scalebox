#!/usr/bin/env julia

# Naive prime number check function
# Strictly follows pseudocode: iterate from 2 to n-1, test divisibility
function is_prime_naive(n)
    # Numbers less than 2 are not prime
    if n < 2
        return false
    end
    
    # Iterate from 2 to n-1
    for i in 2:(n-1)
        # If n is divisible by i, it is not prime
        if n % i == 0
            return false
        end
    end
    
    # If no divisor found, n is prime
    return true
end

function main()
    if length(ARGS) < 2
        println("Usage: julia primes.jl <start> <length>")
        println("Example: julia primes.jl 1 100")
        return
    end
    
    # Parse command-line arguments
    start_num = parse(Int, ARGS[1])
    len_num = parse(Int, ARGS[2])
    
    # Check validity of length
    if len_num <= 0
        println("Error: Length must be positive.")
        return
    end
    
    # Compute the end of the interval
    end_num = start_num + len_num - 1
    
    # Count prime numbers in the interval
    prime_count = 0
    for i in start_num:end_num
        if is_prime_naive(i)
            prime_count += 1
        end
    end
    
    # Output only the prime count (as required)
    println(prime_count)
end

# Program entry point
if abspath(PROGRAM_FILE) == @__FILE__
    main()
end