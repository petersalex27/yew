- [ ] Types cannot end with the type of an implicit parameter--unclear what this would even mean given that the output of a function is not a parameter ...
  - this could be done during syntax analysis, name analysis, or type analysis
    - probably easiest to do this during type analysis--name analysis could work too, but makes less sense

- [ ] This is important: aliases can ***ONLY*** alias *type* constructors, make sure this gets enforced during name analysis

- [x] add `auto` modifier to type sigs