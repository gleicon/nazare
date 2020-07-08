package server

import (
	"log"
	"strings"

	"github.com/tidwall/redcon"
)

/*
redis protocol parser
*/

func (nzs *NZServer) redisCommandParser(conn redcon.Conn, cmd redcon.Command) {
	redisCmd := strings.ToLower(string(cmd.Args[0]))

	switch redisCmd {
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")

	case "ping":
		conn.WriteString("PONG")

	case "quit":
		conn.WriteString("OK")
		conn.Close()

	case "get":
		var err error
		var val []byte

		if len(cmd.Args) < 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		key := cmd.Args[1]

		if val, err = nzs.localDatastorage.Get(key); err != nil {
			log.Println("Error getting data from ", string(key), err.Error())
			conn.WriteError("ERR: GET - " + err.Error())
			return
		}
		conn.WriteString(string(val))

	case "set":
		var err error

		if len(cmd.Args) < 3 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		key := cmd.Args[1]
		val := cmd.Args[2]

		if err = nzs.localDatastorage.Add(key, val); err != nil {
			log.Println("Error setting data: ", string(key), err.Error())
			conn.WriteError("ERR: GET - " + err.Error())
			return
		}
		conn.WriteString("OK")

	case "del":
		var err error
		var count int64

		if len(cmd.Args) < 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		for _, key := range cmd.Args[1:] {
			var ok bool
			if ok, err = nzs.localDatastorage.Delete(key); err != nil {
				log.Println("Error deleting data from ", string(key), err.Error())
				conn.WriteError("ERR: GET - " + err.Error())
				return
			}
			if ok {
				count++
			}
		}

		conn.WriteInt64(count)

	case "pfadd":
		if len(cmd.Args) < 3 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		counterName := cmd.Args[1]
		values := cmd.Args[2:]
		for _, v := range values {
			if err := nzs.localCounters.IncrementCounter(counterName, v); err != nil {
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
			if cc, err = nzs.localCounters.RetrieveAndMergeCounterEstimates(keys...); err != nil {
				log.Println("Error retrieving counters " + string(counterName) + ":" + err.Error())
				conn.WriteError("ERR: pfcount retrieving and merging values: " + err.Error())
				return
			}
		} else {
			if cc, err = nzs.localCounters.RetrieveCounterEstimate(keys[0]); err != nil {
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

		if err = nzs.localSets.SAdd(setName, member); err != nil {
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

		if ok, err = nzs.localSets.SisMember(setName, member); err != nil {
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

		if ok, err = nzs.localSets.SRem(setName, member); err != nil {
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

		if cardinality, err = nzs.localSets.SCard(setName); err != nil {
			log.Printf("Error fetching cardinality for %s: %s\n", string(setName), err.Error())
			conn.WriteError("ERR: srem " + err.Error())
			return
		}
		conn.WriteInt64(int64(cardinality))
	}
}

func (nzs *NZServer) newConnection(conn redcon.Conn) bool {
	log.Printf("New connection: %s", conn.RemoteAddr())
	return true
}
func (nzs *NZServer) closeConnection(conn redcon.Conn, err error) {
	log.Printf("Connection closed: %s, err: %v", conn.RemoteAddr(), err)
}
