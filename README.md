# bruteforce-referendum-1O

This repository contains some scripts that can brute-force the voting locations of the illegal referendum of
independence of Catalonia of the 1st of October, 2017.

The aim of this project is just to highlight the very low level of security of the referendum's main web page as well
as to highlight how easily personal information can be leaked from there. Note that the brute-forcing is done
locally.

The dumped database files are not included, but the javascript script automatically downloads it from one of the
referendum's mirrors. If it is down, please send me a message and I will consider sending you a dumped copy, but the
decryption is your totally up to you.

The golang script is configurable, for seeing the options, just use '-h'.

Note that this is only a sketch, more optimizations may be added in the future.