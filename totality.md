# Totality (from Idris Docs)
For a function `func` to be considered (conservatively) total, it must ...
- cover all possible inputs
- be *well-founded*: by the time a sequence of (possibly mutually) recursive calls reaches `func` again, it must be possible to show that one of its arguments has decreased
- not use any data types which are not *strictly positive*
- not call any non-total functions

***_NOTE_***: this is not perfect. There are lots of total programs that aren't counted as total and there might be things it says are total that are not!