#!/usr/bin/env julia

# 朴素质数判断函数
# 严格按照伪代码: 从2到n-1，通过是否可以整除判断
function is_prime_naive(n)
    # 小于2的数不是质数
    if n < 2
        return false
    end
    
    # 从2遍历到n-1
    for i in 2:(n-1)
        # 如果n能被i整除，则不是质数
        if n % i == 0
            return false
        end
    end
    
    # 如果没有找到能整除n的数，则是质数
    return true
end

# 质数计数函数
# 严格按照伪代码: 从2到n，判断每个数是否为质数
function count_primes_naive(n)
    count = 0
    
    # 从2遍历到n
    for i in 2:n
        if is_prime_naive(i)
            count += 1
        end
    end
    
    return count
end

function main()
    if length(ARGS) < 1
        println("Usage: julia naive_prime.jl <positive_integer>")
        println("Example: julia naive_prime.jl 30")
        return
    end
    
    # 解析输入参数
    n = parse(Int, ARGS[1])
    
    # 检查输入有效性
    if n < 1
        println("Error: Please provide a positive integer (>= 1)")
        return
    end
    
    # 计算质数数量
    start_time = time()
    prime_count = count_primes_naive(n)
    end_time = time()
    
    # 输出结果
    println("Range: 1 to $n")
    println("Number of primes: $prime_count")
    println("Execution time: $(round(end_time - start_time, digits=4)) seconds")
    
    # 可选：输出找到的质数（对于较大的n，建议注释掉这部分）
    if n <= 100
        print("Primes found: ")
        primes = []
        for i in 2:n
            if is_prime_naive(i)
                push!(primes, i)
            end
        end
        println(primes)
    end
end

# 程序入口
if abspath(PROGRAM_FILE) == @__FILE__
    main()
end