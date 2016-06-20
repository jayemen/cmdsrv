cmdsrv
======
Launches an HTTP server that outputs the result of running an executable.

Example usage:
`cmdsrv -cmd "echo" -args "test" -cache-time 10 -listen :7777`

This would start a server listening on port 7777, that outputs the result of 'echo test'. Results are cached for 10s.


 Notes
 -----
 This is literally the only Go program I have ever written. It is likely riddled with bugs and anti-patterns.