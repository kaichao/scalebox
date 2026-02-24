function num_primes = get_num_primes(start, length)
    num_primes = 0;
    for i= 0:length-1
        if isprime(start+i)
            num_primes=num_primes+1;
        end;
    end;
endfunction;
