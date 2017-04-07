all:
	gcc -c unix_test.c
	gcc -c flex_test.c
	gcc -o l4_unix_test unix_test.o l4_test.c
	gcc -o l4_flex_test flex_test.o l4_test.c

clean:
	rm -rf unix_test.o flex_test.o l4_unix_test l4_flex_test
