package am

import (
	"database/sql"
	"encoding/json"
)

type DBHandler interface {
	Ping() // 연결확인
}

const (
	ADMINISTRATOR = 1
	OPERATOR      = 2
	MONITOR       = 3
	SUPERUSER     = 99
)

// select count(*) 삭제하는 대신, reqRow + N개로 검색해서 다음 데이터 있는지 확인
const (
	SEARCH_MORE_ROW = 1
)

// Nullable String that overrides sql.NullString
type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal("")
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

// ---------------------------------------------------------------------------
// DBGroup
// ---------------------------------------------------------------------------
type DBGroup struct {
	GRP_ID          int        `json:"id"`
	GRP_NM          NullString `json:"name"`
	GRP_NUM         int        `json:"gNum"`
	GRP_DEFCTRLMODE int        `json:"defMode"`
	AUTO_ONLINEYN   int        `json:"autoOnlineYN"`
	GRP_LAT         float64    `json:"lat"`
	GRP_LON         float64    `json:"lon"`
	GRP_ZOMMLV      float64    `json:"zoomlvl"`
	MIN_CYCLE       int        `json:"minCyc"`
	MAX_CYCLE       int        `json:"maxCyc"`
	TRANS_MODE      int        `json:"tranMode"`
	USE_YN          int        `json:"useYn"`
	DEL_YN          int        `json:"delYn"`
	R28GRP_YN       int        `json:"r28Yn"`
}

type QRY_DBGroup struct {
	GRP_ID          string `json:"id"`
	GRP_NM          string `json:"name"`
	GRP_NUM         string `json:"gNum"`
	GRP_DEFCTRLMODE string `json:"defMode"`
	AUTO_ONLINEYN   string `json:"autoOnlineYN"`
	GRP_LAT         string `json:"lat"`
	GRP_LON         string `json:"lon"`
	GRP_ZOMMLV      string `json:"zoomlvl"`
	MIN_CYCLE       string `json:"minCyc"`
	MAX_CYCLE       string `json:"maxCyc"`
	TRANS_MODE      string `json:"tranMode"`
	USE_YN          string `json:"useYn"`
	DEL_YN          string `json:"delYn"`
	R28GRP_YN       string `json:"r28Yn"`
}

// ---------------------------------------------------------------------------
// DBLocal
// ---------------------------------------------------------------------------
type DBLocal struct {
	LOC_ID    int        `json:"id"`
	LOC_NM    NullString `json:"name"`
	LOC_NUM   int        `json:"number"`
	GRP_ID    int        `json:"gId"`
	GRP_NUM   int        `json:"gNum"`
	LOC_TYPE  int        `json:"locType"`
	LC_TYPE   int        `json:"lcType"`
	LC_TP     int        `json:"lcTp"`
	LAMP_TYPE int        `json:"lamp"`

	DELTA_LIMIT         int        `json:"delta"`
	TRANSCYCLE          int        `json:"trsCyc"`
	AUTOERR_CORR        int        `json:"aCorrect"`
	AUTOONLINE_YN       int        `json:"aOnYn"`
	COMMTYPE            int        `json:"comType"`
	IPADDR              NullString `json:"ip"`
	NET_PORT            int        `json:"port"`
	OFFLINE_LOCK        int        `json:"offLock"`
	MAX_TRANSCYCLE      int        `json:"maxTrsCyc"`
	OFFSET_FINE_TUNE_YN int        `json:"offsetTuneYn"`
	PHASE_OPR_MODE      int        `json:"phaseOprMode"`

	NODE_ID NullString `json:"nodeId"`
	NODELAT float64    `json:"-"`
	NODELON float64    `json:"-"`
	GEOM    string     `json:"-"`

	LOC_ROT_ANG float64 `json:"angle"`

	USE_YN int `json:"useYn"`
	DEL_YN int `json:"delYn"`

	SCM_ID     int `json:"scmId"`
	MST_LOC_ID int `json:"mstLocId"`

	OPEN_STATE_INVS int        `json:"inversedoor"`
	MANUF_DT        NullString `json:"manufdt"`
	MEMO            NullString `json:"comment"`

	PPC_CTRL_TYPE int        `json:"ppcusetp"`
	PPC_IP        NullString `json:"ppcip"`
}

type QRY_DBLocal struct {
	LOC_ID    string `json:"id"`
	LOC_NM    string `json:"name"`
	LOC_NUM   string `json:"number"`
	GRP_ID    string `json:"gId"`
	GRP_NUM   string `json:"gNum"`
	LOC_TYPE  string `json:"locType"`
	LC_TYPE   string `json:"lcType"`
	LC_TP     string `json:"lcTp"`
	LAMP_TYPE string `json:"lamp"`

	DELTA_LIMIT         string `json:"delta"`
	TRANSCYCLE          string `json:"trsCyc"`
	AUTOERR_CORR        string `json:"aCorrect"`
	AUTOONLINE_YN       string `json:"aOnYn"`
	COMMTYPE            string `json:"comType"`
	IPADDR              string `json:"ip"`
	NET_PORT            string `json:"port"`
	OFFLINE_LOCK        string `json:"offLock"`
	MAX_TRANSCYCLE      string `json:"maxTrsCyc"`
	OFFSET_FINE_TUNE_YN string `json:"offsetTuneYn"`
	PHASE_OPR_MODE      string `json:"phaseOprMode"`

	NODE_ID string `json:"nodeId"`
	NODELAT string `json:"-"`
	NODELON string `json:"-"`
	GEOM    string `json:"-"`

	LOC_ROT_ANG     string `json:"angle"`
	IMG_COLL_SVR_ID string

	USE_YN string `json:"useYn"`
	DEL_YN string `json:"delYn"`

	SCM_ID     string `json:"scmId"`
	MST_LOC_ID string `json:"mstLocId"`

	OPEN_STATE_INVS string `json:"inversedoor"`
	MANUF_DT        string `json:"manufdt"`
	MEMO            string `json:"comment"`

	PPC_CTRL_TYPE string `json:"ppcusetp"`
	PPC_IP        string `json:"ppcip"`
}

// ---------------------------------------------------------------------------
// REPORTLocal
// ---------------------------------------------------------------------------
type REPORTLocal struct {
	LOC_ID  int        `json:"id"`
	LOC_NM  string     `json:"-"`
	LOC_NUM int        `json:"-"`
	GRP_NUM int        `json:"gnum"`
	GRP_NM  NullString `json:"gname"`

	GRP_DEFCTRLMODE int `json:"defMode"`
	LC_TYPE         int `json:"lcType"`
	MAINPHASE       int `json:"mainphase"`
}

type QRY_REPORTLocal struct {
	LOC_ID  string `json:"id"`
	LOC_NM  string `json:"-"`
	LOC_NUM string `json:"-"`
	GRP_NUM string `json:"gnum"`
	GRP_NM  string `json:"gname"`

	GRP_DEFCTRLMODE string `json:"defMode"`
	LC_TYPE         string `json:"lcType"`
	MAINPHASE       string `json:"mainphase"`
}

// ---------------------------------------------------------------------------
// DBGroupLocal
// ---------------------------------------------------------------------------
type DBGroupLocal struct {
	LOC_ID    int    `json:"id"`
	LOC_NM    string `json:"name"`
	LOC_NUM   int    `json:"number"`
	GRP_ID    int    `json:"gId"`
	GRP_NUM   int    `json:"gNum"`
	LOC_TYPE  int    `json:"locType"`
	LC_TYPE   int    `json:"lcType"`
	LAMP_TYPE int    `json:"lamp"`

	DELTA_LIMIT         int    `json:"delta"`
	TRANSCYCLE          int    `json:"trsCyc"`
	AUTOERR_CORR        int    `json:"aCorrect"`
	AUTOONLINE_YN       int    `json:"aOnYn"`
	COMMTYPE            int    `json:"comType"`
	IPADDR              string `json:"ip"`
	NET_PORT            int    `json:"port"`
	OFFLINE_LOCK        int    `json:"offLock"`
	MAX_TRANSCYCLE      int    `json:"maxTrsCyc"`
	OFFSET_FINE_TUNE_YN int    `json:"offsetTuneYn"`
	PHASE_OPR_MODE      int    `json:"phaseOprMode"`

	NODE_ID string  `json:"nodeId"`
	NODELAT float64 `json:"-"`
	NODELON float64 `json:"-"`
	GEOM    string  `json:"GEOM"`

	LOC_ROT_ANG float64 `json:"angle"`

	USE_YN int `json:"useYn"`
	DEL_YN int `json:"delYn"`
}

// ---------------------------------------------------------------------------
// DBGRP_LOC_LANE
// ---------------------------------------------------------------------------
type DBGrpLocLane struct {
	GRP_ID      int    `json:"id"`
	GRP_NM      string `json:"name"`
	GRP_NUM     int    `json:"gNum"`
	ST_LOC_ID   int    `json:"stLocId"`
	END_LOC_ID  int    `json:"endLocId"`
	ST_NODE_ID  string `json:"stNodeId"`
	END_NODE_ID string `json:"endNodeId"`
	ST_GEOM     string `json:"-"`
	END_GEOM    string `json:"-"`
}

// ---------------------------------------------------------------------------
// DBLOC_PHASE
// ---------------------------------------------------------------------------
type DBLOC_PHASE struct {
	LOC_ID      int `json:"locId"`
	LOC_NUM     int `json:"number"`
	BANK_ID     int `json:"bankId"`
	RING_NUM    int `json:"ringNum"`
	PHASE_NUM   int `json:"phaseNum"`
	MIN_GREEN   int `json:"min"`
	MAX_GREEN   int `json:"max"`
	YELLOW      int `json:"yellow"`
	FIXPHASE_YN int `json:"fixPhaseYn"`
	FLOW_NUM    int `json:"flowNum"`
}

type QRY_DBLOC_PHASE struct {
	LOC_ID      string `json:"locId"`
	LOC_NUM     string `json:"number"`
	BANK_ID     string `json:"bankId"`
	RING_NUM    string `json:"ringNum"`
	PHASE_NUM   string `json:"phaseNum"`
	MIN_GREEN   string `json:"min"`
	MAX_GREEN   string `json:"max"`
	YELLOW      string `json:"yellow"`
	FIXPHASE_YN string `json:"fixPhaseYn"`
	FLOW_NUM    string `json:"flowNum"`
	Aphase      string `json:"-"`
	Bphase      string `json:"-"`
	AFlow       string `json:"-"`
	BFlow       string `json:"-"`
}

// ---------------------------------------------------------------------------
// REPORTLOC_PHASE
// ---------------------------------------------------------------------------
type REPORTLOC_PHASE struct {
	LOC_ID      int `json:"-"`
	LOC_NUM     int `json:"-"`
	BANK_ID     int `json:"-"`
	RING_NUM    int `json:"ringNum"`
	PHASE_NUM   int `json:"phaseNum"`
	MIN_GREEN   int `json:"min"`
	MAX_GREEN   int `json:"max"`
	YELLOW      int `json:"yellow"`
	FIXPHASE_YN int `json:"-"`
	FLOW_NUM    int `json:"flowNum"`

	BeforePed int `json:"beforePed"`
	PedFlash  int `json:"pedFlash"`
	PedGreen  int `json:"pedGreen"`
}

type QRY_REPORTLOC_PHASE struct {
	LOC_ID      string `json:"-"`
	LOC_NUM     string `json:"-"`
	BANK_ID     string `json:"-"`
	RING_NUM    string `json:"ringNum"`
	PHASE_NUM   string `json:"phaseNum"`
	MIN_GREEN   string `json:"min"`
	MAX_GREEN   string `json:"max"`
	YELLOW      string `json:"yellow"`
	FIXPHASE_YN string `json:"-"`
	FLOW_NUM    string `json:"flowNum"`

	BeforePed string `json:"beforePed"`
	PedFlash  string `json:"pedFlash"`
	PedGreen  string `json:"pedGreen"`
}

// ---------------------------------------------------------------------------
// REPORTLOC_PHASE2
// ---------------------------------------------------------------------------
type REPORTLOC_PHASE2 struct {
	LOC_ID     int `json:"-"`
	LOC_NUM    int `json:"-"`
	BANK_ID    int `json:"-"`
	PHASE_NUM  int `json:"phaseNum"`
	MIN_GREENA int `json:"amin"`
	MAX_GREENA int `json:"amax"`
	YELLOWA    int `json:"ayellow"`
	FLOW_NUMA  int `json:"aflowNum"`

	MIN_GREENB int `json:"bmin"`
	MAX_GREENB int `json:"bmax"`
	YELLOWB    int `json:"byellow"`
	FLOW_NUMB  int `json:"bflowNum"`
}

type QRY_REPORTLOC_PHASE2 struct {
	LOC_ID     string `json:"-"`
	LOC_NUM    string `json:"-"`
	BANK_ID    string `json:"-"`
	PHASE_NUM  string `json:"phaseNum"`
	MIN_GREENA string `json:"amin"`
	MAX_GREENA string `json:"amax"`
	YELLOWA    string `json:"ayellow"`
	FLOW_NUMA  string `json:"aflowNum"`

	MIN_GREENB string `json:"bmin"`
	MAX_GREENB string `json:"bmax"`
	YELLOWB    string `json:"byellow"`
	FLOW_NUMB  string `json:"bflowNum"`
}

