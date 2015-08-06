Dim InfoType, ServiceName
Dim NamedArgs
Set NamedArgs = WScript.Arguments.Named

Set wmi = GetObject("winmgmts://./root/cimv2")


For Each Arg In NamedArgs
	Select Case Arg
		Case "q", "query"
			InfoType = NamedArgs.Item(Arg)
		Case "s", "service"
			ServiceName = NamedArgs.Item(Arg)
	'end type select'
	End Select
Next

if InfoType = "" or ServiceName = "" Then
	WScript.Echo "Parameters incorrects- q:[status|startmode] query, s: service name"
	WScript.Quit(4)
End If

Select Case InfoType
	Case "status"
		status = wmi.Get("Win32_Service.Name='" & serviceName & "'").Status
		WScript.Echo "The service status is " & status

		if status = "OK" Then
			WScript.Quit(2)
		Else
			WScript.Quit(0)
		End if

	Case "startmode"
		startmode = wmi.Get("Win32_Service.Name='" & serviceName & "'").StartMode
		WScript.Echo "The service start mode is " & startmode

	Case Else
		WScript.Echo "Parameters incorrects- q:[status|startmode] query, s: service name"
		WScript.Quit(4)
'end InfoType'
End Select

