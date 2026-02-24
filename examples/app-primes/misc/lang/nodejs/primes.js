function isPrime(num) {
    if(num < 2) return false;
    for (let i = 2; i < num; i++) {
        if (num % i == 0) {
          return false;
        }
    }
    return true; 
}

function getNumPrimes(start,length){
  var ret = 0;
  for (var k = 0; k < length; k++) {
    let n = parseInt(start) + parseInt(k);
    if (isPrime(n)) {
      ret++;
    }
  }
  return ret; 
}

const args = process.argv.slice(2)
console.log(getNumPrimes(args[0],args[1]));