type DBLocalPhasePlan struct {
	LOC_ID            int `json:"locId"`
	LOC_NUM           int `json:"number"`
	LOC_PLANNUM       int `json:"plan"`
	LOC_PHASEPLAN_IDX int `json:"phaseIdx"`
	RINGA_1PHASE      int `json:"-"`
	RINGA_2PHASE      int `json:"-"`
	RINGA_3PHASE      int `json:"-"`
	RINGA_4PHASE      int `json:"-"`
	RINGA_5PHASE      int `json:"-"`
	RINGA_6PHASE      int `json:"-"`
	RINGA_7PHASE      int `json:"-"`
	RINGA_8PHASE      int `json:"-"`
	RINGB_1PHASE      int `json:"-"`
	RINGB_2PHASE      int `json:"-"`
	RINGB_3PHASE      int `json:"-"`
	RINGB_4PHASE      int `json:"-"`
	RINGB_5PHASE      int `json:"-"`
	RINGB_6PHASE      int `json:"-"`
	RINGB_7PHASE      int `json:"-"`
	RINGB_8PHASE      int `json:"-"`

	OFFSET   int    `json:"offset"`
	Aphases  [8]int `json:"aPhases"`
	Bphases  [8]int `json:"bPhases"`
	CYCLE_LV int    `json:"cyclLv"`
	CYCLELEN int    `json:"cycLen"`
}

type QRY_DBLocalPhasePlan struct {
	LOC_ID            string `json:"locId"`
	LOC_NUM           string `json:"number"`
	LOC_PLANNUM       string `json:"plan"`
	LOC_PHASEPLAN_IDX string `json:"phaseIdx"`
	RINGA_1PHASE      string `json:"-"`
	RINGA_2PHASE      string `json:"-"`
	RINGA_3PHASE      string `json:"-"`
	RINGA_4PHASE      string `json:"-"`
	RINGA_5PHASE      string `json:"-"`
	RINGA_6PHASE      string `json:"-"`
	RINGA_7PHASE      string `json:"-"`
	RINGA_8PHASE      string `json:"-"`
	RINGB_1PHASE      string `json:"-"`
	RINGB_2PHASE      string `json:"-"`
	RINGB_3PHASE      string `json:"-"`
	RINGB_4PHASE      string `json:"-"`
	RINGB_5PHASE      string `json:"-"`
	RINGB_6PHASE      string `json:"-"`
	RINGB_7PHASE      string `json:"-"`
	RINGB_8PHASE      string `json:"-"`

	OFFSET   string `json:"offset"`
	Aphases  string `json:"aPhases"`
	Bphases  string `json:"bPhases"`
	CYCLE_LV string `json:"cyclLv"`
	CYCLELEN string `json:"cycLen"`
}

// ---------------------------------------------------------------------------
// DBStartUpCode
// ---------------------------------------------------------------------------
type DBStartUpCode struct {
	LOC_ID  int        `json:"id"`
	LOC_NM  NullString `json:"name"`
	LOC_NUM int        `json:"number"`
	Data    []int      `json:"data"`
	//Data []int8 `json:"data"`
}

type QRY_DBStartUpCode struct {
	LOC_ID  string `json:"id"`
	LOC_NM  string `json:"name"`
	LOC_NUM string `json:"number"`
	Data    string `json:"data"`
	//Data string `json:"data"`
}

// ---------------------------------------------------------------------------
// DBlocalFuncTable
// ---------------------------------------------------------------------------
type DBlocalFuncTable struct {
	LOC_ID      int `json:"locId"`
	LOC_NUM     int `json:"number"`
	CMD_NUM     int `json:"num"`
	MON         int `json:"month"`
	DAY         int `json:"day"`
	DAYWEEK     int `json:"dayWeek"`
	ST_HOUR     int `json:"stHour"`
	ST_MIN      int `json:"stMin"`
	END_HOUR    int `json:"endHour"`
	END_MIN     int `json:"endMin"`
	FUNC_TYPE   int `json:"funcType"`
	LOC_PLANNUM int `json:"sicaPlan"`
}

type QRY_DBlocalFuncTable struct {
	LOC_ID      string `json:"locId"`
	LOC_NUM     string `json:"number"`
	CMD_NUM     string `json:"num"`
	MON         string `json:"month"`
	DAY         string `json:"day"`
	DAYWEEK     string `json:"dayWeek"`
	ST_HOUR     string `json:"stHour"`
	ST_MIN      string `json:"stMin"`
	END_HOUR    string `json:"endHour"`
	END_MIN     string `json:"endMin"`
	FUNC_TYPE   string `json:"funcType"`
	LOC_PLANNUM string `json:"sicaPlan"`
}

// ---------------------------------------------------------------------------
// REPORTlocalFuncTable
// ---------------------------------------------------------------------------
type REPORTlocalFuncTable struct {
	LOC_ID      int `json:"-"`
	LOC_NUM     int `json:"-"`
	CMD_NUM     int `json:"num"`
	MON         int `json:"month"`
	DAY         int `json:"day"`
	DAYWEEK     int `json:"dayWeek"`
	ST_HOUR     int `json:"stHour"`
	ST_MIN      int `json:"stMin"`
	END_HOUR    int `json:"endHour"`
	END_MIN     int `json:"endMin"`
	FUNC_TYPE   int `json:"funcType"`
	LOC_PLANNUM int `json:"sicaPlan"`
}

type QRY_REPORTlocalFuncTable struct {
	LOC_ID      string `json:"-"`
	LOC_NUM     string `json:"-"`
	CMD_NUM     string `json:"num"`
	MON         string `json:"month"`
	DAY         string `json:"day"`
	DAYWEEK     string `json:"dayWeek"`
	ST_HOUR     string `json:"stHour"`
	ST_MIN      string `json:"stMin"`
	END_HOUR    string `json:"endHour"`
	END_MIN     string `json:"endMin"`
	FUNC_TYPE   string `json:"funcType"`
	LOC_PLANNUM string `json:"sicaPlan"`
}

// ---------------------------------------------------------------------------
// DBSignalmap
// ---------------------------------------------------------------------------
type DBSignalmap struct {
	LOC_ID    int     `json:"-"`
	LC_TYPE   int     `json:"-"`
	LAMP_TYPE int     `json:"-"`
	BANK_ID   int     `json:"-"`
	Data      []uint8 `json:"-"`
}

type QRY_DBSignalmap struct {
	LOC_ID    string `json:"-"`
	LC_TYPE   string `json:"-"`
	LAMP_TYPE string `json:"-"`
	BANK_ID   string `json:"-"`
	Data      string `json:"-"`
}

// ---------------------------------------------------------------------------
// DBflashmap
// ---------------------------------------------------------------------------
type DBflashmap struct {
	LOC_ID    int   `json:"locId"`
	LOC_NUM   int   `json:"number"`
	LAMP_TYPE int   `json:"lamp"`
	Data      []int `json:"data"`
}

type QRY_DBflashmap struct {
	LOC_ID    string `json:"locId"`
	LOC_NUM   string `json:"number"`
	LAMP_TYPE string `json:"lamp"`
	Data      string `json:"data"`
}

// ---------------------------------------------------------------------------
// DisSigoutput
// ---------------------------------------------------------------------------
type DBDisSigoutput struct {
	LOC_ID  int        `json:"id"`
	LOC_NM  NullString `json:"name"`
	LOC_NUM int        `json:"number"`
	//Data    [119]uint8 `json:"data"`
	//Data []uint8 `json:"data"` // R26부터 175, R26전버전 119
	Data []int `json:"data"` // R26부터 175, R26전버전 119
}

type QRY_DBDisSigoutput struct {
	LOC_ID  string `json:"id"`
	LOC_NM  string `json:"name"`
	LOC_NUM string `json:"number"`
	//Data    [119]uint8 `json:"data"`
	//Data []uint8 `json:"data"` // R26부터 175, R26전버전 119
	Data string `json:"data"` // R26부터 175, R26전버전 119
}

// ---------------------------------------------------------------------------
// DBConflictMap
// ---------------------------------------------------------------------------
type DBConflictMap struct {
	LOC_ID  int        `json:"id"`
	LOC_NM  NullString `json:"name"`
	LOC_NUM int        `json:"number"`
	Data    [128]uint8 `json:"data"`
}

type QRY_DBConflictMap struct {
	LOC_ID  string `json:"id"`
	LOC_NM  string `json:"name"`
	LOC_NUM string `json:"number"`
	Data    string `json:"data"`
}

// ---------------------------------------------------------------------------
// DBDetInfo
// ---------------------------------------------------------------------------
type DBDetInfo struct {
	LOC_ID           int `json:"id"`
	DET_NUM          int `json:"detNum"`
	DET_TYPE         int `json:"detType"`
	FLOW_NUM         int `json:"flow"`
	IN_CONTROL       int `json:"inCtrl"`
	SATUR_WEIGHT     int `json:"saturW"`
	SPEED_WEIGHT     int `json:"speedW"`
	TRFFICVOL_WEIGHT int `json:"trffW"`

	LOOP_TYPE int `json:"loopType"`
	DIR       int `json:"dir"`
	POS       int `json:"pos"`
	LANE      int `json:"lane"`
	DET_USE   int `json:"detUse"`
	RING      int `json:"ring"`
	PHASE     int `json:"phase"`
}

type QRY_DBDetInfo struct {
	LOC_ID           string `json:"id"`
	DET_NUM          string `json:"detNum"`
	DET_TYPE         string `json:"detType"`
	FLOW_NUM         string `json:"flow"`
	IN_CONTROL       string `json:"inCtrl"`
	SATUR_WEIGHT     string `json:"saturW"`
	SPEED_WEIGHT     string `json:"speedW"`
	TRFFICVOL_WEIGHT string `json:"trffW"`

	LOOP_TYPE string `json:"loopType"`
	DIR       string `json:"dir"`
	POS       string `json:"pos"`
	LANE      string `json:"lane"`
	DET_USE   string `json:"detUse"`
	RING      string `json:"ring"`
	PHASE     string `json:"phase"`
}

// ---------------------------------------------------------------------------
// DBImgDetInfo
// ---------------------------------------------------------------------------
type DBImgDetInfo struct {
	LOC_ID               int    `json:"id"`
	CAM_ID               string `json:"cam"`
	CHAN                 int    `json:"chan"`
	LANE_NUM             int    `json:"lane"`
	ACCROAD_DIR          int    `json:"dir"`
	COMM_TYPE            int    `json:"net"`
	RING_NUM             int    `json:"ring"`
	PHASE_NUM            int    `json:"phase"`
	LANE_DIR             int    `json:"lineflow"`
	NMIXSEC_ST_POS       int    `json:"ssecS"`
	NMIXSEC_END_POS      int    `json:"ssecE"`
	MIXSEC_ST_POS        int    `json:"msecS"`
	MIXSEC_END_POS       int    `json:"msecE"`
	BUS_LANE_YN          int    `json:"buslaneyn"`
	CTRL_YN              int    `json:"ctrlyn"`
	TRFFICVOL_ALLOT_RATE int    `json:"trffAllotRate"`
}

type QRY_DBImgDetInfo struct {
	LOC_ID               string `json:"id"`
	CAM_ID               string `json:"cam"`
	CHAN                 string `json:"chan"`
	LANE_NUM             string `json:"lane"`
	ACCROAD_DIR          string `json:"dir"`
	COMM_TYPE            string `json:"net"`
	RING_NUM             string `json:"ring"`
	PHASE_NUM            string `json:"phase"`
	LANE_DIR             string `json:"lineflow"`
	NMIXSEC_ST_POS       string `json:"ssecS"`
	NMIXSEC_END_POS      string `json:"ssecE"`
	MIXSEC_ST_POS        string `json:"msecS"`
	MIXSEC_END_POS       string `json:"msecE"`
	BUS_LANE_YN          string `json:"buslaneyn"`
	CTRL_YN              string `json:"ctrlyn"`
	TRFFICVOL_ALLOT_RATE string `json:"trffAllotRate"`
	REG_DT               string
	MOD_DT               string
}

// ---------------------------------------------------------------------------
// DBImgCamInfo
// ---------------------------------------------------------------------------
type DBImgCamInfo struct {
	LOC_ID      int    `json:"id"`
	CAM_ID      string `json:"cam"`
	ACCROAD_DIR int    `json:"dir"`
	COMM_TYPE   int    `json:"net"`
}

type QRY_DBImgCamInfo struct {
	LOC_ID      string `json:"id"`
	CAM_ID      string `json:"cam"`
	ACCROAD_DIR string `json:"dir"`
	COMM_TYPE   string `json:"net"`
	MOD_DT      string
}

type DBLocAddSigInfo struct {
	LOC_ID  int     `json:"-"`
	BANK_ID int     `json:"-"`
	STEP    int     `json:"-"`
	Siginfo [16]int `json:"-"` // 0:북(u), 15 : 북서(좌)
}

