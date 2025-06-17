package logger

import "go_redis/general"

// //---------------------------------------------------------------------------
// // Log
// //---------------------------------------------------------------------------
var Mlog *general.OLog2 = general.InitLogEnv("./log", "main", 1) // level 1~9 main process
