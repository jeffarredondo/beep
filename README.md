beep
====

A simple Go program that is useful for alerting the end of a long running command execution.
It can also play music sheet from stdin. To play a demo music, run: ```$ beep -p | beep -m```

Listen demo: [demo-mozart-k33b.mp3](https://github.com/dbatbold/beep/blob/master/demo-mozart-k33b.mp3))
```
Compiling:
 $ apt-get install golang libasound2-dev  # for Debian and Ubuntu
 $ git clone http://github.com/dbatbold/beep
 $ cd beep
 $ go build
 $ strip beep  # optional
 $ cp beep /usr/bin/  # as root

 Windows port is not implemented yet.

Usage: beep [options]
  -c=1: count
  -d="default": audio device, example: "hw:0,0"
  -f=0.07459: frequency
  -h: help
  -l: beep per line from stdin
  -m: play music notes from stdin (see beep notation)
  -p: print a demo music by Mozart
  -t=1: time duration (1-100)
  -v=100: volume (1-100)
  -b: send bell to PC speaker
  -q: quiet stdout while playing music
  -n: print notes while playing music
  -o=file: output music wave to file (WAV format)

Beep notation:
  | | | | | | | | | | | | | | | | | | | | | | 
  |2|3| |5|6|7| |9|0| |=|a|s| |f|g| |j|k|l| |
 | | | | | | | | | | | | | | | | | | | | | | 
 |q|w|e|r|t|y|u|i|o|p|[|]|z|x|c|v|b|n|m|,|.|

 q - middle C (261.6 hertz)

 Left and right hand keys are same. Uppercase 
 letters are control keys. Lowercase letters
 are music notes. Space bar is current duration
 rest. Spaces after first space are ignored.

 Control keys:

 Rest:
 RW     - whole rest
 RH     - half rest
 RQ     - quarter rest
 RE     - eighth rest
 RS     - sixteenth rest
 RT     - thirty-second rest

 Space  - eighth rest, depends on current duration

 Durations:
 DW     - whole note
 DH     - half note
 DQ     - quarter note
 DE     - eighth note
 DS     - sixteenth note
 DT     - thirty-second note

 Octave:
 HL     - switch to left hand keys
 HR     - switch to right hand keys

 Clef:
 CB     - G and F clef partition (Base)

 Measures:
 |      - bar (ignored)

Demo Music: Mozart K33b:
# Mozart K33b
HRDEc c DSc s z s |DEc DQz DE[ CB
HLDE[ n   z   ,      c HRq HLz ,

HRDEc c DSc s z s |DEc DQz DE[ CB
HLDE[ n   z   ,      c HRq HLz ,

HRDEv v DSv c s c |DEv s ] v CB
HLDEs l   z   ,      ] m p b

HRDEc c DSc s z s |DEc z [ c CB
HLDEz ,   ]   m      [ n o v 

HRDEs s DSs z ] z |DEs ] p s CB
HLDE] m   [   n      p b i c 

HRDSs z ] [ z ] [ p |DE[ DSi y DQr CB
HLDEn   z   s   c      n   c     [

Usage Examples:

 $ cp -vr directory target; beep
 $ ffmpeg -i video.mp4 -vn -acodec libmp3lame sound.mp3; beep -t 3
 
 # alarm for download completion
 $ curl -O http://host.com/bigfile.tgz; beep -c 4 -f 0.012
 
 # beep for every text file found under home
 $ find ~ -name '*.txt' | beep -l
 
 # set an alarm for 1 hour from now
 $ sh -c 'sleep 3600; beep -t 3 -c 6' &
 
 # play all music notes
 # echo "q2w3er5t6y7ui9o0p[=]azsxcfvgbnjmk,l." | beep -m
 
 # play Mozart K33b
 $ beep -p | beep -m
```