type QRY_DBLocAddSigInfo struct {
	LOC_ID  string `json:"-"`
	BANK_ID string `json:"-"`
	STEP    string `json:"-"`
	Siginfo string `json:"-"` // 0:북(u), 15 : 북서(좌)
}

// ---------------------------------------------------------------------------
// DBLocOprReqRst
// ---------------------------------------------------------------------------
type DBLocOprReqRst struct {
	LOC_ID  int        `json:"locId"`
	LOC_NUM int        `json:"number"`
	REQ_DT  string     `json:"reqDate"`
	COLL_DT NullString `json:"rstDate"`

	REQ_CTRLMODE int       `json:"reqCtrlMode"`
	REQ_CYCLE    int       `json:"reqCycle"`
	REQ_OFFSET   int       `json:"reqOffset"`
	REQ_PHASE    [16]uint8 `json:"reqPhase"`

	RST_CTRLMODE  int       `json:"rstCtrlMode"`
	RST_CTRLSTATE int       `json:"rstCtrlStt"`
	RST_PHASE     [16]uint8 `json:"rstPhase"`
	RST_CYCLE     int       `json:"rstCycle"`
}

type QRY_DBLocOprReqRst struct {
	LOC_ID  string `json:"locId"`
	LOC_NUM string `json:"number"`
	REQ_DT  string `json:"reqDate"`
	COLL_DT string `json:"rstDate"`

	REQ_CTRLMODE string `json:"reqCtrlMode"`
	REQ_CYCLE    string `json:"reqCycle"`
	REQ_OFFSET   string `json:"reqOffset"`
	REQ_PHASE    string `json:"reqPhase"`

	RST_CTRLMODE  string `json:"rstCtrlMode"`
	RST_CTRLSTATE string `json:"rstCtrlStt"`
	RST_PHASE     string `json:"rstPhase"`
	RST_CYCLE     string `json:"rstCycle"`
}

type DBLocalState struct {
	LOC_ID  int        `json:"id"`
	LOC_NM  NullString `json:"-"`
	LOC_NUM int        `json:"-"`
	COLL_DT string     `json:"dateTm"`
	State   [25]uint8  `json:"state"`
}

type QRY_DBLocalState struct {
	LOC_ID  string `json:"id"`
	LOC_NM  string `json:"-"`
	LOC_NUM string `json:"-"`
	COLL_DT string `json:"dateTm"`
	State   string `json:"state"`
}

// ---------------------------------------------------------------------------
// DBGRP_DAYPLAN
// ---------------------------------------------------------------------------
type DBGRP_DAYPLAN struct {
	GRP_ID             int `json:"gId"`
	GRP_NUM            int `json:"gNum"`
	GRP_PLANNUM        int `json:"planNum"`
	HOUR               int `json:"hour"`
	MIN                int `json:"min"`
	CYCLE_LV           int `json:"lv"`
	LOC_OFFSETPLAN_IDX int `json:"offIdx"`
	LOC_PHASEPLAN_IDX  int `json:"phaseIdx"`
	SIGOPR_MODE        int `json:"mode"`
}

type QRY_DBGRP_DAYPLAN struct {
	GRP_ID             string `json:"gId"`
	GRP_NUM            string `json:"gNum"`
	GRP_PLANNUM        string `json:"planNum"`
	HOUR               string `json:"hour"`
	MIN                string `json:"min"`
	CYCLE_LV           string `json:"lv"`
	LOC_OFFSETPLAN_IDX string `json:"offIdx"`
	LOC_PHASEPLAN_IDX  string `json:"phaseIdx"`
	SIGOPR_MODE        string `json:"mode"`
}

// ---------------------------------------------------------------------------
// REPORTGRP_DAYPLAN
// ---------------------------------------------------------------------------
type REPORTGRP_DAYPLAN struct {
	GRP_PLANNUM       int    `json:"-"`
	TM                string `json:"tm"`
	CYCLE_LV          int    `json:"cyclv"`
	LOC_PHASEPLAN_IDX int    `json:"phaseIdx"`
}

type QRY_REPORTGRP_DAYPLAN struct {
	GRP_PLANNUM       string `json:"-"`
	TM                string `json:"tm"`
	CYCLE_LV          string `json:"cyclv"`
	LOC_PHASEPLAN_IDX string `json:"phaseIdx"`
}

// ---------------------------------------------------------------------------
// DBGRP_WEEKPLAN
// ---------------------------------------------------------------------------
type DBGRP_WEEKPLAN struct {
	GRP_ID      int `json:"gId"`
	LOC_ID      int `json:"-"`
	GRP_NUM     int `json:"gNum"`
	DAYWEEK     int `json:"dayWeek"`
	GRP_PLANNUM int `json:"planNum"`
}

type QRY_DBGRP_WEEKPLAN struct {
	GRP_ID      string `json:"gId"`
	LOC_ID      string `json:"-"`
	GRP_NUM     string `json:"gNum"`
	DAYWEEK     string `json:"dayWeek"`
	GRP_PLANNUM string `json:"planNum"`
}

// ---------------------------------------------------------------------------
// DBGRP_HOLIPLAN
// ---------------------------------------------------------------------------
type DBGRP_HOLIPLAN struct {
	GRP_ID      int `json:"gId"`
	LOC_ID      int `json:"-"`
	GRP_NUM     int `json:"gNum"`
	MONTH       int `json:"month"`
	DAY         int `json:"day"`
	GRP_PLANNUM int `json:"planNum"`
}

type QRY_DBGRP_HOLIPLAN struct {
	GRP_ID      string `json:"gId"`
	LOC_ID      string `json:"-"`
	GRP_NUM     string `json:"gNum"`
	MONTH       string `json:"month"`
	DAY         string `json:"day"`
	GRP_PLANNUM string `json:"planNum"`
}

// ---------------------------------------------------------------------------
// DBGRP_CYCLEPLAN
// ---------------------------------------------------------------------------
type DBGRP_CYCLEPLAN struct {
	GRP_ID      int `json:"gId"`
	GRP_NUM     int `json:"gNum"`
	GRP_PLANNUM int `json:"planNum"`
	CYCLE_LV    int `json:"cyclv"`
	CYCLELEN    int `json:"cycLen"`
	OIRATE      int `json:"oirate"`
	IORATE      int `json:"iorate"`
	USE_YN      int `json:"useYn"`
	//USE_YN int `json:"USE_YN"`
}

type QRY_DBGRP_CYCLEPLAN struct {
	GRP_ID      string `json:"gId"`
	GRP_NUM     string `json:"gNum"`
	GRP_PLANNUM string `json:"planNum"`
	CYCLE_LV    string `json:"cyclv"`
	CYCLELEN    string `json:"cycLen"`
	OIRATE      string `json:"oirate"`
	IORATE      string `json:"iorate"`
	USE_YN      string `json:"useYn"`
}

// ---------------------------------------------------------------------------
// 센터 교차로 데이플랜정보
// ---------------------------------------------------------------------------
type DBGrpLocDayPlan struct {
	GRP_ID             int `json:"grpId"`
	GRP_PLANNUM        int `json:"planNum"`
	HOUR               int `json:"hour"`
	MIN                int `json:"min"`
	LOC_ID             int `json:"locId"`
	LOC_PHASEPLAN_IDX  int `json:"phaseIdx"`
	LOC_OFFSETPLAN_IDX int `json:"offsetIdx"`
	RINGA_1PHASE       int `json:"ringA1"`
	RINGA_2PHASE       int `json:"ringA2"`
	RINGA_3PHASE       int `json:"ringA3"`
	RINGA_4PHASE       int `json:"ringA4"`
	RINGA_5PHASE       int `json:"ringA5"`
	RINGA_6PHASE       int `json:"ringA6"`
	RINGA_7PHASE       int `json:"ringA7"`
	RINGA_8PHASE       int `json:"ringA8"`

	RINGB_1PHASE int `json:"ringB1"`
	RINGB_2PHASE int `json:"ringB2"`
	RINGB_3PHASE int `json:"ringB3"`
	RINGB_4PHASE int `json:"ringB4"`
	RINGB_5PHASE int `json:"ringB5"`
	RINGB_6PHASE int `json:"ringB6"`
	RINGB_7PHASE int `json:"ringB7"`
	RINGB_8PHASE int `json:"ringB8"`

	OFFSET int `json:"offset"`
}

type QRY_DBGrpLocDayPlan struct {
	GRP_ID             string `json:"grpId"`
	GRP_PLANNUM        string `json:"planNum"`
	HOUR               string `json:"hour"`
	MIN                string `json:"min"`
	LOC_ID             string `json:"locId"`
	LOC_PHASEPLAN_IDX  string `json:"phaseIdx"`
	LOC_OFFSETPLAN_IDX string `json:"offsetIdx"`
	RINGA_1PHASE       string `json:"ringA1"`
	RINGA_2PHASE       string `json:"ringA2"`
	RINGA_3PHASE       string `json:"ringA3"`
	RINGA_4PHASE       string `json:"ringA4"`
	RINGA_5PHASE       string `json:"ringA5"`
	RINGA_6PHASE       string `json:"ringA6"`
	RINGA_7PHASE       string `json:"ringA7"`
	RINGA_8PHASE       string `json:"ringA8"`

	RINGB_1PHASE string `json:"ringB1"`
	RINGB_2PHASE string `json:"ringB2"`
	RINGB_3PHASE string `json:"ringB3"`
	RINGB_4PHASE string `json:"ringB4"`
	RINGB_5PHASE string `json:"ringB5"`
	RINGB_6PHASE string `json:"ringB6"`
	RINGB_7PHASE string `json:"ringB7"`
	RINGB_8PHASE string `json:"ringB8"`

	OFFSET string `json:"offset"`
}

// ---------------------------------------------------------------------------
// REPORTLocTimePlan
// ---------------------------------------------------------------------------
type REPORTLocTimePlan struct {
	LOC_ID            int `json:"-"`
	LOC_PHASEPLAN_IDX int `json:"-"`
	CYCLE_LV          int `json:"-"`
	CYCLELEN          int `json:"-"`
	RINGA_1PHASE      int `json:"-"`
	RINGA_2PHASE      int `json:"-"`
	RINGA_3PHASE      int `json:"-"`
	RINGA_4PHASE      int `json:"-"`
	RINGA_5PHASE      int `json:"-"`
	RINGA_6PHASE      int `json:"-"`
	RINGA_7PHASE      int `json:"-"`
	RINGA_8PHASE      int `json:"-"`

	RINGB_1PHASE int `json:"-"`
	RINGB_2PHASE int `json:"-"`
	RINGB_3PHASE int `json:"-"`
	RINGB_4PHASE int `json:"-"`
	RINGB_5PHASE int `json:"-"`
	RINGB_6PHASE int `json:"-"`
	RINGB_7PHASE int `json:"-"`
	RINGB_8PHASE int `json:"-"`

	Aphases [8]int `json:"aPhases"`
	Bphases [8]int `json:"bPhases"`

	OFFSET int `json:"offset"`
}

type QRY_REPORTLocTimePlan struct {
	LOC_ID            string `json:"-"`
	LOC_PHASEPLAN_IDX string `json:"-"`
	CYCLE_LV          string `json:"-"`
	CYCLELEN          string `json:"-"`
	RINGA_1PHASE      string `json:"-"`
	RINGA_2PHASE      string `json:"-"`
	RINGA_3PHASE      string `json:"-"`
	RINGA_4PHASE      string `json:"-"`
	RINGA_5PHASE      string `json:"-"`
	RINGA_6PHASE      string `json:"-"`
	RINGA_7PHASE      string `json:"-"`
	RINGA_8PHASE      string `json:"-"`

	RINGB_1PHASE string `json:"-"`
	RINGB_2PHASE string `json:"-"`
	RINGB_3PHASE string `json:"-"`
	RINGB_4PHASE string `json:"-"`
	RINGB_5PHASE string `json:"-"`
	RINGB_6PHASE string `json:"-"`
	RINGB_7PHASE string `json:"-"`
	RINGB_8PHASE string `json:"-"`

	Aphases string `json:"aPhases"`
	Bphases string `json:"bPhases"`

	OFFSET string `json:"offset"`
}

// ---------------------------------------------------------------------------
// REPORTCycleInfo
// ---------------------------------------------------------------------------
type REPORTCycleInfo struct {
	LOC_ID   int `json:"-"`
	CYCLE_LV int `json:"-"`
	CYCLELEN int `json:"cycle"`
}

