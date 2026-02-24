#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>

#define min(m,n) ((m)<(n)) ? (m) : (n)

// The starting value of the interval to be calculated
long n_start;
// The length of the interval to be calculated
long length;
// Calculated unit length by group
int group_size;
// Result values
int num_primes;

// single value calculation
int is_prime(int n){
    if(n<2) return 0;
    for(int i=2;i<n;i++)
        if(n % i == 0) return 0;
    return 1;
}

// Calculate by group
int get_num_primes(int start, int len){
    int ret = 0;
    for(int k=0;k<len;k++){
        ret += is_prime(start+k);
    }
    return ret;
}

static struct timeval t_start, t_end;
void init_timestamp() {
    gettimeofday(&t_start, NULL); 
}

void print_timestamp() {
    gettimeofday(&t_end, NULL); 
    long interval = (t_end.tv_sec-t_start.tv_sec)*1000
                  + (t_end.tv_usec-t_start.tv_usec)/1000;
    printf("%ld miliseconds\n", interval);
    t_start = t_end;
}

void print_result(){
    printf("[%ld,%ld] prime count:%d,",n_start, n_start+length-1, num_primes);
    print_timestamp();
}

static char buf[20];
void print_int(int n){
    int i=0;
    while(n){
        buf[i++]=n%10;
        n/=10;
    }
    while(i--){
        putchar('0'+buf[i]);
    }
    putchar('\n');
}
