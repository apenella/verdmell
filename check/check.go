/*
Package check is used by verdmell to manage the monitoring checks defined by user
*/
package check

import (
	"errors"
	"verdmell/sample"
	"verdmell/utils"
)

// Check stuct
// Check defines a check to be executed
type Check struct {
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Command        string      `json:"command"`
	Depend         []string    `json:"depend"`
	ExpirationTime int         `json:"expiration_time"`
	Interval       int         `json:"interval"`
	Custom         interface{} `json:"-"`
	Timeout        int         `json:"timeout"`
	// Timestamp
	Timestamp int64 `json:"timestamp"`
	//Queues
	TaskQueue chan *Check `json:"-"`
	//StatusChan chan int
	SampleChan chan *sample.CheckSample `json:"-"`
}

//
// ValidateCheck ensures that the Check has all the required data set. The method returns an error object once a definition method is found
func (c *Check) ValidateCheck() error {
	if c.Name == "" {
		return errors.New("(Check::ValidateCheck) Check requires a Name")
	}

	if c.Command == "" {
		return errors.New("(Check::ValidateCheck) Check '" + c.Name + "' requires a Command")
	}

	if c.ExpirationTime < 0 {
		err := errors.New("(Check::ValidateCheck) Check '" + c.Name + "' has an invalid expiration time")
		return err
	}

	if c.Interval < 0 {
		err := errors.New("(Check::ValidateCheck) Check '" + c.Name + "' has an invalid interval. Interval is lower than 0")
		return err
	}

	if c.Interval < c.Timeout {
		err := errors.New("(Check::ValidateCheck) Check '" + c.Name + "' has an invalid interval. Timeout should not be greater than interval")
		return err
	}

	return nil
}

//
// String method converts a Check to a string
func (c *Check) String() string {
	var str string
	var err error

	str, err = utils.ObjectToJSONString(c)
	if err != nil {
		return err.Error()
	}

	return str
}

//---------TODEL----------------------------------------------------------------

//
// StartQueue: method starts a queue for receive check
// func (c *Check) StartQueue() {
// 	env.Output.WriteChDebug("(Check::StartQueue) Starting queue for check '" + c.Name + "'")
// 	var err error
// 	expired := make(chan bool)
// 	result := -1
// 	queue := c.TaskQueue
// 	sample := new(sample.CheckSample)

// 	defer close(queue)

// 	//function to clean up the result to enforce that the check is not started while the sample is already valid
// 	sampleExpiration := func() {
// 		env.Output.WriteChDebug("(Check::StartQueue::sampleExpiration) Countdown for " + c.Name + "'s sample")
// 		timeout := time.After(time.Duration(c.ExpirationTime) * time.Second)
// 		for {
// 			select {
// 			case <-timeout:
// 				expired <- true
// 			}
// 		}
// 	}

// 	scheduleCheckTask := func() {
// 		env.Output.WriteChDebug("(Check::StartQueue::scheduleCheckTask) Scheduling a new task for '" + c.Name + "'. It will be launched " + strconv.Itoa(c.Interval) + "s latter.")
// 		timeout := time.After(time.Duration(c.Interval) * time.Second)
// 		for {
// 			select {
// 			case <-timeout:
// 				c.Timestamp = c.Timestamp + 1
// 				checkEngine := env.GetCheckEngine().(*CheckEngine)
// 				checkEngine.Start(c)
// 			}
// 		}
// 	}

// 	// waiting for task to be queued by EnqueueCheck
// 	for {
// 		select {
// 		case checkObj := <-queue:
// 			if result >= 0 {
// 				env.Output.WriteChDebug("(Check::StartQueue) ObjectTask alive and won't be started again. Check '" + checkObj.Name + "' already has a sample")
// 			} else {
// 				env.Output.WriteChDebug("(Check::StartQueue) ObjectTask started. Check '" + checkObj.Name + "' has no sample")
// 				if err, sample = checkObj.StartCheckTask(); err != nil {
// 					env.Output.WriteChWarn("(Check::StartQueue) Task for '" + checkObj.Name + "' has not finished properly")
// 				}
// 				result = sample.GetExit()
// 				env.Output.WriteChDebug("(Check::StartQueue) ObjectTask finished. Exit code for '" + checkObj.Name + "' is '" + strconv.Itoa(result) + "'")

