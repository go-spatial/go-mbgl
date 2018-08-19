# mbgl

This directory contains direct go bindings against the Mapbox GL Native C++ API.

## Design overview

* follow the C++ api as closely as possible, not idiomatic Go
    * keeping in mind Go can only interface with C defintions
* manual memory management through `(T).Destruct()` methods
* C++ types are aliased in C as empty structs
    * they are named `MbglTypeName`
* [WIP] C++ classes are aliased in Go as the C type
    * ie. `type TypeName C.TypeName`
* C++ classes will always be passed around as pointers to C aliases and `reinterpret_cast<>`ed 
* [WIP] C++ structs are represented in Go as normal structs with Go memebers
    * ie. `type TypeName struct{ field1 *TypeNameX, field2 *TypeNameY }`
    * they must be accesse through a `(TypeName).cPtr() *C.MbglTypeName` method which updates the values of a C++ instance of the struct with the current values of the Go struct
        * this keeps us from having to write a getter and setter for every struct
* Memory management is the responsibility of the programmer, through use of `(TypeName).Destruct()` methods which will call `delete` on the object

## Using the bindings

* in the higher level package, `runtime.SetFinalizer` is used for memory management, note that this cannot be done for classes (which are aliased as empty structs) as: 

> It is not guaranteed that a finalizer will run if the size of *obj is zero bytes.
> https://golang.org/pkg/runtime/#SetFinalizer

* after calling `(T).Destruct()` make sure to not reuse the variable, this will result in strange errors; another `(T).Destruct()` call might yeild something like "freed block must be allocated with malloc".