module main

go 1.25

require (
	github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef
)


replace (
	github.com/coreycole/datastarui => ../../datastarui
	github.com/coreycole/datastarui/serviceworker => ../../serviceworker
)