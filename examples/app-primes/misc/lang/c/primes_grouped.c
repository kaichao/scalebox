#include "util.h"
#include <omp.h>

int main(int argc, char *argv[]) {
    n_start = atol(argv[1]);
    length = atoi(argv[2]);
    group_size = atoi(argv[3]);

    init_timestamp();

    num_primes = 0;
#pragma omp parallel for schedule(dynamic) reduction(+ : num_primes)
    for(int n=0; n<length; n+=group_size){
        num_primes +=  get_num_primes(n_start+n,group_size);
    }

    printf("openmp grouped version\n");
    print_result();

    return 0;
}