// ---------------------------------------------------------------------------
// DBGrpOprState
// ---------------------------------------------------------------------------
type DBGrpOprState struct {
	GRP_ID                 int    `json:"id"`
	GRP_NUM                int    `json:"gNum"`
	CREATE_DT              string `json:"createDt"`
	GRP_CTRLMODE           int    `json:"ctrlMode"`
	GRP_CTRLSTATE          int    `json:"ctrlStt"`
	GRP_PLANNUM            int    `json:"planNum"`
	NOW_GRP_CYCLELEN       int    `json:"nowCycLen"`
	NOW_GRP_CYCLE_LV       int    `json:"nowCycLv"`
	NOW_LOC_OFFSETPLAN_IDX int    `json:"nowOffsetIdx"`
	NOW_LOC_PHASEPLAN_IDX  int    `json:"nowPhaseIdx"`

	TOD_GRP_CYCLELEN       int `json:"todCycLen"`
	TOD_GRP_CYCLE_LV       int `json:"todCycLv"`
	TOD_LOC_OFFSETPLAN_IDX int `json:"todOffsetIdx"`
	TOD_LOC_PHASEPLAN_IDX  int `json:"todPhaseIdx"`

	TRC_GRP_CYCLELEN       int `json:"trcCycLen"`
	TRC_GRP_CYCLE_LV       int `json:"trcCycLv"`
	TRC_LOC_OFFSETPLAN_IDX int `json:"trcOffsetIdx"`
	TRC_LOC_PHASEPLAN_IDX  int `json:"trcPhaseIdx"`

	MAN_GRP_CYCLELEN       int `json:"manCycLen"`
	MAN_GRP_CYCLE_LV       int `json:"manCycLv"`
	MAN_LOC_OFFSETPLAN_IDX int `json:"manOffsetIdx"`
	MAN_LOC_PHASEPLAN_IDX  int `json:"manPhaseIdx"`
}

type QRY_DBGrpOprState struct {
	GRP_ID                 string `json:"id"`
	GRP_NUM                string `json:"gNum"`
	CREATE_DT              string `json:"createDt"`
	GRP_CTRLMODE           string `json:"ctrlMode"`
	GRP_CTRLSTATE          string `json:"ctrlStt"`
	GRP_PLANNUM            string `json:"planNum"`
	NOW_GRP_CYCLELEN       string `json:"nowCycLen"`
	NOW_GRP_CYCLE_LV       string `json:"nowCycLv"`
	NOW_LOC_OFFSETPLAN_IDX string `json:"nowOffsetIdx"`
	NOW_LOC_PHASEPLAN_IDX  string `json:"nowPhaseIdx"`

	TOD_GRP_CYCLELEN       string `json:"todCycLen"`
	TOD_GRP_CYCLE_LV       string `json:"todCycLv"`
	TOD_LOC_OFFSETPLAN_IDX string `json:"todOffsetIdx"`
	TOD_LOC_PHASEPLAN_IDX  string `json:"todPhaseIdx"`

	TRC_GRP_CYCLELEN       string `json:"trcCycLen"`
	TRC_GRP_CYCLE_LV       string `json:"trcCycLv"`
	TRC_LOC_OFFSETPLAN_IDX string `json:"trcOffsetIdx"`
	TRC_LOC_PHASEPLAN_IDX  string `json:"trcPhaseIdx"`

	MAN_GRP_CYCLELEN       string `json:"manCycLen"`
	MAN_GRP_CYCLE_LV       string `json:"manCycLv"`
	MAN_LOC_OFFSETPLAN_IDX string `json:"manOffsetIdx"`
	MAN_LOC_PHASEPLAN_IDX  string `json:"manPhaseIdx"`
}

// ---------------------------------------------------------------------------
// NODE
// ---------------------------------------------------------------------------
type DBNODE struct {
	NODE_ID string     `json:"id"`
	NODE_NM NullString `json:"name"`
	//LOC_ID           int     `json:"locId"`
	NODE_TYPE        int        `json:"type"`
	NODE_LAT         float64    `json:"lat"`
	NODE_LON         float64    `json:"lon"`
	TURN_P           NullString `json:"turnp"`
	REMARK           NullString `json:"remark"`
	GEOM             NullString `json:"-"`
	NODE_CREATE_TYPE int        `json:"createType"`
}

type QRY_DBNODE struct {
	NODE_ID string `json:"id"`
	NODE_NM string `json:"name"`
	//LOC_ID           int     `json:"locId"`
	NODE_TYPE        string `json:"type"`
	NODE_LAT         string `json:"lat"`
	NODE_LON         string `json:"lon"`
	TURN_P           string `json:"turnp"`
	REMARK           string `json:"remark"`
	GEOM             string `json:"-"`
	NODE_CREATE_TYPE string `json:"createType"`
}

// ---------------------------------------------------------------------------
// DBUnitEvent
// ---------------------------------------------------------------------------
type DBUnitEvent struct {
	UNIT_TYPE int        `json:"uType"`
	UNIT_NM   NullString `json:"uName"`
}

type QRY_DBUnitEvent struct {
	UNIT_TYPE string `json:"uType"`
	UNIT_NM   string `json:"uName"`
}

// ---------------------------------------------------------------------------
// DBEventCode
// ---------------------------------------------------------------------------
type DBEventCode struct {
	UNIT_TYPE  int        `json:"uType"`
	UNIT_NM    NullString `json:"uName"`
	EVENT_CODE int        `json:"code"`
	EVENT_DESC NullString `json:"codeDesc"`
	CODE_TYPE  int        `json:"type"`
}

type QRY_DBEventCode struct {
	UNIT_TYPE  string `json:"uType"`
	UNIT_NM    string `json:"uName"`
	EVENT_CODE string `json:"code"`
	EVENT_DESC string `json:"codeDesc"`
	CODE_TYPE  string `json:"type"`
}

// ---------------------------------------------------------------------------
// DBEventLowCode
// ---------------------------------------------------------------------------
type DBEventLowCode struct {
	UNIT_TYPE           int        `json:"uType"`
	UNIT_NM             NullString `json:"uName"`
	EVENT_CODE          int        `json:"code"`
	EVENT_DESC          NullString `json:"codeDesc"`
	EVENT_LOW_CODE      int        `json:"lowCode"`
	EVENT_LOW_CODE_DESC NullString `json:"lowCodeDesc"`
	LV                  int        `json:"lv"`
	CODE_TYPE           int        `json:"type"`
}

type QRY_DBEventLowCode struct {
	UNIT_TYPE           string `json:"uType"`
	UNIT_NM             string `json:"uName"`
	EVENT_CODE          string `json:"code"`
	EVENT_DESC          string `json:"codeDesc"`
	EVENT_LOW_CODE      string `json:"lowCode"`
	EVENT_LOW_CODE_DESC string `json:"lowCodeDesc"`
	LV                  string `json:"lv"`
	CODE_TYPE           string `json:"type"`
}

// ---------------------------------------------------------------------------
// DBEventCodeInfo (key value struct), for event code management
// ---------------------------------------------------------------------------
type DBEventCodeInfo struct {
	EVKEY               int        `json:"evKey"`
	UNIT_TYPE           int        `json:"uType"`
	UNIT_NM             NullString `json:"uName"`
	EVENT_CODE          int        `json:"code"`
	EVENT_DESC          NullString `json:"codeDesc"`
	EVENT_LOW_CODE      int        `json:"lowCode"`
	EVENT_LOW_CODE_DESC NullString `json:"lowCodeDesc"`
	LV                  int        `json:"lv"`
	CODE_TYPE           int        `json:"type"`
}

type QRY_DBEventCodeInfo struct {
	EVKEY               string `json:"evKey"`
	UNIT_TYPE           string `json:"uType"`
	UNIT_NM             string `json:"uName"`
	EVENT_CODE          string `json:"code"`
	EVENT_DESC          string `json:"codeDesc"`
	EVENT_LOW_CODE      string `json:"lowCode"`
	EVENT_LOW_CODE_DESC string `json:"lowCodeDesc"`
	LV                  string `json:"lv"`
	CODE_TYPE           string `json:"type"`
}

// ---------------------------------------------------------------------------
// DBEventStateCode
// ---------------------------------------------------------------------------
type DBEventStateCode struct {
	CODE      int
	CODE_VAL  int
	CODE_NM   NullString
	STATE_VAL NullString
}

type QRY_DBEventStateCode struct {
	CODE      string
	CODE_VAL  string
	CODE_NM   string
	STATE_VAL string
}

// ---------------------------------------------------------------------------
// DBEventFormat
// ---------------------------------------------------------------------------
type DBEventFormat struct {
	EVENT_CODE   int
	EVENT_FORMAT NullString
}

type QRY_DBEventFormat struct {
	EVENT_CODE   string
	EVENT_FORMAT string
}

// ---------------------------------------------------------------------------
// DBEventLog
// ---------------------------------------------------------------------------
type DBEventLog struct {
	UNIT_TYPE      int    `json:"uType"`
	EVENT_CODE     int    `json:"code"`
	EVENT_LOW_CODE int    `json:"lowCode"`
	UNIT_ID        int    `json:"uId"`
	CREATE_DATE    string `json:"createDt"`
	CREATE_TIME    string `json:"-"`
	CTOR           int    `json:"ctor"`
	STR_CTOR       string `json:"strCtor"`
	TARGET         int    `json:"target"`
	STR_TARGET     string `json:"strTarget"`
	PARAM_1        int    `json:"parma1"`
	PARAM_2        int    `json:"parma2"`
	PARAM_3        int    `json:"parma3"`
	PARAM_4        int    `json:"parma4"`
	PARAM_5        int    `json:"parma5"`
	Format         string `json:"format"`
}

type QRY_DBEventLog struct {
	UNIT_TYPE      string `json:"uType"`
	EVENT_CODE     string `json:"code"`
	EVENT_LOW_CODE string `json:"lowCode"`
	UNIT_ID        string `json:"uId"`
	CREATE_DATE    string `json:"createDt"`
	CREATE_TIME    string `json:"-"`
	SEQ            string `json:"-"`
	CTOR           string `json:"ctor"`
	STR_CTOR       string `json:"strCtor"`
	TARGET         string `json:"target"`
	STR_TARGET     string `json:"strTarget"`
	PARAM_1        string `json:"parma1"`
	PARAM_2        string `json:"parma2"`
	PARAM_3        string `json:"parma3"`
	PARAM_4        string `json:"parma4"`
	PARAM_5        string `json:"parma5"`
	Format         string `json:"format"`
}

type DBGrpLocalPhasePlan struct {
	GRP_ID            int `json:"gId"`
	GRP_NUM           int `json:"gNum"`
	LOC_ID            int `json:"locId"`
	LOC_NUM           int `json:"number"`
	LOC_PLANNUM       int `json:"planNum"`
	LOC_PHASEPLAN_IDX int `json:"phaseIdx"`

	RINGA_1PHASE int `json:"ringA1"`
	RINGA_2PHASE int `json:"ringA2"`
	RINGA_3PHASE int `json:"ringA3"`
	RINGA_4PHASE int `json:"ringA4"`
	RINGA_5PHASE int `json:"ringA5"`
	RINGA_6PHASE int `json:"ringA6"`
	RINGA_7PHASE int `json:"ringA7"`
	RINGA_8PHASE int `json:"ringA8"`

	RINGB_1PHASE int `json:"ringB1"`
	RINGB_2PHASE int `json:"ringB2"`
	RINGB_3PHASE int `json:"ringB3"`
	RINGB_4PHASE int `json:"ringB4"`
	RINGB_5PHASE int `json:"ringB5"`
	RINGB_6PHASE int `json:"ringB6"`
	RINGB_7PHASE int `json:"ringB7"`
	RINGB_8PHASE int `json:"ringB8"`
	OFFSET       int `json:"offset"`
}

type QRY_DBGrpLocalPhasePlan struct {
	GRP_ID            string `json:"gId"`
	GRP_NUM           string `json:"gNum"`
	LOC_ID            string `json:"locId"`
	LOC_NUM           string `json:"number"`
	LOC_PLANNUM       string `json:"planNum"`
	LOC_PHASEPLAN_IDX string `json:"phaseIdx"`

	RINGA_1PHASE string `json:"ringA1"`
	RINGA_2PHASE string `json:"ringA2"`
	RINGA_3PHASE string `json:"ringA3"`
	RINGA_4PHASE string `json:"ringA4"`
	RINGA_5PHASE string `json:"ringA5"`
	RINGA_6PHASE string `json:"ringA6"`
	RINGA_7PHASE string `json:"ringA7"`
	RINGA_8PHASE string `json:"ringA8"`

	RINGB_1PHASE string `json:"ringB1"`
	RINGB_2PHASE string `json:"ringB2"`
	RINGB_3PHASE string `json:"ringB3"`
	RINGB_4PHASE string `json:"ringB4"`
	RINGB_5PHASE string `json:"ringB5"`
	RINGB_6PHASE string `json:"ringB6"`
	RINGB_7PHASE string `json:"ringB7"`
	RINGB_8PHASE string `json:"ringB8"`
	OFFSET       string `json:"offset"`
}

type OprtStateColor struct {
	Defalut   string `json:"default"`
	Local     string `json:"msttLocal"`
	Center    string `json:"msttCenterTod"`
	CommErr   string `json:"msttCommErr"`
	Trans     string `json:"msttTrans"`
	Flash     string `json:"msttFlash"`
	LightOff  string `json:"msttLightOff"`
	Manual    string `json:"msttManual"`
	ManFlash  string `json:"msttManFlash"`
	ManOff    string `json:"msttManOff"`
	Keepphase string `json:"msttKeepPhase"`
	Noreg     string `json:"msttNoreg"`
	Conflict  string `json:"msttConflic"`
	CenterTrc string `json:"msttCenterTrc"`
}

