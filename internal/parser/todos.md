- [ ] Types cannot end with implicit types 
  - this could be done during syntax analysis, name analysis, or type analysis
    - probably easiest to do this during type analysis--name analysis could work too, but makes less sense

- [ ] This is important: aliases can ***ONLY*** alias *type* constructors, make sure this gets enforced during name analysis

- [x] add `auto` modifier to type sigs