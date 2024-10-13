# Examples

Here you can find example implementations of the interfaces in API

### Translator: translator.go
The notable thing here is an idea for how to create a bound on the translation. While not creative or particularly interesting, it should provide a starting point. 

One thing worth noting is certain obvious "translators" cannot be bounded. For example ... 
  - A general turing machine, `&Tm` can expose the Translator interface
  - Each call to `(&Tm).Translate()` modifies the internal state of `Tm` by doing one instruction on `Tm`'s "tape", returning `Tm`
  - Assume `(&Tm).Done() == true` if and only if `Tm` has halted  
  - Because of the Halting Problem, it's not possible to say whether `&Tm` will halt on arbitrary input
  - Consequently, one can't *always* provide a number of steps before `(&Tm).Done() == true`

Nonetheless, it is not required, and not expected, that `MaxTranslations` returns the same integer each time it is called. Merely that on a given state configuration of a translator it always returns the same integer. This is useful because as one gains more information, say one finds that `Tm` has a state configuration that we know an upper bound of, then one can reduce the number of unbounded cases.