// -----------------------------------------------------
// ####################################################
// -----------------------------------------------------
type DBOPRT struct {
	OPRT_ID  string     `json:"id"`
	ID_IDX   int        `json:"ididx"`
	OPRT_PW  NullString `json:"pw"`
	OPRT_NM  NullString `json:"nm"`
	OPRT_TEL NullString `json:"tel"`
	GRP_ID   int        `json:"grpId"`
	GRP_NM   NullString `json:"-"`
	LoginStt bool       `json:"-"`
	RefToken string     `json:"-"`
	AccToken string     `json:"-"`
}

type QRY_DBOPRT struct {
	OPRT_ID  string `json:"id"`
	ID_IDX   string `json:"ididx"`
	OPRT_PW  string `json:"pw"`
	OPRT_NM  string `json:"nm"`
	OPRT_TEL string `json:"tel"`
	GRP_ID   string `json:"grpId"`
	GRP_NM   string `json:"-"`
	LoginStt string `json:"-"`
	RefToken string `json:"-"`
	AccToken string `json:"-"`
}

type DBOPRT2 struct {
	ID_IDX   int        `json:"ididx"`
	OPRT_ID  string     `json:"id"`
	OPRT_PW  NullString `json:"pw"`
	OPRT_NM  NullString `json:"nm"`
	OPRT_TEL NullString `json:"tel"`
	//SIGN_KEY string         `json:"key"`
	GRP_NM  NullString     `json:"grpNm"`
	STT_COL OprtStateColor `json:"-"`
}

type QRY_DBOPRT2 struct {
	ID_IDX   int    `json:"ididx"`
	OPRT_ID  string `json:"id"`
	OPRT_PW  string `json:"pw"`
	OPRT_NM  string `json:"nm"`
	OPRT_TEL string `json:"tel"`
	//SIGN_KEY string         `json:"key"`
	GRP_NM  string         `json:"grpNm"`
	STT_COL OprtStateColor `json:"-"`
}

type DBOPRTGRP_MENU struct {
	GRP_ID      int        `json:"-"`
	GRP_NM      NullString `json:"-"`
	MENU_GRP_ID int        `json:"-"`
	MENU_GRP_NM NullString `json:"menuNm"`
}

type QRY_DBOPRTGRP_MENU struct {
	GRP_ID      string `json:"-"`
	GRP_NM      string `json:"-"`
	MENU_GRP_ID string `json:"-"`
	MENU_GRP_NM string `json:"menuNm"`
}

type OPRTGRPInfo struct {
	GRP_ID      int    `json:"-"`
	GRP_NM      string `json:"grpNm"`
	MENU_GRP_ID []int  `json:"menu"`
}

type DBOPRT_GRP struct {
	GRP_ID int        `json:"grpId"`
	GRP_NM NullString `json:"grpNm"`
}

type DBMENU struct {
	MENU_GRP_ID int        `json:"menuId"`
	MENU_GRP_NM NullString `json:"menuNm"`
}

type QRY_DBMENU struct {
	MENU_GRP_ID string `json:"menuId"`
	MENU_GRP_NM string `json:"menuNm"`
}

type UserMenuList struct {
	ID_IDX      int        `json:"ididx"`
	OPRT_ID     string     `json:"id"`
	GRP_ID      int        `json:"grpId"`
	MENU_GRP_ID int        `json:"menuGrpId"`
	MENU_GRP_NM NullString `json:"menuGrpNm"`
	AUTH        int        `json:"auth"`
	JSON_NM     NullString `json:"-"`
}

type QRY_UserMenuList struct {
	ID_IDX      string `json:"ididx"`
	OPRT_ID     string `json:"id"`
	GRP_ID      string `json:"grpId"`
	MENU_GRP_ID string `json:"menuGrpId"`
	MENU_GRP_NM string `json:"menuGrpNm"`
	AUTH        string `json:"auth"`
	JSON_NM     string `json:"-"`
}

type UserMapPresetList struct {
	ID_IDX int     `json:"-"`
	SEQ    int     `json:"-"`
	LAT    float64 `json:"lat"`
	LON    float64 `json:"lon"`
	ZOOMLV float64 `json:"zoom"`
}

type QRY_UserMapPresetList struct {
	ID_IDX string `json:"idx"`
	SEQ    string `json:"presetno"`
	LAT    string `json:"lat"`
	LON    string `json:"lon"`
	ZOOMLV string `json:"zoom"`
}

type DBROUTE struct {
	ROUTE_ID int    `json:"routeId"`
	ROUTE_NM string `json:"routeNm"`
	COLOR_A  int    `json:"colorA"`
	COLOR_R  int    `json:"colorR"`
	COLOR_G  int    `json:"colorG"`
	COLOR_B  int    `json:"colorB"`
	REMARK   string `json:"remark"`
}

type QRY_DBROUTE struct {
	ROUTE_ID string `json:"routeId"`
	ROUTE_NM string `json:"routeNm"`
	COLOR_A  string `json:"colorA"`
	COLOR_R  string `json:"colorR"`
	COLOR_G  string `json:"colorG"`
	COLOR_B  string `json:"colorB"`
	REMARK   string `json:"remark"`
}

type DBROUTE_PATH_INFO struct {
	ROUTE_ID int     `json:"routeId"`
	PATH_NUM int     `json:"-"`
	SEQ      int     `json:"-"`
	POS_LAT  float64 `json:"lat"`
	POS_LON  float64 `json:"lon"`
}

type QRY_DBROUTE_PATH_INFO struct {
	ROUTE_ID string `json:"routeId"`
	PATH_NUM string `json:"-"`
	SEQ      string `json:"-"`
	POS_LAT  string `json:"lat"`
	POS_LON  string `json:"lon"`
}

type DBGRP_PATH_INFO struct {
	GRP_ID  int        `json:"grpId"`
	GRP_NM  string     `json:"grpNm"`
	COLOR_A int        `json:"colorA"`
	COLOR_R int        `json:"colorR"`
	COLOR_G int        `json:"colorG"`
	COLOR_B int        `json:"colorB"`
	REMARK  NullString `json:"remark"`
}

type QRY_DBGRP_PATH_INFO struct {
	GRP_ID  string `json:"grpId"`
	GRP_NM  string `json:"grpNm"`
	COLOR_A string `json:"colorA"`
	COLOR_R string `json:"colorR"`
	COLOR_G string `json:"colorG"`
	COLOR_B string `json:"colorB"`
	REMARK  string `json:"remark"`
}

type DBGRP_PATH struct {
	GRP_ID   int     `json:"grpId"`
	PATH_NUM int     `json:"-"`
	SEQ      int     `json:"-"`
	POS_LAT  float64 `json:"lat"`
	POS_LON  float64 `json:"lon"`
}

type QRY_DBGRP_PATH struct {
	GRP_ID   string `json:"grpId"`
	PATH_NUM string `json:"-"`
	SEQ      string `json:"-"`
	POS_LAT  string `json:"lat"`
	POS_LON  string `json:"lon"`
}

type DB_SYS_ENV_CONF struct {
	SVR_COMM_PORT      int `json:"svr_comm_port"`
	TIME_SYNC_CYCLE    int `json:"time_sync_cycle"`
	STATE_REQ_CYCLE    int `json:"state_req_cycle"`
	SIG_SVR_OPR_MODE   int `json:"sig_svr_opr_mode"`
	FUNC_CMD_PROC_MODE int `json:"func_cmd_proc_mode"`

	EXT_SVR1_NM   NullString `json:"-"`
	EXT_SVR1_IP   NullString `json:"-"`
	EXT_SVR1_PORT int        `json:"-"`
	EXT_SVR1_MODE int        `json:"-"`

	EXT_SVR2_NM   NullString `json:"-"`
	EXT_SVR2_IP   NullString `json:"-"`
	EXT_SVR2_PORT int        `json:"-"`
	EXT_SVR2_MODE int        `json:"-"`

	EXT_SVR3_NM   NullString `json:"-"`
	EXT_SVR3_IP   NullString `json:"-"`
	EXT_SVR3_PORT int        `json:"-"`
	EXT_SVR3_MODE int        `json:"-"`

	EXT_SVR4_NM   NullString `json:"-"`
	EXT_SVR4_IP   NullString `json:"-"`
	EXT_SVR4_PORT int        `json:"-"`
	EXT_SVR4_MODE int        `json:"-"`
}

type QRY_DB_SYS_ENV_CONF struct {
	SVR_COMM_PORT      string `json:"svr_comm_port"`
	TIME_SYNC_CYCLE    string `json:"time_sync_cycle"`
	STATE_REQ_CYCLE    string `json:"state_req_cycle"`
	SIG_SVR_OPR_MODE   string `json:"sig_svr_opr_mode"`
	FUNC_CMD_PROC_MODE string `json:"func_cmd_proc_mode"`

	EXT_SVR1_NM   string `json:"-"`
	EXT_SVR1_IP   string `json:"-"`
	EXT_SVR1_PORT string `json:"-"`
	EXT_SVR1_MODE string `json:"-"`

	EXT_SVR2_NM   string `json:"-"`
	EXT_SVR2_IP   string `json:"-"`
	EXT_SVR2_PORT string `json:"-"`
	EXT_SVR2_MODE string `json:"-"`

	EXT_SVR3_NM   string `json:"-"`
	EXT_SVR3_IP   string `json:"-"`
	EXT_SVR3_PORT string `json:"-"`
	EXT_SVR3_MODE string `json:"-"`

	EXT_SVR4_NM   string `json:"-"`
	EXT_SVR4_IP   string `json:"-"`
	EXT_SVR4_PORT string `json:"-"`
	EXT_SVR4_MODE string `json:"-"`
}

type DB_VER struct {
	DBVER    NullString
	OEM      NullString
	VER_DESC NullString
	CHG_DT   NullString
}

type QRY_DB_VER struct {
	DBVER    string
	OEM      string
	VER_DESC string
	CHG_DT   string
}

type DB_FLOW_INFO struct {
	LOC_ID   int
	LOC_NUM  int
	FLOW_NUM int
	T_POS_X  float64
	T_POS_Y  float64

	C_POS_X float64
	C_POS_Y float64

	H_POS_X float64
	H_POS_Y float64
}

type QRY_DB_FLOW_INFO struct {
	LOC_ID   string
	LOC_NUM  string
	FLOW_NUM string
	T_POS_X  string
	T_POS_Y  string

	C_POS_X string
	C_POS_Y string

	H_POS_X string
	H_POS_Y string
}

type LocalAngleInfo struct {
	Id    int     `json:"id"`
	Angle float64 `json:"angle"`
	//Flow  [][3][2]float64 `json:"flowShape"`
	//Flow  [][3][2]float64 `json:"flowShape"`
	Flow map[int][3][2]float64 `json:"flowShape"`
}

type DB_SigongdoInfo struct {
	LOC_ID            int        `json:"-"`
	LOC_NM            NullString `json:"nm"`
	LOC_PLANNUM       int        `json:"plan"`
	LOC_PHASEPLAN_IDX int        `json:"entry"`
	MainPhase         int        `json:"mainPhase"`
	OFFSET            int        `json:"offset"`
	APhaseCnt         int        `json:"aPhaseCnt"`
	BPhaseCnt         int        `json:"bPhaseCnt"`
	Aphase            [8]int     `json:"aRing"`
	Bphase            [8]int     `json:"bRing"`
	AFlow             [8]int     `json:"aFlow"`
	BFlow             [8]int     `json:"bFlow"`
}

type DBSGD_AXIS struct {
	ID       int        `json:"-"`
	AXIS_NM  NullString `json:"name"`
	GRP_ID   int        `json:"gid"`
	PLAN_NUM int        `json:"-"`
	HOUR     int        `json:"-"`
	MIN      int        `json:"-"`

	//AXIS_TMS []DBSGD_AXIS_TM `json:"tms"`
	AXIS_TMS map[int]DBSGD_AXIS_TM `json:"tms"`
}

type DBSGD_AXIS_TM struct {
	ID           int        `json:"id"`
	SEQ          int        `json:"idx"`
	AXIS_NM      NullString `json:"-"`
	GRP_ID       int        `json:"-"`
	HOUR         int        `json:"hour"`
	MIN          int        `json:"min"`
	PLAN_NUM     int        `json:"plan"`
	PEAK_TM_TYPE int        `json:"peektm"`
}

type QRY_DBSGD_AXIS_TM struct {
	ID           string `json:"id"`
	SEQ          string `json:"idx"`
	AXIS_NM      string `json:"-"`
	GRP_ID       string `json:"-"`
	HOUR         string `json:"hour"`
	MIN          string `json:"min"`
	PLAN_NUM     string `json:"plan"`
	PEAK_TM_TYPE string `json:"peektm"`
}

