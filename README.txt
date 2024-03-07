preloader
---------
This project serves as a scaffold to my DLL hijack-based preloaders for various patches.

Usage:
 - Build the Go project in _sourcegenerator
 - Run the source generator to create the required files
 - Place all of them in preloader/Proxy/gen
 - Load the project up in Visual Studio
 - Extend Core::Init