// 				go sampleExpiration()
// 				go scheduleCheckTask()
// 			}
// 			//Send sample to check object sampleChan.
// 			checkObj.SampleChan <- sample
// 		case <-expired:
// 			env.Output.WriteChDebug("(Check::StartQueue) Sample for " + c.Name + " has expired")
// 			result = -1
// 		}
// 	}
// }

// //
// // EnqueueCheck: enqueu a Check to be run
// func (c *Check) EnqueueCheck() error {
// 	env.Output.WriteChDebug("(Check::EnqueueCheck) Enqueing Check '" + c.Name + "'")
// 	c.TaskQueue <- c

// 	return nil
// }

// //
// // StartCheckTask: executes the command defined on check an return the result
// func (c *Check) StartCheckTask() (error, *sample.CheckSample) {
// 	env.Output.WriteChDebug("(Check::StartCheckTask) Running a check '" + c.Name + "'")

// 	exit := 0
// 	var output string
// 	var messageError string
// 	//Exit codes
// 	// OK: 0
// 	// WARN: 1
// 	// ERROR: 2
// 	// UNKNOWN: other (-1)
// 	cmdSplitted := strings.SplitN(c.Command, " ", 2)
// 	time_init := time.Now()
// 	out, err := exec.Command(cmdSplitted[0], strings.Split(cmdSplitted[1], " ")...).Output()
// 	elapsedtime := time.Since(time_init)

// 	// When the exec has exit code, these code is achived. If is not possible to achive it, then is set to '-1', the unknown code.
// 	if err != nil {
// 		if exiterr, ok := err.(*exec.ExitError); ok {
// 			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
// 				exit = status.ExitStatus()
// 				env.Output.WriteChDebug("(Check::StartCheckTask) Exit status for '" + c.Name + "': " + strconv.Itoa(exit))
// 				if exit > 2 || exit < 0 {
// 					exit = -1
// 				}
// 			} else {
// 				exit = -1
// 			}
// 		} else {
// 			messageError = "(Check::StartCheckTask) The task for '" + c.Name + "has ended with errors"
// 			exit = -1
// 		}
// 	}

// 	if len(out) > 0 {
// 		output = string(out[:len(out)-1])
// 	}
// 	_, sample := c.GenerateCheckSample(exit, output, elapsedtime, time.Duration(c.ExpirationTime)*time.Second, c.Timestamp)

// 	if messageError != "" {
// 		env.Output.WriteChWarn(messageError)
// 		return errors.New(messageError), sample
// 	} else {
// 		return nil, sample
// 	}
// }

//
// GenerateCheckSample: method prepares the system to gather check's data
// func (c *Check) GenerateCheckSample(e int, o string, elapsedtime time.Duration, expirationtime time.Duration, timestamp int64) (error, *sample.CheckSample) {
// 	env.Output.WriteChDebug("(Check::GenerateCheckSample) CheckSample for '" + c.Name + "'")
// 	checkEngine := env.GetCheckEngine().(*CheckEngine)
// 	cs := new(sample.CheckSample)

// 	cs.SetCheck(c.Name)
// 	cs.SetExit(e)
// 	cs.SetOutput(o)
// 	cs.SetElapsedTime(elapsedtime)
// 	cs.SetSampletime(time.Now())
// 	cs.SetExpirationTime(expirationtime)
// 	cs.SetTimestamp(timestamp)

// 	env.Output.WriteChDebug("(Check::GenerateCheckSample) " + cs.String())

// 	// send the sample to CheckEngines's sendSample method to write its value into output channels
// 	checkEngine.sendSample(cs)

// 	return nil, cs
// }
