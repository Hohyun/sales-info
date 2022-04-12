;; Vectis Sales Data Download Script
;
;  : Written by hohynkim@jinair.com
;  : Last edited date: 2020-09-15
;
;  : Arguments -- id, password, fromDate, toDate
;
;  : How to run
;  :: AutoIt3.exe sales_download.au3 vectis_id vectis_pswd 2020-09-01 2020-09-07


#include <Date.au3>
#include <MsgBoxConstants.au3>

Local $paramCnt = $CmdLine[0]
Local $id         = $CmdLine[1]
Local $pswd     = $CmdLine[2]
Local $fromDate, $toDate

if  $paramCnt == 4 Then
	Local $fromDate   = $CmdLine[3]
    Local $toDate     = $CmdLine[4]
ElseIf $paramCnt == 2 Then
	$fromDate = DefaultFromDate()
	$toDate = _DateTimeFormat(_DateAdd('d', -1, _NowCalcDate()), 2)
Else
	Exit
EndIf

; MsgBox($MB_SYSTEMMODAL, "", "ID: " & $id & " Password: " & $pswd & " From: " & $fromDate & ", To: " & $toDate)

VectisLogin($id, $pswd)
DownloadReport($fromDate, $toDate)
WinClose("Vectis")
Send("{Enter}")

; copy file
FileCopy("C:\VectisClient\VectisTemp\VectisReport.csv", "D:\Projects\sales-info\data\VectisReport_sales.csv", $FC_OVERWRITE)
Exit


Func DefaultFromDate()
	Local $fromDate
	Local $iWeekday = _DateToDayOfWeek(@YEAR, @MON, @MDAY)
	if $iWeekday == 2 Then  ; Monday
		$fromDate = _DateTimeFormat(_DateAdd('d', -3, _NowCalcDate()), 2)
	Else
		$fromDate = _DateTimeFormat(_DateAdd('d', -1, _NowCalcDate()), 2)
	EndIf
	Return $fromDate
EndFunc

Func VectisLogin($id, $pswd)
   ;~ ; Run Vectis
   Run("C:\VectisClient\bin\jade.exe appServer=10.23.34.4 appServerPort=6021 app=Vectis schema=AppSchema")

   Local $hWnd = WinWaitActive("Welcome to Vectis - Logon")
   ControlSetText($hWnd, "", "Jade:Edit1", $id)
   ControlSetText($hWnd, "", "Jade:Edit2", $pswd)
   ControlClick($hWnd, "", "Jade:JadeMask1")

   Local $hWnd = WinWaitActive("Vectis")
EndFunc

Func DownloadReport($fromDate, $toDate)
   Local $hWnd = WinActivate("Vectis")
   Sleep(1000)
   Send("!R{RIGHT}{ENTER}")    ; Report - Passenger Revenue - Sales

   Local $hWnd = WinWaitActive("Vectis - [Passenger Sales Reports]")
   ControlClick($hWnd, "", "Jade:ListBox1", "left", 2, 134, 118)

   Local $hWnd = WinWaitActive("Vectis - [Sale FOP Manager]")
   ; Filter On: Settlement Date
   ControlClick($hWnd, "", "Jade:Edit1")
   ControlSetText($hWnd, "", "Jade:Edit1", $fromDate)  ; Date From:
   Send("{TAB}")
   ControlClick($hWnd, "", "Jade:Edit2")
   ControlSetText($hWnd, "", "Jade:Edit2", $toDate)  ; Date To:
   Send("{TAB}")
   ControlClick($hWnd, "", "Jade:Button23")

   ; Sometimes following confirm window appears
   WinWait("Report Period", "", 5)
   ; select "Yes" - Date is correct?
   If WinExists("Report Period") Then
      ControlClick("Report Period", "", "Button1")
   EndIf

   Local $hWnd = WinWaitActive("Printer Options")
   ControlClick($hWnd, "", "Button10")  ; Export
   Sleep(10000)
   ControlClick($hWnd, "", "ComboBox1")      ; CSV
   Sleep(3000)
   Send("{UP}{TAB}")
   ControlClick($hWnd, "", "Jade:Button3")        ; OK --> Run

   Local $hWnd = WinWaitActive("File/s Created")
   ControlClick($hWnd, "", "Button1")
   WinActivate("Vectis")
EndFunc