#include "util.h"
#include <omp.h>

int main(int argc, char *argv[]) {
    n_start = atol(argv[1]);
    length = atoi(argv[2]);

#ifndef CONTAINER_APP
    init_timestamp();
#endif
    num_primes = 0;
#pragma omp parallel for schedule(dynamic) reduction(+ : num_primes)
    for(int i=0; i<length; i++){
        num_primes +=  is_prime(n_start+i);
    }

#ifdef CONTAINER_APP
    print_int(num_primes);
#else
    printf("openmp version\n");
    print_result();
#endif

    return 0;
}
