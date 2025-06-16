package general

import "time"

//------------------------------------------------------------------------------
// Constant
//------------------------------------------------------------------------------
const def_check_interval = 1000 // 1sec
const (
	RST_OK         int = 1
	RST_UNEXIST    int = 2
	RST_ABNOMAL    int = 3
	MAX_WAIT_COUNT     = 5
)

//------------------------------------------------------------------------------
// struct
//------------------------------------------------------------------------------
type RunInfo struct {
	wdc       int16 // Watch dog counter
	prevWdc   int16 // 이전 Watch dog counter
	elapsedTm int   // process cycle time
	waitCount int8  // wait count
	maxEStm   int   // max elapsed Time
}

type RunBase struct {
	ID     int   // process ID
	THRID  int64 // routine ID
	Active bool

	runInfo RunInfo
	markTm  time.Time // loop start time
}

//------------------------------------------------------------------------------
// CheckRunInfo
//------------------------------------------------------------------------------
func (r *RunBase) checkRunInfo() bool {
	// check wdc
	r.runInfo.prevWdc = r.runInfo.wdc
	r.runInfo.wdc = 0

	if r.runInfo.prevWdc == 0 && r.runInfo.waitCount > MAX_WAIT_COUNT {
		return false
	}

	defer func() {
		r.runInfo.waitCount++
	}()

	return true
}

func (r *RunBase) MarkTime() {
	r.markTm = time.Now()
}

//------------------------------------------------------------------------------
// Register (process)
//------------------------------------------------------------------------------
func (r *RunBase) register(id int) {
	r.ID = id
}

//------------------------------------------------------------------------------
// Deregister (process)
//------------------------------------------------------------------------------
func (r *RunBase) deregister(id int) {
	if r.ID == id {
		r.ID = 0
		r.Active = false
	}
}

//------------------------------------------------------------------------------
// threadRegister (routine)
//------------------------------------------------------------------------------
func (r *RunBase) threadRegister(id int64) {
	r.THRID = id
}

//------------------------------------------------------------------------------
// threadDeregister (routine)
//------------------------------------------------------------------------------
func (r *RunBase) threadDeregister(id int64) {
	if r.THRID == id {
		r.THRID = 0
		r.Active = false
	}
}

//------------------------------------------------------------------------------
// UpdateRunInfo
//------------------------------------------------------------------------------
func (r *RunBase) UpdateRunInfo() {
	// update timer
	endTime := time.Now()
	elapsed := endTime.Sub(r.markTm)

	// nano sec to milli sec
	msTime := elapsed / 1000000
	r.runInfo.elapsedTm = int(msTime)

	if r.runInfo.elapsedTm > r.runInfo.maxEStm {
		r.runInfo.maxEStm = r.runInfo.elapsedTm
	}

	// update wdc
	r.runInfo.wdc++
	// reset waitcount
	r.runInfo.waitCount = 0
}
