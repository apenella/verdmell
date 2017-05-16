/*
Sample System management

The package 'sample' is used by verdmell to manage the samples generated by monitoring checks

-SampleEngine
-CheckSamples
*/
package sample

 import (
  "errors"
  "sync"
  "strconv"
  "verdmell/environment"
  "verdmell/utils"
 )

//
var env *environment.Environment
//#
//#
//# SampleEngine struct:
//# SampleEngine defines a map to store the maps
type SampleEngine struct{
  Samples map[string]*CheckSampleSync `json: "samples"`
  inputChannel chan interface{} `json:"-"`
}

type CheckSampleSync struct {
	Sample *CheckSample 
	mutex sync.RWMutex
}

//
//# NewSampleEngine: method prepare to system gather information
func NewSampleEngine(e *environment.Environment) (error, *SampleEngine) {
  e.Output.WriteChDebug("(SampleEngine::NewSampleEngine)")
	sys := &SampleEngine{
    Samples: make(map[string]*CheckSampleSync),
  }
  //var err error
	env = e

  // start the sample receiver
  sys.Start()

  // Set the environment's sample engine
  env.SetSampleEngine(sys)

  env.Output.WriteChInfo("(SampleEngine::NewSampleEngine) Hi! I'm your new sample engine instance")

	return nil, sys
}

//
//# SetInputChannel: methods sets the inputChannel's value
func (s *SampleEngine) SetInputChannel(c chan interface{}) {
  s.inputChannel = c
}

//
//# GetInputChannel: methods sets the inputChannel's value
func (s *SampleEngine) GetInputChannel() chan interface{} {
  return s.inputChannel
}

//#
//# Specific methods
//#----------------------------------------------------------------------------------------

//
//# SayHi: 
func (sys *SampleEngine) SayHi() {
  env.Output.WriteChInfo("(SampleEngine::SayHi) Hi! I'm your new sample engine instance")
}
//
//# StartServiceEngine: method prepares the system to wait sample and calculate the results for services
func (s *SampleEngine) Start() error {
  env.Output.WriteChDebug("(SampleEngine::Start) Starting sample receiver")
  s.inputChannel = make(chan interface{})

  go func() {
    defer close (s.inputChannel)
    for{
      select{
      case obj := <-s.inputChannel:
        sample := obj.(*CheckSample)
        env.Output.WriteChDebug("(SampleEngine::Start) New sample received for '"+sample.GetCheck()+"'")
        s.AddSample(sample)
      }
    }
  }()
  return nil
}
//
//# SendSample: method prepares the system to wait samples
func (s *SampleEngine) SendSample(sample *CheckSample) {
  env.Output.WriteChDebug("(SampleEngine::SendSample) Send sample "+sample.String())
  s.inputChannel <- sample
}
//
// AddSample method creaty a new entry to CheckSample or modify its value
func (sys *SampleEngine) AddSample(cs *CheckSample) error {
  env.Output.WriteChDebug("(SampleEngine::AddSample) ["+strconv.Itoa(int(cs.GetTimestamp()))+"] '"+cs.GetCheck()+"'")
  var sam *CheckSampleSync
  var exist bool

  name := cs.GetCheck()
  
  //If now sample exist for this check, initialize it
  if _, exist = sys.Samples[name]; !exist{
  	sam = new(CheckSampleSync)
  	sys.Samples[name] = sam
  } else {
    sam = sys.Samples[name]
  }

  //write lock
  sam.mutex.Lock()
	defer sam.mutex.Unlock()
	sam.Sample = cs

  return nil
}
//
//# GetSample: method returns the CheckSample object for a CheckObject
func (sys *SampleEngine) GetSample(name string) (error, *CheckSample) {
  env.Output.WriteChDebug("(SampleEngine::GetSample) '"+name+"'")
  var sam *CheckSampleSync
  var exist bool

  // if no sample for the check, an error is thrown
  if sam, exist = sys.Samples[name]; !exist {
    msg := "(SampleEngine::GetSample) There is not a sample for the check '"+name+"'"
    env.Output.WriteChDebug(msg)
    return errors.New(msg),nil
  }

  sam = sys.Samples[name]

  //read lock
  sam.mutex.RLock()
  defer sam.mutex.RUnlock()

  return nil,sam.Sample
}
//
//# DeleteSample: method deletes the Sample for a CheckObject
func (sys *SampleEngine) DeleteSample(name string) error {
  var sam *CheckSampleSync
  var exist bool

  // if no sample for the check, an error is thrown
  if sam, exist = sys.Samples[name]; !exist {
    msg := "(SampleEngine::GetSample) There is not a sample to be deleled for the check '"+name+"'"
    env.Output.WriteChDebug(msg)
    return errors.New(msg)
  }

  sam.mutex.Lock()
  defer sam.mutex.Unlock()
  delete(sys.Samples,name)
  return nil
}
//
//# GetAllSamples: return the status of all checks
func (sys *SampleEngine) GetAllSamples() (error, []byte) {
  env.Output.WriteChDebug("(SampleEngine::GetAllSamples)")
  sample := make(map[string] *CheckSample)

  for name, obj := range sys.Samples {
    sample[name] = obj.Sample
  }

  return utils.ObjectToJsonByte(sample) 
}
//
//# GetSampleForCheck: return the status of all checks
func (sys *SampleEngine) GetSampleForCheck(name string) (error, []byte) {
  env.Output.WriteChDebug("(SampleEngine::GetSampleForCheck)")
  var sample *CheckSample
  var err error

  if err,sample = sys.GetSample(name); err!=nil{
    msg := "(SampleEngine::GetSampleForCheck) There is not a sample for the check '"+name+"'"
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  }

  return utils.ObjectToJsonByte(sample)
}


//#
//# Common methods
//#---------------------------------------------------------------------

//
//# String: converts a SampleEngine object to string
func (sys *SampleEngine) String() string {
  if err, str := utils.ObjectToJsonString(sys); err != nil{
    return err.Error()
  } else{
    return str
  }
}

//
//# Itoa: method transform a integer to status string
func Itoa(i int) string {
  switch(i){
  case 0:
    return "OK"
  case 1:
    return "WARN"
  case 2:
    return "ERROR"
  default:
    return "UNKNOWN"
  }
}

//#######################################################################################################