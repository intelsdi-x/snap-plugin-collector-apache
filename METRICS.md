| Namespace                                  | Data Type | Description (optional)                                              | Safe               | Default       |
|--------------------------------------------|-----------|---------------------------------------------------------------------|--------------------|---------------|
| /intel/apache/BusyWorkers                  | float64   | Busy workers                                                        | :white_check_mark: | 0             |
| /intel/apache/BytesPerReq                  | float64   | Bytes transferred per request                                       | :x:                | 0             |
| /intel/apache/BytesPerSec                  | float64   | Bytes transferred per second                                        | :white_check_mark: | 0             |
| /intel/apache/CPULoad                      | float64   | CPU load                                                            | :white_check_mark: | 0             |
| /intel/apache/ConnsAsyncClosing            | float64   | Asynchronous closing connections                                    | :white_check_mark: | 0             |
| /intel/apache/ConnsAsyncKeepAlive          | float64   | Asynchronous keepalive connections                                  | :white_check_mark: | 0             |
| /intel/apache/ConnsAsyncWriting            | float64   | Asynchronous writing connections                                    | :white_check_mark: | 0             |
| /intel/apache/ConnsTotal                   | float64   | Total connections                                                   | :white_check_mark: | 0             |
| /intel/apache/IdleWorkers                  | float64   | Idle workers                                                        | :white_check_mark: | 0             |
| /intel/apache/ReqPerSec                    | float64   | Requests per second                                                 | :white_check_mark: | 0             |
| /intel/apache/Total_Accesses               | float64   | Total accesses                                                      | :white_check_mark: | 0             |
| /intel/apache/Total_kBytes                 | float64   | Total kBytes                                                        | :white_check_mark: | 0             |
| /intel/apache/Uptime                       | float64   | Server uptime                                                       | :white_check_mark: | 0             |
| /intel/apache/workers/Closing              | float64   | Closing connection                                                  | :white_check_mark: | 0             |
| /intel/apache/workers/DNSLookup            | float64   | DNS Lookup                                                          | :white_check_mark: | 0             |
| /intel/apache/workers/Finishing            | float64   | Gracefully finishing                                                | :white_check_mark: | 0             |
| /intel/apache/workers/Idle_Cleanup         | float64   | Idle cleanup of worker                                              | :white_check_mark: | 0             |
| /intel/apache/workers/Keepalive            | float64   | Keepalive (read)                                                    | :white_check_mark: | 0             |
| /intel/apache/workers/Logging              | float64   | Logging                                                             | :white_check_mark: | 0             |
| /intel/apache/workers/Open                 | float64   | Open slot with no current process                                   | :white_check_mark: | 0             |
| /intel/apache/workers/Reading              | float64   | Reading Request                                                     | :white_check_mark: | 0             |
| /intel/apache/workers/Sending              | float64   | Sending Reply                                                       | :white_check_mark: | 0             |
| /intel/apache/workers/Starting             | float64   | Starting up                                                         | :white_check_mark: | 0             |
| /intel/apache/workers/Waiting              | float64   | Waiting for Connection                                              | :white_check_mark: | 0             |
| /intel/apache/ServerVersion                | string    | Apache server version                                               | :x:                | `"Not Found"` |
| /intel/apache/ServerMPM                    | string    | Apache servers selected MPM                                         | :x:                | `"Not Found"` |
| /intel/apache/Server_Built                 | string    | Build time for server                                               | :x:                | `"Not Found"` |
| /intel/apache/CurrentTime                  | string    | Current server time                                                 | :x:                | `"Not Found"` |
| /intel/apache/RestartTime                  | string    | Last server restart time                                            | :x:                | `"Not Found"` |
| /intel/apache/ParentServerConfigGeneration | int       | # times apache has reread configs and restarted child processes     | :x:                | 0             |
| /intel/apache/ParentServerMPMGeneration    | int       | # times apache has reread MPM configs and restarted child processes | :x:                | 0             |
| /intel/apache/ServerUptime                 | string    | Server uptime in readable string                                    | :x:                | `"Not Found"` |
| /intel/apache/CPUUser                      | float64   | jiffs used in User Mode                                             | :x:                | 0             |
| /intel/apache/CPUSystem                    | float64   | jiffs used in System Mode                                           | :x:                | 0             |
| /intel/apache/CPUChildrenUser              | float64   | jiffs used in User Mode by child processes                          | :x:                | 0             |
| /intel/apache/CPUChildrenSystem            | float64   | jiffs used in System Mode by child processes                        | :x:                | 0             |
| /intel/apache/Load1                        | float64   | Server load over last 1 minute                                      | :x:                | 0             |
| /intel/apache/Load5                        | float64   | Server load over last 5 minutes                                     | :x:                | 0             |
| /intel/apache/Load15                       | float64   | Server load over last 15 minutes                                    | :x:                | 0             |
