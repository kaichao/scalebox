function ret = is_prime(n)
    if n < 2
        ret = false;
        return;
    end;
    for i = 2:n-1
        if ~mod(n,i)
            ret = false;
            return;
        end;
    end;
    ret=true;
endfunction;