type DBSGD_AXIS_LOC struct {
	ID       int        `json:"-"`
	IDX      int        `json:"-"`
	AXIS_NM  NullString `json:"-"`
	GRP_ID   int        `json:"-"`
	PLAN_NUM int        `json:"-"`
	HOUR     int        `json:"-"`
	MIN      int        `json:"-"`

	SEQ                int        `json:"-"`
	LOC_ID             int        `json:"locId"`
	LOC_NM             NullString `json:"name"`
	MAIN_PHASE         int        `json:"mainPhase"`
	DIST               int        `json:"distance"`
	ARING_OFFSET_SPEED int        `json:"aSpeed"`
	BRING_OFFSET_SPEED int        `json:"bSpeed"`
	ARING_OFFSET_FLOW  int        `json:"aflownum"`
	BRING_OFFSET_FLOW  int        `json:"bflownum"`
	UPDIR_RING         int        `json:"dir"`
	DirAring           int        `json:"dira"`
	DirBring           int        `json:"dirb"`

	OFFSET int `json:"offset"`

	APhaseCnt int    `json:"aPhaseCnt"`
	BPhaseCnt int    `json:"bPhaseCnt"`
	Aphase    [8]int `json:"-"`
	Bphase    [8]int `json:"-"`
	AFlow     [8]int `json:"-"`
	BFlow     [8]int `json:"-"`

	StrAphase string `json:"aRing"`
	StrBphase string `json:"bRing"`
	StrAFlow  string `json:"aFlow"`
	StrBFlow  string `json:"bFlow"`
}

type QRY_DBSGD_AXIS_LOC struct {
	ID       string `json:"-"`
	IDX      string `json:"-"`
	AXIS_NM  string `json:"-"`
	GRP_ID   string `json:"-"`
	PLAN_NUM string `json:"-"`
	HOUR     string `json:"-"`
	MIN      string `json:"-"`

	SEQ                string `json:"-"`
	LOC_ID             string `json:"locId"`
	LOC_NM             string `json:"name"`
	MAIN_PHASE         string `json:"mainPhase"`
	DIST               string `json:"distance"`
	ARING_OFFSET_SPEED string `json:"aSpeed"`
	BRING_OFFSET_SPEED string `json:"bSpeed"`
	ARING_OFFSET_FLOW  string `json:"aflownum"`
	BRING_OFFSET_FLOW  string `json:"bflownum"`
	UPDIR_RING         string `json:"dir"`
	DirAring           string `json:"dira"`
	DirBring           string `json:"dirb"`

	OFFSET string `json:"offset"`

	RINGA_1PHASE string
	RINGA_2PHASE string
	RINGA_3PHASE string
	RINGA_4PHASE string
	RINGA_5PHASE string
	RINGA_6PHASE string
	RINGA_7PHASE string
	RINGA_8PHASE string
	RINGB_1PHASE string
	RINGB_2PHASE string
	RINGB_3PHASE string
	RINGB_4PHASE string
	RINGB_5PHASE string
	RINGB_6PHASE string
	RINGB_7PHASE string
	RINGB_8PHASE string

	APhaseCnt string `json:"aPhaseCnt"`
	BPhaseCnt string `json:"bPhaseCnt"`
	Aphase    string `json:"-"`
	Bphase    string `json:"-"`
	AFlow     string `json:"-"`
	BFlow     string `json:"-"`

	StrAphase string `json:"aRing"`
	StrBphase string `json:"bRing"`
	StrAFlow  string `json:"aFlow"`
	StrBFlow  string `json:"bFlow"`
}

type DB_LOC_TRFFIC_INFO_LOG struct {
	LOC_ID   int        `json:"id"`
	LOC_NUM  int        `json:"num"`
	LOC_NM   NullString `json:"nm"`
	CAM_ID   string     `json:"cam"`
	LANE_NUM int        `json:"lane"`
	COLL_DT  string     `json:"dt"`

	GO_TRFFICVOL_S    int `json:"gSm"`
	GO_TRFFICVOL_BUS  int `json:"gBus"`
	GO_TRFFICVOL_B    int `json:"gBg"`
	GO_TRFFICVOL_BIKE int `json:"gBk"`
	LT_TRFFICVOL_S    int `json:"lSm"`
	LT_TRFFICVOL_BUS  int `json:"lBus"`
	LT_TRFFICVOL_B    int `json:"lBg"`
	LT_TRFFICVOL_BIKE int `json:"lBk"`
	RT_TRFFICVOL_S    int `json:"rSm"`
	RT_TRFFICVOL_BUS  int `json:"rBus"`
	RT_TRFFICVOL_B    int `json:"rBg"`
	RT_TRFFICVOL_BIKE int `json:"rBk"`
	U_TRFFICVOL_S     int `json:"uSm"`
	U_TRFFICVOL_BUS   int `json:"uBus"`
	U_TRFFICVOL_B     int `json:"uBg"`
	U_TRFFICVOL_BIKE  int `json:"uBk"`

	NMIXSEC_REMAIN_VH_S    int `json:"nsecRmSm"`
	NMIXSEC_REMAIN_VH_B    int `json:"nsecRmBg"`
	NMIXSEC_REMAIN_VH_BUS  int `json:"nsecRmBus"`
	NMIXSEC_REMAIN_VH_BIKE int `json:"nsecRmBk"`
	MIXSEC_REMAIN_VH_S     int `json:"msecRmSm"`
	MIXSEC_REMAIN_VH_B     int `json:"msecRmBg"`
	MIXSEC_REMAIN_VH_BUS   int `json:"msecRmBus"`
	MIXSEC_REMAIN_VH_BIKE  int `json:"msecRmBk"`

	NMIXSEC_STDBY_VH_S    int `json:"nsecStdbySm"`
	NMIXSEC_STDBY_VH_B    int `json:"nsecStdbyBg"`
	NMIXSEC_STDBY_VH_BUS  int `json:"nsecStdbyBus"`
	NMIXSEC_STDBY_VH_BIKE int `json:"nsecStdbyBk"`
	MIXSEC_STDBY_VH_S     int `json:"msecStdbySm"`
	MIXSEC_STDBY_VH_B     int `json:"msecStdbyBg"`
	MIXSEC_STDBY_VH_BUS   int `json:"msecStdbyBus"`
	MIXSEC_STDBY_VH_BIKE  int `json:"msecStdbyBk"`

	NMIXSEC_AVG_STDBY_VH_COUNT int     `json:"nsecAvgStdbyCnt"`
	MIXSEC_AVG_STDBY_VH_COUNT  int     `json:"msecAvgStdbyCnt"`
	NMIXSEC_MOM_SP_OCC_RATIO   float64 `json:"nsecMomOccRat"`
	MIXSEC_MOM_SP_OCC_RATIO    float64 `json:"msecMomOccRat"`
	NMIXSEC_SP_OCC_RATIO       float64 `json:"nsecOccRat"`
	MIXSEC_SP_OCC_RATIO        float64 `json:"msecOccRat"`
	AVG_SPEED                  float64 `json:"avgSpeed"`
}

type QRY_DB_LOC_TRFFIC_INFO_LOG struct {
	LOC_ID   string `json:"id"`
	LOC_NUM  string `json:"num"`
	LOC_NM   string `json:"nm"`
	CAM_ID   string `json:"cam"`
	LANE_NUM string `json:"lane"`
	COLL_DT  string `json:"dt"`

	GO_TRFFICVOL_S    string `json:"gSm"`
	GO_TRFFICVOL_BUS  string `json:"gBus"`
	GO_TRFFICVOL_B    string `json:"gBg"`
	GO_TRFFICVOL_BIKE string `json:"gBk"`
	LT_TRFFICVOL_S    string `json:"lSm"`
	LT_TRFFICVOL_BUS  string `json:"lBus"`
	LT_TRFFICVOL_B    string `json:"lBg"`
	LT_TRFFICVOL_BIKE string `json:"lBk"`
	RT_TRFFICVOL_S    string `json:"rSm"`
	RT_TRFFICVOL_BUS  string `json:"rBus"`
	RT_TRFFICVOL_B    string `json:"rBg"`
	RT_TRFFICVOL_BIKE string `json:"rBk"`
	U_TRFFICVOL_S     string `json:"uSm"`
	U_TRFFICVOL_BUS   string `json:"uBus"`
	U_TRFFICVOL_B     string `json:"uBg"`
	U_TRFFICVOL_BIKE  string `json:"uBk"`

	NMIXSEC_REMAIN_VH_S    string `json:"nsecRmSm"`
	NMIXSEC_REMAIN_VH_B    string `json:"nsecRmBg"`
	NMIXSEC_REMAIN_VH_BUS  string `json:"nsecRmBus"`
	NMIXSEC_REMAIN_VH_BIKE string `json:"nsecRmBk"`
	MIXSEC_REMAIN_VH_S     string `json:"msecRmSm"`
	MIXSEC_REMAIN_VH_B     string `json:"msecRmBg"`
	MIXSEC_REMAIN_VH_BUS   string `json:"msecRmBus"`
	MIXSEC_REMAIN_VH_BIKE  string `json:"msecRmBk"`

	NMIXSEC_STDBY_VH_S    string `json:"nsecStdbySm"`
	NMIXSEC_STDBY_VH_B    string `json:"nsecStdbyBg"`
	NMIXSEC_STDBY_VH_BUS  string `json:"nsecStdbyBus"`
	NMIXSEC_STDBY_VH_BIKE string `json:"nsecStdbyBk"`
	MIXSEC_STDBY_VH_S     string `json:"msecStdbySm"`
	MIXSEC_STDBY_VH_B     string `json:"msecStdbyBg"`
	MIXSEC_STDBY_VH_BUS   string `json:"msecStdbyBus"`
	MIXSEC_STDBY_VH_BIKE  string `json:"msecStdbyBk"`

	NMIXSEC_AVG_STDBY_VH_COUNT string `json:"nsecAvgStdbyCnt"`
	MIXSEC_AVG_STDBY_VH_COUNT  string `json:"msecAvgStdbyCnt"`
	NMIXSEC_MOM_SP_OCC_RATIO   string `json:"nsecMomOccRat"`
	MIXSEC_MOM_SP_OCC_RATIO    string `json:"msecMomOccRat"`
	NMIXSEC_SP_OCC_RATIO       string `json:"nsecOccRat"`
	MIXSEC_SP_OCC_RATIO        string `json:"msecOccRat"`
	AVG_SPEED                  string `json:"avgSpeed"`
}

type DB_LOC_IMGDETCAM_STT_LOG struct {
	LOC_ID      int        `json:"id"`
	LOC_NUM     int        `json:"num"`
	LOC_NM      NullString `json:"nm"`
	CAM_ID      string     `json:"cam"`
	COLL_DT     string     `json:"dt"`
	ACCROAD_DIR int        `json:"dir"`
	STATE       int        `json:"stt"`
}

type QRY_DB_LOC_IMGDETCAM_STT_LOG struct {
	LOC_ID      string `json:"id"`
	LOC_NUM     string `json:"num"`
	LOC_NM      string `json:"nm"`
	CAM_ID      string `json:"cam"`
	COLL_DT     string `json:"dt"`
	ACCROAD_DIR string `json:"dir"`
	STATE       string `json:"stt"`
}

type DB_TOD_SS_INFO struct {
	SS_ID     float64 `json:"id"`
	SA_NO     int     `json:"-"`
	PLAN_NO   int     `json:"plan"`
	CREATE_DT string  `json:"dt"`
	VALID_YN  int     `json:"-"`
}

type DB_SCM_BOD struct {
	SCM_ID   int        `json:"-"`
	IPADDR   NullString `json:"ip"`
	COMMTYPE int        `json:"commType"`
	NET_PORT int        `json:"port"`
}

type QRY_DB_SCM_BOD struct {
	SCM_ID   string `json:"-"`
	IPADDR   string `json:"ip"`
	COMMTYPE string `json:"commType"`
	NET_PORT string `json:"port"`
}

type SCMBOD_INFO struct {
	LOC_ID   int        `json:"id"`
	SCM_ID   int        `json:"scmid"`
	IPADDR   NullString `json:"ip"`
	COMMTYPE int        `json:"commType"`
	NET_PORT int        `json:"port"`
}

type DBLOC_FLOW_TRFFIC_INFO_LOG struct {
	LOC_ID               int        `json:"id"`
	LOC_NM               NullString `json:"nm"`
	CAM_ID               string     `json:"cam"`
	FLOW_NUM             int        `json:"flow"`
	COLL_DT              string     `json:"dt"`
	SATUR                float64    `json:"satur"`
	AVG_STDBYVH          float64    `json:"avgstdvh"`
	MAX_STDBYVH          float64    `json:"maxstdvh"`
	PASS_TRFFICVOL       float64    `json:"passtrff"`
	REMAIN_TRFFICVOL     float64    `json:"remaintrff"`
	DEMAND_TRFFICVOL     float64    `json:"demandtrff"`
	AVG_DEMAND_TRFFICVOL float64    `json:"avgdemandtrff"`
	LANE_FACTOR          float64    `json:"lanefactor"`
}

