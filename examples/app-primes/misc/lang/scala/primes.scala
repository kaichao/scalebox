if (args.length < 2) {
  println("Usage: scala prime.scala <start> <length>")
  sys.exit(1)
}
val start = args(0).toInt
val length = args(1).toInt
val end = start + length - 1

def isPrime(n: Int): Boolean = {
  if (n < 2) false
  else (2 until n).forall(n % _ != 0)   // 朴素算法：试除从2到n-1
}

val primeCount = (start to end).filter(isPrime).length
println(primeCount)