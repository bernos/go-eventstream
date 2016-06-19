# go-eventstream

A package for creating composable concurrent event streams. 

The go programming language features primitives such as channels and go routines for writing concurrent programs. While these primitives can be used to implement powerful concurrency patterns, doing so over and over again can be tedious and error-prone.

The go-eventstream package provides a small set of simple, composable abstractions over these base primitives. Inspired by concepts from functional reactive programming, these abstractions can be used to compose powerful, concurrent programs that are easy to reason about, and easy to test.

## Event