type QRY_DBLOC_FLOW_TRFFIC_INFO_LOG struct {
	LOC_ID               string `json:"id"`
	LOC_NM               string `json:"nm"`
	CAM_ID               string `json:"cam"`
	FLOW_NUM             string `json:"flow"`
	COLL_DT              string `json:"dt"`
	SATUR                string `json:"satur"`
	AVG_STDBYVH          string `json:"avgstdvh"`
	MAX_STDBYVH          string `json:"maxstdvh"`
	PASS_TRFFICVOL       string `json:"passtrff"`
	REMAIN_TRFFICVOL     string `json:"remaintrff"`
	DEMAND_TRFFICVOL     string `json:"demandtrff"`
	AVG_DEMAND_TRFFICVOL string `json:"avgdemandtrff"`
	LANE_FACTOR          string `json:"lanefactor"`
}

type DBLOC_FLOW_30MIN_LOS_INFO_LOG struct {
	LOC_ID              int     `json:"-"`
	FLOW_NUM            int     `json:"-"`
	CREATE_DT           string  `json:"-"`
	REF_TM_ID           int     `json:"-"`
	INFO_VALID_ID       int     `json:"-"`
	CAPA                float64 `json:"-"`
	TRFFICVOLRATE       float64 `json:"-"`
	SATUR               float64 `json:"-"`
	FIR_STDBYVH         float64 `json:"-"`
	HVY_VH_CORRVAL      float64 `json:"-"`
	HVY_VH_MIXRATE      float64 `json:"-"`
	SFR                 float64 `json:"-"`
	TRFFICVOL           float64 `json:"trffvol"`
	AVGCTRLDLY          float64 `json:"-"`
	UNIFDLY             float64 `json:"-"`
	INCRDLY             float64 `json:"-"`
	INITQDLY            float64 `json:"-"`
	AVG_STDBYVH         float64 `json:"-"`
	AVG_STDBYLEN        float64 `json:"waitlen"`
	DEMAND_TRFFICVOLLEN float64 `json:"demand"`
	LOS                 int     `json:"los"`
}

type QRY_DBLOC_FLOW_30MIN_LOS_INFO_LOG struct {
	LOC_ID              string `json:"-"`
	FLOW_NUM            string `json:"-"`
	CREATE_DT           string `json:"-"`
	REF_TM_ID           string `json:"-"`
	INFO_VALID_ID       string `json:"-"`
	CAPA                string `json:"-"`
	TRFFICVOLRATE       string `json:"-"`
	SATUR               string `json:"-"`
	FIR_STDBYVH         string `json:"-"`
	HVY_VH_CORRVAL      string `json:"-"`
	HVY_VH_MIXRATE      string `json:"-"`
	SFR                 string `json:"-"`
	TRFFICVOL           string `json:"trffvol"`
	AVGCTRLDLY          string `json:"-"`
	UNIFDLY             string `json:"-"`
	INCRDLY             string `json:"-"`
	INITQDLY            string `json:"-"`
	AVG_STDBYVH         string `json:"-"`
	AVG_STDBYLEN        string `json:"waitlen"`
	DEMAND_TRFFICVOLLEN string `json:"demand"`
	LOS                 string `json:"los"`
}

type DBLOC_30MIN_LOS_INFO_LOG struct {
	LOC_ID        int    `json:"-"`
	CREATE_DT     string `json:"-"`
	REF_TM_ID     int    `json:"-"`
	INFO_VALID_ID int    `json:"-"`
	AVG_LOS       int    `json:"-"`
	MAX_LOS       int    `json:"-"`
	MIN_LOS       int    `json:"-"`
}

type QRY_DBLOC_30MIN_LOS_INFO_LOG struct {
	LOC_ID        string `json:"-"`
	CREATE_DT     string `json:"-"`
	REF_TM_ID     string `json:"-"`
	INFO_VALID_ID string `json:"-"`
	AVG_LOS       string `json:"-"`
	MAX_LOS       string `json:"-"`
	MIN_LOS       string `json:"-"`
}

type SMART_LOC_FLOW_DIR struct {
	LOC_ID   int `json:"-"`
	FLOW_NUM int `json:"-"`
	LANE_DIR int `json:"-"`
}

type QRY_SMART_LOC_FLOW_DIR struct {
	LOC_ID   string `json:"-"`
	FLOW_NUM string `json:"-"`
	LANE_DIR string `json:"-"`
}

type DBTPO_TOD_SS_INFO struct {
	SS_ID     int    `json:"id"`
	GRP_ID    int    `json:"-"`
	PLANNUM   int    `json:"plan"`
	CREATE_DT string `json:"dt"`
	VALID_YN  int    `json:"-"`
}

type QRY_DBTPO_TOD_SS_INFO struct {
	SS_ID     string `json:"id"`
	GRP_ID    string `json:"-"`
	PLANNUM   string `json:"plan"`
	CREATE_DT string `json:"dt"`
	VALID_YN  string `json:"-"`
}

type DBTPO_GRP_DAYPLAN_RES struct {
	SS_ID              int `json:"-"`
	IDX                int `json:"idx"`
	HOUR               int `json:"hour"`
	MIN                int `json:"min"`
	CYCLE_LV           int `json:"lv"`
	LOC_OFFSETPLAN_IDX int `json:"offIdx"`
	LOC_PHASEPLAN_IDX  int `json:"phaseIdx"`
	CREATE_TYPE        int `json:"createtp"`
}

type QRY_DBTPO_GRP_DAYPLAN_RES struct {
	SS_ID              string `json:"-"`
	IDX                string `json:"idx"`
	HOUR               string `json:"hour"`
	MIN                string `json:"min"`
	CYCLE_LV           string `json:"lv"`
	LOC_OFFSETPLAN_IDX string `json:"offIdx"`
	LOC_PHASEPLAN_IDX  string `json:"phaseIdx"`
	CREATE_TYPE        string `json:"createtp"`
}

type DBTPO_GRP_CYCLEPLAN_RES struct {
	SS_ID    int `json:"-"`
	CYCLE_LV int `json:"cyclv"`
	CYCLELEN int `json:"cycLen"`
}

type QRY_DBTPO_GRP_CYCLEPLAN_RES struct {
	SS_ID    string `json:"-"`
	CYCLE_LV string `json:"cyclv"`
	CYCLELEN string `json:"cycLen"`
}

type DBTPO_GRP_DAYPLAN_CHK struct {
	SS_ID         int `json:"-"`
	IDX           int `json:"idx"`
	HOUR          int `json:"hour"`
	MIN           int `json:"min"`
	CYCLE_LV      int `json:"lv"`
	ORIGCYCLE     int `json:"originCyc"`
	CHG_REQ_CYCLE int `json:"reqCyc"`
}

type QRY_DBTPO_GRP_DAYPLAN_CHK struct {
	SS_ID         string `json:"-"`
	IDX           string `json:"idx"`
	HOUR          string `json:"hour"`
	MIN           string `json:"min"`
	CYCLE_LV      string `json:"lv"`
	ORIGCYCLE     string `json:"originCyc"`
	CHG_REQ_CYCLE string `json:"reqCyc"`
}

type DBTPO_LOC_PHASEPLAN_RES struct {
	SS_ID             int `json:"-"`
	LOC_ID            int `json:"-"`
	LOC_TYPE          int `json:"locType"`
	LOC_PHASEPLAN_IDX int `json:"idx"`
	OFFSET            int `json:"offset"`
	RINGA_1PHASE      int `json:"-"`
	RINGA_2PHASE      int `json:"-"`
	RINGA_3PHASE      int `json:"-"`
	RINGA_4PHASE      int `json:"-"`
	RINGA_5PHASE      int `json:"-"`
	RINGA_6PHASE      int `json:"-"`
	RINGA_7PHASE      int `json:"-"`
	RINGA_8PHASE      int `json:"-"`
	RINGB_1PHASE      int `json:"-"`
	RINGB_2PHASE      int `json:"-"`
	RINGB_3PHASE      int `json:"-"`
	RINGB_4PHASE      int `json:"-"`
	RINGB_5PHASE      int `json:"-"`
	RINGB_6PHASE      int `json:"-"`
	RINGB_7PHASE      int `json:"-"`
	RINGB_8PHASE      int `json:"-"`
	CREATE_TYPE       int `json:"createtp"`

	Aphases [8]int `json:"aPhases"`
	Bphases [8]int `json:"bPhases"`
}

type DBTPO_LOC_PHASEPLAN_RES2 struct {
	SS_ID             int    `json:"-"`
	HHMM              string `json:"hhmm"`
	LOC_ID            int    `json:"-"`
	LOC_TYPE          int    `json:"locType"`
	LOC_PHASEPLAN_IDX int    `json:"idx"`
	OFFSET            int    `json:"offset"`

	TOD_AVGCTRLDLY float64 `json:"avgdelay"`
	TOD_LOS        int     `json:"los"`

	OPRTOD_AVGCTRLDLY float64 `json:"-"`
	OPRTOD_LOS        int     `json:"-"`
	NEWTOD_AVGCTRLDLY float64 `json:"-"`
	NEWTOD_LOS        int     `json:"-"`

	CREATE_TYPE int `json:"createtp"`

	Aphases [8]int `json:"aPhases"`
	Bphases [8]int `json:"bPhases"`
}

type QRY_DBTPO_LOC_PHASEPLAN_RES2 struct {
	SS_ID             string `json:"-"`
	HHMM              string `json:"hhmm"`
	LOC_ID            string `json:"-"`
	LOC_TYPE          string `json:"locType"`
	LOC_PHASEPLAN_IDX string `json:"idx"`
	OFFSET            string `json:"offset"`

	TOD_AVGCTRLDLY string `json:"avgdelay"`
	TOD_LOS        string `json:"los"`

	OPRTOD_AVGCTRLDLY string `json:"-"`
	OPRTOD_LOS        string `json:"-"`
	NEWTOD_AVGCTRLDLY string `json:"-"`
	NEWTOD_LOS        string `json:"-"`

	CREATE_TYPE string `json:"createtp"`

	Aphases string `json:"aPhases"`
	Bphases string `json:"bPhases"`
}

type DBTPO_LOC_PHASEPLAN_CHK struct {
	SS_ID             int `json:"-"`
	LOC_ID            int `json:"-"`
	LOC_TYPE          int `json:"locType"`
	LOC_PHASEPLAN_IDX int `json:"idx"`
	CYCLELEN          int `json:"cycle"`
	RINGA_1PHASE      int `json:"-"`
	RINGA_2PHASE      int `json:"-"`
	RINGA_3PHASE      int `json:"-"`
	RINGA_4PHASE      int `json:"-"`
	RINGA_5PHASE      int `json:"-"`
	RINGA_6PHASE      int `json:"-"`
	RINGA_7PHASE      int `json:"-"`
	RINGA_8PHASE      int `json:"-"`
	RINGB_1PHASE      int `json:"-"`
	RINGB_2PHASE      int `json:"-"`
	RINGB_3PHASE      int `json:"-"`
	RINGB_4PHASE      int `json:"-"`
	RINGB_5PHASE      int `json:"-"`
	RINGB_6PHASE      int `json:"-"`
	RINGB_7PHASE      int `json:"-"`
	RINGB_8PHASE      int `json:"-"`

	Aphases [8]int `json:"aPhases"`
	Bphases [8]int `json:"bPhases"`
}

type DBTPO_LOC_PHASEPLAN_CHK2 struct {
	SS_ID             int    `json:"-"`
	HHMM              string `json:"hhmm"`
	LOC_ID            int    `json:"-"`
	LOC_TYPE          int    `json:"locType"`
	LOC_PHASEPLAN_IDX int    `json:"idx"`
	HOUR              int    `json:"hour"`
	MIN               int    `json:"min"`
	CYCLELEN          int    `json:"cycle"`

	Aphases [8]int `json:"aPhases"`
	Bphases [8]int `json:"bPhases"`

	TOD_AVGCTRLDLY float64 `json:"avgdelay"`
	TOD_LOS        int     `json:"los"`
}

type QRY_DBTPO_LOC_PHASEPLAN_CHK2 struct {
	SS_ID             string `json:"-"`
	HHMM              string `json:"hhmm"`
	LOC_ID            string `json:"-"`
	LOC_TYPE          string `json:"locType"`
	LOC_PHASEPLAN_IDX string `json:"idx"`
	HOUR              string `json:"hour"`
	MIN               string `json:"min"`
	CYCLELEN          string `json:"cycle"`

	Aphases string `json:"aPhases"`
	Bphases string `json:"bPhases"`

	TOD_AVGCTRLDLY string `json:"avgdelay"`
	TOD_LOS        string `json:"los"`
}

type TPO_TOD_LOC_DELAYINFO struct {
	SS_ID          int     `json:"-"`
	HOUR           int     `json:"-"`
	MIN            int     `json:"-"`
	HHMM           string  `json:"-"`
	LOC_ID         int     `json:"-"`
	CREATE_TYPE    int     `json:"-"`
	TOD_AVGCTRLDLY float64 `json:"avgdelay"`
	TOD_LOS        int     `json:"los"`

	OPRTOD_AVGCTRLDLY float64 `json:"-"`
	OPRTOD_LOS        int     `json:"-"`
	NEWTOD_AVGCTRLDLY float64 `json:"-"`
	NEWTOD_LOS        int     `json:"-"`
}

