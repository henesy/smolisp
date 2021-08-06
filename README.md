# Smolisp

A toy lisp-like language implementation.

## Build

	go build

## Usage

EOF exits the repl.

## Examples

	$ ./smolisp
	» (+ 2.14 3.2)
	5.34
	» (- 5.34 π)
	2.1984073464102067
	» (/ 4 4)
	1
	» (* 2 billion)
	2000000000
	»

## Specification

The syntax is as per S-expression where the form is:

	(procedureName argA argB ...)

The implemented procedures:

	+
	-
	*
	/

The implemented types:

	Integral
	Floating

The implemented keywords:

	π
	billion
