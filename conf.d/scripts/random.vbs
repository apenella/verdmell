Dim range, exitvalue
Dim output

if WScript.Arguments.Count <> 2 then
    WScript.Echo "Parameters incorrect"
    WScript.Quit(4)
end if

range = Wscript.Arguments(0)
output = Wscript.Arguments(1)

Randomize
exitvalue = Int((range + 1) * Rnd )

Wscript.Echo output
Wscript.Quit(exitvalue)