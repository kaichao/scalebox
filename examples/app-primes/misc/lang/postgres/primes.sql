CREATE OR REPLACE FUNCTION get_primes(INT4, INT4)
RETURNS INT4 AS
-- Returns the prime numbers in the range of $1 to $2
-- Contibuted by Melvin Davidson
-- Computer & Communication Technologies, Inc.
-- mdavidson(at)cctus(dot)com
$BODY$
DECLARE 
  v_start ALIAS FOR $1;
  v_end ALIAS FOR $2;
  v_test INT4;
  v_divisor INT4;
  v_prime_list TEXT DEFAULT '';
  v_msg TEXT;
  v_num INT4 DEFAULT 0;
BEGIN
	v_test = v_start;
	WHILE (v_test <= v_end) LOOP
		v_divisor = 2;
		WHILE (v_divisor <= v_test) LOOP
			IF mod(v_test, v_divisor) = 0 AND v_divisor < v_test THEN
				EXIT;
			ELSE 
				IF mod(v_test, v_divisor) = 0 AND v_divisor = v_test THEN
                    IF v_prime_list > '' THEN
						v_prime_list = v_prime_list ||  ',';
					END IF;
                    v_prime_list = v_prime_list || v_test::text;
					v_num = v_num + 1;
				END IF;
			END IF;
			v_divisor = v_divisor +1;
		END LOOP;
		v_test = v_test + 1;
	END LOOP;
	-- RETURN v_prime_list;
	RETURN v_num;
END;
$BODY$
LANGUAGE 'plpgsql' VOLATILE;

GRANT EXECUTE ON FUNCTION get_primes(INT4, INT4) TO public;
