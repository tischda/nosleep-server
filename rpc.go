package main

import "log"

type ShutdownArgs struct{}
type ShutdownReply struct{}

type SleepControl struct {
	shutdown chan bool
}

type ReadFlagsReply struct {
	Flags uint32
}

func (c *SleepControl) Clear(args *struct{}, reply *struct{}) error {
	log.Println("Clear RPC called — clearing sleep flags")
	ClearSleepFlags()
	return nil
}

func (c *SleepControl) Display(args *struct{}, reply *struct{}) error {
	log.Println("Display RPC called — forcing display on")
	ForceDisplayOn()
	return nil
}

func (c *SleepControl) System(args *struct{}, reply *struct{}) error {
	log.Println("System RPC called — forcing system on")
	ForceSystemOn()
	return nil
}

func (c *SleepControl) Critical(args *struct{}, reply *struct{}) error {
	log.Println("Critical RPC called — forcing system critical on")
	ForceSystemCriticalOn()
	return nil
}

func (c *SleepControl) Read(args *struct{}, reply *ReadFlagsReply) error {
	log.Println("Read RPC called — reading current flags")
	reply.Flags = ReadFlags()
	return nil
}

func (c *SleepControl) Shutdown(args *struct{}, reply *struct{}) error {
	log.Println("Shutdown RPC called — shutting down server")
	close(c.shutdown)
	return nil
}
