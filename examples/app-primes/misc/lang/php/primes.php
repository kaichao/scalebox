#!/usr/bin/php 

<?php
function isPrime($num) {
    if($num < 2) return false;
    for ($i = 2; $i < $num; $i++) {
        if ($num % $i == 0) {
          return false;
        }
    }
    return true; 
}

function getNumPrimes($start,$length){
  $ret = 0;
  for ($k = 0; $k < $length; $k++) {
    $n = $start + $k;
    if (isPrime($n)) {
      $ret++;
    }
  }
  return $ret; 
}

$start = $argv[1];
$length = $argv[2];

print(getNumPrimes($start,$length));
?>
