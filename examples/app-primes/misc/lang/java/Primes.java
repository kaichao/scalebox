import java.time.Duration;
import java.time.LocalDateTime;
import java.util.stream.Stream;

public class Primes {
    public static void main(String[] args){
        if(args.length != 2){
            System.out.println("Usage : num_primes ${start_number} ${range_size}");
            System.exit(1);
        }
        int startNumber=0,rangeSize=0;
        try {
            startNumber = Integer.parseInt(args[0]);
            rangeSize = Integer.parseInt(args[1]);
        } catch (Exception ex){
            ex.printStackTrace();
            System.out.println("wrong arguments!");
            System.exit(2);
        }

        int numParallel;
        try{
            numParallel = Integer.parseInt(System.getenv("NUM_PARALLEL"));
        } catch (Exception ex){
            numParallel = Runtime.getRuntime().availableProcessors();
        }
        System.setProperty("java.util.concurrent.ForkJoinPool.common.parallelism",
                String.valueOf(numParallel));

        LocalDateTime start = LocalDateTime.now();
        long numPrimes = Stream.iterate(startNumber, n -> n + 1)
                .limit(rangeSize)
                // disable parallel 
                // .parallel()
                .filter(n->isPrime(n))
                .count();                
        LocalDateTime end = LocalDateTime.now();
        Duration duration = Duration.between(start,end);
        // System.out.format("%d ms\n", duration.toMillis());

        System.out.println(numPrimes);
    }
    private static boolean isPrime(int n){
        if(n < 2) return false;
        for(int i=2;i<n;i++)
            if(n % i == 0) return false;
        return true;
    }
}
