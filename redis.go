package main

import (
	"log"
	"strings"

	"github.com/tidwall/redcon"
)

func redisCommandParser(conn redcon.Conn, cmd redcon.Command) {
	switch strings.ToLower(string(cmd.Args[0])) {
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	case "ping":
		conn.WriteString("PONG")
	case "quit":
		conn.WriteString("OK")
		conn.Close()
	case "pfadd":
		if len(cmd.Args) < 3 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		counterName := cmd.Args[1]
		values := cmd.Args[2:]
		for _, v := range values {
			if err := localCounters.IncrementCounter(string(counterName), []byte(v)); err != nil {
				log.Println("Error incrementing counter ", string(counterName), err.Error())
				conn.WriteError("ERR: PFADD on " + err.Error())
				return
			}
			// TODO: profile as this will generate I/O - maybe cache
			//if cc, err = localCounters.RetrieveCounterEstimate(counterName); err != nil {
			//	log.Println("Error retrieving counter " + counterName + ":" + err)
			//	conn.WriteError("ERR: PFADD retrieving value " + err)
			//	return
			//}

		}
		// TODO: estimate if the cardinality was changed.
		conn.WriteInt(1)
		break
	case "pfcount":
		var err error
		var cc uint64
		if len(cmd.Args) < 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		counterName := cmd.Args[1]
		var keys []string
		for _, k := range cmd.Args[1:] {
			keys = append(keys, string(k))
		}

		if len(keys) > 1 {
			if cc, err = localCounters.RetrieveAndMergeCounterEstimates(keys...); err != nil {
				log.Println("Error retrieving counters " + string(counterName) + ":" + err.Error())
				conn.WriteError("ERR: pfcount retrieving and merging values: " + err.Error())
				return
			}
		} else {
			if cc, err = localCounters.RetrieveCounterEstimate(keys[0]); err != nil {
				log.Println("Error retrieving counters " + string(counterName) + ":" + err.Error())
				conn.WriteError("ERR: pfcount retrieving values: " + err.Error())
				return
			}

		}
		//conn.WriteNull()
		conn.WriteInt64(int64(cc))
		break
	}
}

func newConnection(conn redcon.Conn) bool {
	log.Printf("New connection: %s", conn.RemoteAddr())
	return true
}
func closeConnection(conn redcon.Conn, err error) {
	log.Printf("Connection closed: %s, err: %v", conn.RemoteAddr(), err)
}
