
package kpnmlogger

import(
	io    "io"
	os    "os"
	sync  "sync"
	time  "time"

	util  "github.com/zyxgad/go-util/util"
)

var (
	_logger *Logger = NewLogger(7)
)

type Logger struct{
	out    io.Writer

	block   sync.Locker
	plock   sync.Locker
}

func NewLogger(bufsize int)(lgr *Logger){
	return &Logger{
		out: os.Stderr,

		block: new(sync.Mutex),
		plock: new(sync.Mutex),
	}
}

func (lgr *Logger)SetWriter(out io.Writer){
	lgr.out = out
}

func (lgr *Logger)Println(msg string){
	lgr.plock.Lock()

	lgr.pushLine(([]byte)(msg))

	lgr.plock.Unlock()
}

func (lgr *Logger)pushLine(msg []byte){
	lgr.block.Lock()

	lgr.out.Write(msg)
	lgr.out.Write(([]byte)("\n"))

	lgr.block.Unlock()
}

func GetLogFileName()(string){
	return time.Now().Format("20060102-15.log")
}

func ChangeLogWithDir(dirPath string){
	util.CreateDir(dirPath)
	file := util.JoinPath(dirPath, GetLogFileName())

	logf, err := os.OpenFile(file, os.O_CREATE | os.O_APPEND | os.O_WRONLY | os.O_SYNC, os.ModePerm)
	if err != nil {
		panic(err)
		return
	}
	_logger.SetWriter(io.MultiWriter(logf, os.Stderr))
}

var (
	_ACLklLock chan bool = nil
)

func AutoChangeLog(dirPath string){
	if _ACLklLock == nil {
		_ACLklLock = make(chan bool, 0)
	}else{
		_ACLklLock <- true
	}
	ChangeLogWithDir(dirPath)
	go func(){
		for {
			select{
			case <-_ACLklLock:
				return
			case <-time.After(time.Minute * 15):
				ChangeLogWithDir(dirPath)
			}
		}
	}()
}

type loggerSource int

func (loggerSource)Init(){
	AutoChangeLog(util.JoinPath("/", "var", "server_logs"))
}

