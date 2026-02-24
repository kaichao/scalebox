PROGRAM PRIMES
    INTEGER::n,len,get_num_primes

    CHARACTER(len=32) :: arg
    CALL get_command_argument(1, arg)
    READ(arg,"(I11)") n
    CALL get_command_argument(2, arg)
    READ(arg,"(I11)") len

    WRITE(*,"(I0)") get_num_primes(n,len)

!    STOP 'OK'

END PROGRAM PRIMES

LOGICAL FUNCTION is_prime(n)
    IMPLICIT NONE
    INTEGER::n,i
    is_prime = n > 1
    do i = 2, n-1
        if(mod(n,i)==0)then
            is_prime = .false.
        end if
    end do
END FUNCTION is_prime

INTEGER FUNCTION get_num_primes(num,len)
    IMPLICIT NONE
    LOGICAL is_prime
    INTEGER::num,len,i
    get_num_primes = 0
    do i = num, num+len-1
        if (is_prime(i)) then
            get_num_primes = get_num_primes + 1
        end if
    end do
END FUNCTION get_num_primes