type TPO_TOD_CREATE_MNGINFO struct {
	MNG_ST_TIME            int     `json:"mngTm"`
	OPT_REQ_LOS_DIFF       int     `json:"losDiff"`
	OPT_REF_LOS            int     `json:"refLos"`
	CYCLE_OPT_REF_LOS_DIFF int     `json:"cycOptLosDiff"`
	REFINFO_CNT            int     `json:"refInfoCnt"`
	OPT_REF_RATE           int     `json:"optRefRate"`
	CYCLE_INCRE_VOL        int     `json:"cycIncrVol"`
	AVG_VH_LEN             float32 `json:"avgVhLen"`
	BUS_CONVCOE            float32 `json:"busConvCoe"`
	BGVH_CONVCOE           float32 `json:"bgVhConvCoe"`
	SMVH_CONVCOE           float32 `json:"smVhConvCoe"`
	OPT_MODE               int     `json:"optMode"`
}

type QRY_TPO_TOD_CREATE_MNGINFO struct {
	MNG_ST_TIME            string `json:"mngTm"`
	OPT_REQ_LOS_DIFF       string `json:"losDiff"`
	OPT_REF_LOS            string `json:"refLos"`
	CYCLE_OPT_REF_LOS_DIFF string `json:"cycOptLosDiff"`
	REFINFO_CNT            string `json:"refInfoCnt"`
	OPT_REF_RATE           string `json:"optRefRate"`
	CYCLE_INCRE_VOL        string `json:"cycIncrVol"`
	AVG_VH_LEN             string `json:"avgVhLen"`
	BUS_CONVCOE            string `json:"busConvCoe"`
	BGVH_CONVCOE           string `json:"bgVhConvCoe"`
	SMVH_CONVCOE           string `json:"smVhConvCoe"`
	OPT_MODE               string `json:"optMode"`
}

// ----------------------------------------------------
// 이력정보 조회 struct
// ----------------------------------------------------
type DBLOC_FLOW_30M_LOS_INFO_LOG struct {
	LOC_ID         int     `json:"id"`
	FLOW_NUM       int     `json:"flow"`
	CREATE_DT      string  `json:"dt"`
	REF_TM_ID      int     `json:"refTmId"`
	INFO_VALID_YN  int     `json:"invalidYn"`
	CAPA           float32 `json:"capa"`
	TRFFICVOLRATE  float32 `json:"trffVolRate"`
	SATUR          float32 `json:"satur"`
	FIR_STDBYVH    float32 `json:"firStdbyVh"`
	HVY_VH_CORRVAL float32 `json:"hVhCorrVar"`
	HVY_VH_MIXRATE float32 `json:"hVhMixRate"`
	SFR            float32 `json:"sfr"`
	TRFFICVOL      float32 `json:"trfficVol"`
	AVGCTRLDLY     float32 `json:"avgCtrlDly"`
	UNIFDLY        float32 `json:"unifDly"`
	INCRDLY        float32 `json:"increDly"`
	INITQDLY       float32 `json:"initQdly"`
	LOS            int     `json:"los"`
}

type QRY_DBLOC_FLOW_30M_LOS_INFO_LOG struct {
	LOC_ID         string `json:"id"`
	FLOW_NUM       string `json:"flow"`
	CREATE_DT      string `json:"dt"`
	REF_TM_ID      string `json:"refTmId"`
	INFO_VALID_YN  string `json:"invalidYn"`
	CAPA           string `json:"capa"`
	TRFFICVOLRATE  string `json:"trffVolRate"`
	SATUR          string `json:"satur"`
	FIR_STDBYVH    string `json:"firStdbyVh"`
	HVY_VH_CORRVAL string `json:"hVhCorrVar"`
	HVY_VH_MIXRATE string `json:"hVhMixRate"`
	SFR            string `json:"sfr"`
	TRFFICVOL      string `json:"trfficVol"`
	AVGCTRLDLY     string `json:"avgCtrlDly"`
	UNIFDLY        string `json:"unifDly"`
	INCRDLY        string `json:"increDly"`
	INITQDLY       string `json:"initQdly"`
	LOS            string `json:"los"`
}

type DBLOC_30M_SIGOPRINFO_LOG struct {
	LOC_ID           int     `json:"id"`
	CREATE_DT        string  `json:"dt"`
	REF_TM_ID        int     `json:"refTmId"`
	INFO_VALID_YN    int     `json:"invalidYn"`
	SIGOPRINFO_COUNT int     `json:"infoCnt"`
	AVG_CYCLELEN     float32 `json:"avgCyc"`

	RINGA_1PHASE float32 `json:"ringA1"`
	RINGA_2PHASE float32 `json:"ringA2"`
	RINGA_3PHASE float32 `json:"ringA3"`
	RINGA_4PHASE float32 `json:"ringA4"`
	RINGA_5PHASE float32 `json:"ringA5"`
	RINGA_6PHASE float32 `json:"ringA6"`
	RINGA_7PHASE float32 `json:"ringA7"`
	RINGA_8PHASE float32 `json:"ringA8"`

	RINGB_1PHASE float32 `json:"ringB1"`
	RINGB_2PHASE float32 `json:"ringB2"`
	RINGB_3PHASE float32 `json:"ringB3"`
	RINGB_4PHASE float32 `json:"ringB4"`
	RINGB_5PHASE float32 `json:"ringB5"`
	RINGB_6PHASE float32 `json:"ringB6"`
	RINGB_7PHASE float32 `json:"ringB7"`
	RINGB_8PHASE float32 `json:"ringB8"`
}

type QRY_DBLOC_30M_SIGOPRINFO_LOG struct {
	LOC_ID           string `json:"id"`
	CREATE_DT        string `json:"dt"`
	REF_TM_ID        string `json:"refTmId"`
	INFO_VALID_YN    string `json:"invalidYn"`
	SIGOPRINFO_COUNT string `json:"infoCnt"`
	AVG_CYCLELEN     string `json:"avgCyc"`

	RINGA_1PHASE string `json:"ringA1"`
	RINGA_2PHASE string `json:"ringA2"`
	RINGA_3PHASE string `json:"ringA3"`
	RINGA_4PHASE string `json:"ringA4"`
	RINGA_5PHASE string `json:"ringA5"`
	RINGA_6PHASE string `json:"ringA6"`
	RINGA_7PHASE string `json:"ringA7"`
	RINGA_8PHASE string `json:"ringA8"`

	RINGB_1PHASE string `json:"ringB1"`
	RINGB_2PHASE string `json:"ringB2"`
	RINGB_3PHASE string `json:"ringB3"`
	RINGB_4PHASE string `json:"ringB4"`
	RINGB_5PHASE string `json:"ringB5"`
	RINGB_6PHASE string `json:"ringB6"`
	RINGB_7PHASE string `json:"ringB7"`
	RINGB_8PHASE string `json:"ringB8"`
}

type DBLOC_30M_LOS_INFO_LOG struct {
	LOC_ID        int     `json:"id"`
	CREATE_DT     string  `json:"dt"`
	REF_TM_ID     int     `json:"refTmId"`
	INFO_VALID_YN int     `json:"invalidYn"`
	AVG_LOS       int     `json:"avgLos"`
	MAX_LOS       int     `json:"maxLos"`
	MIN_LOS       int     `json:"minLos"`
	CLG_CDS_SUM   float32 `json:"clgCdsSum"`

	RINGA_1PHASE_CDS float32 `json:"ringA1Cds"`
	RINGA_2PHASE_CDS float32 `json:"ringA2Cds"`
	RINGA_3PHASE_CDS float32 `json:"ringA3Cds"`
	RINGA_4PHASE_CDS float32 `json:"ringA4Cds"`
	RINGA_5PHASE_CDS float32 `json:"ringA5Cds"`
	RINGA_6PHASE_CDS float32 `json:"ringA6Cds"`
	RINGA_7PHASE_CDS float32 `json:"ringA7Cds"`
	RINGA_8PHASE_CDS float32 `json:"ringA8Cds"`

	RINGB_1PHASE_CDS float32 `json:"ringB1Cds"`
	RINGB_2PHASE_CDS float32 `json:"ringB2Cds"`
	RINGB_3PHASE_CDS float32 `json:"ringB3Cds"`
	RINGB_4PHASE_CDS float32 `json:"ringB4Cds"`
	RINGB_5PHASE_CDS float32 `json:"ringB5Cds"`
	RINGB_6PHASE_CDS float32 `json:"ringB6Cds"`
	RINGB_7PHASE_CDS float32 `json:"ringB7Cds"`
	RINGB_8PHASE_CDS float32 `json:"ringB8Cds"`

	OPT_REQ_YN        int     `json:"optReqYn"`
	PLAN_CHK_APPLY_YN int     `json:"planChkApplyYn"`
	AVG_CTRLDLY       float32 `json:"avgCtrlDly"`
	PLAN_NO           int     `json:"planNo"`
}

type QRY_DBLOC_30M_LOS_INFO_LOG struct {
	LOC_ID        string `json:"id"`
	CREATE_DT     string `json:"dt"`
	REF_TM_ID     string `json:"refTmId"`
	INFO_VALID_YN string `json:"invalidYn"`
	AVG_LOS       string `json:"avgLos"`
	MAX_LOS       string `json:"maxLos"`
	MIN_LOS       string `json:"minLos"`
	CLG_CDS_SUM   string `json:"clgCdsSum"`

	RINGA_1PHASE_CDS string `json:"ringA1Cds"`
	RINGA_2PHASE_CDS string `json:"ringA2Cds"`
	RINGA_3PHASE_CDS string `json:"ringA3Cds"`
	RINGA_4PHASE_CDS string `json:"ringA4Cds"`
	RINGA_5PHASE_CDS string `json:"ringA5Cds"`
	RINGA_6PHASE_CDS string `json:"ringA6Cds"`
	RINGA_7PHASE_CDS string `json:"ringA7Cds"`
	RINGA_8PHASE_CDS string `json:"ringA8Cds"`

	RINGB_1PHASE_CDS string `json:"ringB1Cds"`
	RINGB_2PHASE_CDS string `json:"ringB2Cds"`
	RINGB_3PHASE_CDS string `json:"ringB3Cds"`
	RINGB_4PHASE_CDS string `json:"ringB4Cds"`
	RINGB_5PHASE_CDS string `json:"ringB5Cds"`
	RINGB_6PHASE_CDS string `json:"ringB6Cds"`
	RINGB_7PHASE_CDS string `json:"ringB7Cds"`
	RINGB_8PHASE_CDS string `json:"ringB8Cds"`

	OPT_REQ_YN        string `json:"optReqYn"`
	PLAN_CHK_APPLY_YN string `json:"planChkApplyYn"`
	AVG_CTRLDLY       string `json:"avgCtrlDly"`
	PLAN_NO           string `json:"planNo"`
}

// ----------------------------------------------------
type TEMPLOCAL_RST struct {
	LocId  int
	CollDt string

	CtrlMode  int
	CtrlState int

	Aring [8]int
	Bring [8]int

	Cycle int
}

type SIGMAP struct {
	Skey  string `json :"-"`
	LocId int
	Bank  int
	Ring  int
	Step  int     `json :"-"`
	Lsu   [16]int `json :"-"`

	Min int `json :"-"`
	Max int `json :"-"`
	Eop int `json :"-"`
	Ar  [608]byte
	Br  [608]byte
}

type SigMapData struct {
	LocId int       `json:"locid"`
	Data  [1216]int `json:"map"`
}

// 통합 이벤트 이력 db 전달 parameter
type EventLogParam struct {
	// 1. value
	// 1) 필수 파라미터
	SDt      string
	EDt      string
	Page     int
	ReqRows  int
	SearchTp int
	StrSort  string

	Offset int // offset

	// 2) 시스템 / 그룹 / 교차로 필수 파라미터 (uType, uId)
	UType int
	UId   int

	// 3) 운영자 필수 파라미터 (ctor, uId)
	Ctor int

	// 4) 선택 파라미터
	EvCode     int
	EvLowCodes []int  // [...]
	EvLowCode  string // "1, 2, 3, ..."

	// 2. check exist value
	// 1) 필수 파라미터
	SDtOk      bool
	EDtOk      bool
	PageOk     bool
	ReqRowsOk  bool
	SearchTpOk bool

	// 2) 시스템 / 그룹 / 교차로 필수 파라미터 (uType, uId)
	UTypeOk bool
	UIdOk   bool

	// 3) 운영자 필수 파라미터 (ctor, uId)
	CtorOk bool

	// 4) 선택 파라미터
	EvCodeOk    bool
	EvLowCodeOk bool

	SortFields []SortField
	Sortfield  string
}

// 여러 필드값으로 sort
type SortField struct {
	Field string `json:"field"` // 필드명
	Order int    `json:"order"` // 0: 오름차순 / 1: 내림차순
}
