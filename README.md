
GoLang Lib to Check for Dups
============================

Probably should have just used find | xargs sha1 | sort | uniq -c but
I'm doing it in go.  This lib is the build up two maps, one for "keys
seen just once" and one for "keys seen 2 or more tims" so you can
print the keys out quickly without going thru all the uniques.  It is
useful when you have a ton of uniques and just a few dups.

*Is it good to use?*

I'm using it.  

*What is it?*

Call "GetDups" and get the map of dups.  

Call "Add" with a string and it will be added to the
maps eventually.  

*Who owns this code?*

Chris Lane

*Adivce for starting out*

If you integrate, please let me or them know of your experience and
any suggestions for improvement.

The current API can best be seen in the _test files probably.  

*Requirements*

None.    
