#include <stdio.h>
#include <math.h>
#include <string.h>
#include <malloc.h>
#include <stdlib.h>

int main(int argc, char* argv[])
{
    int num_procs = 3;
    int num_groups = 11;
    int n = (int)ceil((double)num_groups/num_procs);
    printf("%d\n",n);

    int* p_start = (int*) malloc(sizeof(int) * n * num_procs);
    memset((void*)p_start, 0,sizeof(int) * n * num_procs);

    int offset=0,direction = 1;
    for(int k=1;k<=num_groups;k++){
        p_start[offset]=k;
        offset += direction * n;
        if(offset<0 || offset >= n * num_procs){
            // next column in the same row
            offset = offset - direction * n + 1;
            // reverse diretion
            direction = -1 * direction;
        }
        printf("offset=%d, direction=%d\n\n",offset,direction);
    }

    printf("n=%d, num_procs=%d\n\n",n,num_procs);

    for(int i=0;i<num_procs;i++){
        for(int j=0;j<n;j++){
            printf("%d\t",p_start[i*n+j]);
        }
        printf("\n");
    }
    printf("\n");

    for(int k=0;k<n * num_procs;k++){
        printf("%d\t",p_start[k]);
    }
    return 0;
}