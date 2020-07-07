package main

import (
	"log"
	"strings"

	"github.com/tidwall/redcon"
)

/*
main redis parser
*/

func redisCommandParser(conn redcon.Conn, cmd redcon.Command) {
	switch strings.ToLower(string(cmd.Args[0])) {
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")

	case "ping":
		conn.WriteString("PONG")

	case "quit":
		conn.WriteString("OK")
		conn.Close()

	// TODO: implement basic set, get and del
	case "get":
		conn.WriteString("OK")
		conn.Close()

	case "set":
		conn.WriteString("OK")
		conn.Close()

	case "del":
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
			if err := localCounters.IncrementCounter(counterName, v); err != nil {
				log.Println("Error incrementing counter ", string(counterName), err.Error())
				conn.WriteError("ERR: PFADD on " + err.Error())
				return
			}
		}
		conn.WriteInt(1)
	case "pfcount":
		var err error
		var cc uint64
		if len(cmd.Args) < 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		counterName := cmd.Args[1]
		var keys [][]byte
		if len(cmd.Args[1:]) < 1 {
			conn.WriteInt64(0)
			return
		}
		for _, k := range cmd.Args[1:] {
			keys = append(keys, k)
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
		conn.WriteInt64(int64(cc))
	case "sadd":
		var err error
		if len(cmd.Args) < 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		setName := cmd.Args[1]
		member := cmd.Args[2]

		if err = localSets.SAdd(setName, member); err != nil {
			log.Printf("Error adding member %s to set %s: %s\n", string(member), string(setName), err.Error())
			conn.WriteError("ERR: sadd " + err.Error())
			return
		}
		conn.WriteInt64(1)

	case "sismember":
		var err error
		var ok bool
		if len(cmd.Args) < 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		setName := cmd.Args[1]
		member := cmd.Args[2]

		if ok, err = localSets.SisMember(setName, member); err != nil {
			log.Printf("Error looking up membership of %s to set %s: %s\n", string(member), string(setName), err.Error())
			conn.WriteError("ERR: ismember " + err.Error())
			return
		}
		if ok {
			conn.WriteInt64(1)
			return
		}
		conn.WriteInt64(0)

	case "srem":
		var err error
		var ok bool
		if len(cmd.Args) < 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		setName := cmd.Args[1]
		member := cmd.Args[2]

		if ok, err = localSets.SRem(setName, member); err != nil {
			log.Printf("Error removing %s from set %s: %s\n", string(member), string(setName), err.Error())
			conn.WriteError("ERR: srem " + err.Error())
			return
		}

		if ok {
			conn.WriteInt64(1)
			return
		}
		conn.WriteInt64(0)

	case "scard":
		var err error
		var cardinality uint
		if len(cmd.Args) < 1 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		setName := cmd.Args[1]

		if cardinality, err = localSets.SCard(setName); err != nil {
			log.Printf("Error fetching cardinality for %s: %s\n", string(setName), err.Error())
			conn.WriteError("ERR: srem " + err.Error())
			return
		}
		conn.WriteInt64(int64(cardinality))
	}
}

func newConnection(conn redcon.Conn) bool {
	log.Printf("New connection: %s", conn.RemoteAddr())
	return true
}
func closeConnection(conn redcon.Conn, err error) {
	log.Printf("Connection closed: %s, err: %v", conn.RemoteAddr(), err)
}
