# Yew Language Guide (Temp.)
## 1. Program Structure
Programs are broken into two main structures: (usually one) package and one or more modules.
### 1.1 Package
A package is a collection of related modules. When a package is imported into a module, the importer gets access to all things exported by the author of the imported package. 
#### 1.1.1 Importing Packages 
There are two parts required to access a package's interface. The general form of an import is
```
[<id>|_ =] import <package_name>
```
The first part is the left-hand-side. When an id is provided as a target, that will be used for accessing the package's interface. When `_` is used, the package's interface is directly accessible. Omitting the left-hand-side is identical to using the imported package's name for `<id>`. 

The second part is the call to import. 

### 1.2 Module
## Keywords
## Operators
