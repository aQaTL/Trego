@echo off
go run trego.go 2> errWindows
if not %ERRORLEVEL% == 0 type errWindows
