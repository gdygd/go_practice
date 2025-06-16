package logger

import "go_redis/general"

// define global log variable
// http server logger
// socket server logger (socket, msg)
// main process logger

// //---------------------------------------------------------------------------
// // Log
// //---------------------------------------------------------------------------
var Apilog *general.OLog2 = general.InitLogEnv("./log", "apiserver", 1) // level 1~9, api server
var Mlog *general.OLog2 = general.InitLogEnv("./log", "main", 1)        // level 1~9 main process
