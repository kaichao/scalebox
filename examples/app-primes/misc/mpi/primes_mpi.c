// Author: Wes Kendall
// Copyright 2012 www.mpitutorial.com
// This code is provided freely with the tutorials on mpitutorial.com. Feel
// free to modify it for your own use. Any distribution of the code must
// either provide a link to www.mpitutorial.com or keep this header intact.
//
// Computes the number of prime numbers in parallel using MPI_Scatter and MPI_Gather
//
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <time.h>
#include <mpi.h>
#include <assert.h>

#include "util.h"

// Creates an array of task id list.
int *create_taskid_lists(int num_procs,int max) {
    int num_groups = max/group_size;
    int num_per_proc = (int)ceil((double)num_groups/num_procs);

    int* taskid_lists = (int*) malloc(sizeof(int) * num_per_proc * num_procs);
    memset((void*)taskid_lists, 0,sizeof(int) * num_per_proc * num_procs);

    int offset=0,direction = 1;
    for(int k=1;k<=num_groups;k++){
        taskid_lists[offset]=max - k * group_size + 1;
        offset += direction * num_per_proc;
        if(offset<0 || offset >= num_per_proc * num_procs){
            // next column in the same row
            offset = offset - direction * num_per_proc + 1;
            // reverse diretion
            direction = -1 * direction;
        }
    }

    return taskid_lists;
}

int compute_num_of_primes(int* sub_taskid_list, int num_per_proc) {
    int sum_primes = 0;
    for(int i=0;i<num_per_proc;i++){
        n_start = sub_taskid_list[i];
        if(n_start > 0) {
        // do prime calculation
            int num_primes = 0;
#pragma omp parallel for schedule(dynamic) reduction(+ : num_primes)
            for(int n=0; n<group_size; n++){
                num_primes +=  is_prime(n_start+n);
            }
            sum_primes += num_primes;
        }
    }
    return sum_primes;
}

int main(int argc, char** argv) {
    if (argc != 2) {
        fprintf(stderr, "Usage: prime max_num\n");
        exit(1);
    }

    int num_max = atoi(argv[1]);
    group_size = 10000;
    n_start=1;
    length=num_max;

    MPI_Init(&argc, &argv);
    int world_rank, world_size;
    MPI_Comm_rank(MPI_COMM_WORLD, &world_rank);
    MPI_Comm_size(MPI_COMM_WORLD, &world_size);

  // Create taskid list on the root process.
    int *taskid_list = NULL;
    if (world_rank == 0) {
        init_timestamp();
        taskid_list = create_taskid_lists(world_size, num_max);
    }

  // For each process, create a buffer that will hold a subset of the entire array
    int num_per_proc = (int)ceil((double)num_max/group_size/world_size);
    int* sub_taskid_list = (int *)malloc(sizeof(int) * num_per_proc);

  // Scatter the task list from the root process to all processes in the MPI world
    MPI_Scatter(taskid_list, num_per_proc, MPI_INT, sub_taskid_list, 
            num_per_proc, MPI_INT, 0, MPI_COMM_WORLD);

  // Compute the number of your subset
    int sub_result = compute_num_of_primes(sub_taskid_list,num_per_proc);

#ifdef DEBUG
    printf("id=%d, num_primes=%d\n",world_rank,sub_result);
#endif

  // Gather all partial results down to the root process
    int *sub_results = NULL;
    if (world_rank == 0) {
        sub_results = (int *)malloc(sizeof(int) * world_size);
        assert(sub_results != NULL);
    }
    MPI_Gather(&sub_result, 1, MPI_INT, sub_results, 1, MPI_INT, 0, MPI_COMM_WORLD);

  // compute the final result
    if (world_rank == 0) {
        num_primes = 0;
        for(int i=0;i<world_size;i++){
            num_primes += sub_results[i];
        }
        n_start=1;
        print_result();
    }

  // Clean up
    if (world_rank == 0) {
        free(taskid_list);
        free(sub_results);
    }
    free(sub_taskid_list);

    MPI_Barrier(MPI_COMM_WORLD);
    MPI_Finalize();
}
