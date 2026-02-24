#include "util.h"
#include <math.h>
#include <stdlib.h>

// primes in [1..65535]
int is_short_prime(unsigned n){
    int i;
    int k=(int)sqrt((double)n);
    for(i=2;i<=k;i++)
        if(n % i == 0) break;
    return i>k;
}

#define NUM 6542
static unsigned short_primes[NUM];

int is_int_prime(unsigned n){
    int i;
    for(i=0;i<NUM;i++)
        if(n % short_primes[i] == 0) break;
    return i==NUM;
}

int cmp (const void *a , const void *b){
    return *(int*)a - *(int*)b; 
}
int main(int argc, char *argv[]) {
    n_start=1;
    // length=1L<<32;
    length=atol(argv[1]);

    init_timestamp();

    FILE* fp = fopen("/tmp/primes.dat", "wb");
    int k=0;
    for (unsigned n=2;n<0xffff;n++){
        if(is_short_prime(n)){
            short_primes[k++]=n;
        }
    }
    fwrite(short_primes, sizeof(int), NUM, fp);
    num_primes += k;

    unsigned buff[10000];
    for (long n=1<<16;n<n_start+length-1;n+=(1<<16)){
        k = 0;
        unsigned upper = min(n+(1<<16), n_start+length-1);
// #pragma omp parallel for shared(buff)
#pragma omp parallel for
        for(unsigned i=n;i<upper;i++){
            if(is_int_prime(i)){
#pragma omp critical
                buff[k++]=i;
            }
        }
        num_primes += k;
        qsort(buff,k,sizeof(int),cmp);
        fwrite(&n, sizeof(int), k, fp);
        printf("n=%ld,\tupper=%ld,\t%ld\n",n,upper,n>>16);
    }
    fclose(fp);
    print_result();
}
