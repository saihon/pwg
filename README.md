## pwg

<br/>

Password and username generator  
Download: [releases page](https://github.com/saihon/pwg/releases)

<br/>
<br/>

## example

<br/>
<br/>

#### password

<br/>

* Passowrd length
```
    $ pwg 10
    or
    $ pwg -d 10
```

* Number. same as no specification (0123456789)
```
    $ pwg -n
```

* Lowercase letters (abcdefghijklmnopqrstuvwxyz)
```
    $ pwg -l
```

* Uppercase letters (ABCDEFGHIJKLMNOPQRSTUVWXYZ)
```
    $ pwg -L
```

* Symbols (!"#$%&'()-=^~\|@`\[{;+:*]},<.>/?_)
```
    $ pwg -s
```

* Specify all
```
    $ pwg -nslL
    or
    $ pwg -a
    $ pwg --all
```

* Any string you want to use 
```
    $ pwg -c 'abc$%&@'
```

* Generate one or more
```
    $ pwg -g 10
```

* Evenly
```
    $ pwg -d 40 --all --evenly | ruby per.rb

    [input  ]: `5]WPTrj&%2'yd8b7HM`ji^0B3Afq9b*5.ED>14V
    [length ]: [  40]
    [lower  ]: [  10,  25%]
    [upper  ]: [  10,  25%]
    [number ]: [  10,  25%]
    [others ]: [  10,  25%]

```

<br/>
<br/>

#### username

<br/>

* Help
```
    $ pwg username --help
```

* Username length
```
    $ pwg username -d 6
    or
    $ pwg username 6
```

* Random length 5 ~ 9
```
    $ pwg username -r
```

* Generate one or more
```
    $ pwg username -g 10
```

* Capitalize initial letter
```
    $ pwg username -c
```

* Initial letter only 'a'
```
    $ pwg username -r -g 100 | grep ^a
```
