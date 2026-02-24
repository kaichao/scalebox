#include "util.h"
#include <math.h>
#include <stdlib.h>

#define NUM 10000
static unsigned short_primes[NUM];

int binary_search(int key){
    int low = 0;
    int high = 6541;
    int mid;
    while(low<= high){
        mid = (low + high)/2;
        int midVal = short_primes[mid];
        if(midVal<key)
            low = mid + 1;
        else if(midVal>key)
            high = mid - 1;
        else
            return mid;
    }
    return mid;
}
// primes in [1..65535]
int is_short_prime(unsigned n){
    int i;
    int k=(int)sqrt((double)n);
    for(i=2;i<=k;i++)
        if(n % i == 0) break;
    return i>k;
}


int is_int_prime(unsigned n){
    int k=(int)sqrt((double)n);
    int upper = binary_search(k);

    int i;
    for(i=0;i<=upper;i++)
        if(n % short_primes[i] == 0) break;
    return i==(upper+1);
}

int main(int argc, char *argv[]) {
    int num = atol(argv[1]);
    int num_primes = 0;
    init_timestamp();

    int k=0;
    for (unsigned n=2;n<0xffff;n++){
        if(is_short_prime(n)){
            short_primes[k++]=n;
        }
    }
    num_primes = k;
    print_timestamp();

    if(num < 0xffff){
        num_primes = binary_search(num);
        printf("[%d..%d]: %d primes.", 1, num, num_primes);
        print_timestamp();
        return 0;
    }
    
// #pragma omp parallel for schedule(dynamic) reduction(+ : num_primes)
    for(int n=0x10001; n<num; n+=2){
        num_primes +=  is_int_prime(n);
    }
    printf("[%d..%d]: %d primes, ", 1, num, num_primes);
    print_timestamp();
    return 0;
}
