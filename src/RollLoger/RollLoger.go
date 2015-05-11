package RollLoger

import (
	//"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"winsvc/osext"
)

const (
	infoStr  = "[Info ]- "
	debugStr = "[Debug]- "
	errorStr = "[Error]- "
	fatalStr = "[Fatal]- "
)

func init() {
	homedir, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatalln(err)
	}
	rLogger = NewRollLogger(1024*5, "log", homedir+"/log")
}

var rLogger *RollLog

func NewRollLogger(MaxSize int64, FileName string, FileDirPath string) *RollLog {
	_rLogger := new(RollLog)
	_rLogger.innerWriter = &Writers_rollingfilewriter{MaxSize, FileName, FileDirPath, "", 0, nil, ""}
	writers := []io.Writer{
		_rLogger.innerWriter,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	_rLogger.innerlogger = log.New(fileAndStdoutWriter, "", log.Ldate|log.Ltime|log.Llongfile)
	return _rLogger
}
func Debug(msgs ...string) {

	rLogger.Output(2, debugStr+strings.Join(msgs, ""))
}
func Info(msgs ...string) {
	rLogger.Output(2, infoStr+strings.Join(msgs, ""))
}
func Error(msgs ...string) {
	rLogger.Output(2, errorStr+strings.Join(msgs, ""))
}

func Fatal(err error) {
	rLogger.innerlogger.Fatalln(err)

}

func Close() error {
	err := rLogger.Close()
	return err
}

type RollLog struct {
	innerlogger *log.Logger
	innerWriter *Writers_rollingfilewriter
}

func (p *RollLog) Close() error {
	err := p.innerWriter.Close()
	return err
}

func (p *RollLog) Output(callDepth int, msg string) {
	p.innerlogger.Output(callDepth+1, msg)
	//p.innerlogger.Println(msg)
}

type Writers_rollingfilewriter struct {
	MaxSize         int64 //文件最大SIZE 单位KB
	FileName        string
	FileDirPath     string
	currentFileName string
	currentFileSize int64
	currentFile     *os.File
	currentFileDate string
}

func NewRollFileWriter(MaxSize int64, FileName string, FileDirPath string) (*Writers_rollingfilewriter, error) {
	return &Writers_rollingfilewriter{MaxSize, FileName, FileDirPath, "", 0, nil, ""}, nil
	//return nil, nil
}

func (w *Writers_rollingfilewriter) Write(p []byte) (n int, err error) {

	if w.currentFile == nil {
		err := w.createFileAndFolderIfNeeded()
		if err != nil {
			return 0, err
		}
	}

	//检查日期是否已经改变
	nrByDate, err := w.checkToRollByDate()
	if err != nil {
		return 0, err
	}
	if nrByDate {
		w.currentFile.Close()
		w.currentFile = nil
		w.createFileAndFolderIfNeeded()
	}

	//检查是否超过最大SIZE
	nrBySize, err := w.checkToRollBySize()
	if err != nil {
		return 0, err
	}
	if nrBySize {
		w.currentFile.Close()
		w.currentFile = nil

		err := w.RollFileBySize(1)
		if err != nil {
			return 0, err
		}

		err = os.Rename(w.getLogName(0), w.getLogName(1))
		if err != nil {
			return 0, err
		}

		err = w.createFileAndFolderIfNeeded()
		if err != nil {
			return 0, err
		}
	}

	w.currentFileSize += int64(len(p))
	return w.currentFile.Write(p)
}

//获取日志文件名
func (w *Writers_rollingfilewriter) getLogName(i int) string {
	if i == 0 {
		return filepath.Join(w.FileDirPath, w.FileName+"_"+w.currentFileDate+".log")
	} else {
		return filepath.Join(w.FileDirPath, w.FileName+"_"+w.currentFileDate+"_"+strconv.Itoa(i)+".log")
	}

}

//检查文件序号为i 的文件是否存在，如果存在则递增改名（递归）
func (w *Writers_rollingfilewriter) RollFileBySize(i int) error {

	nextName := w.getLogName(i + 1)
	currentName := w.getLogName(i)

	_, err := os.Lstat(currentName)
	if err == nil {

		err = w.RollFileBySize(i + 1)
		if err != nil {
			return err
		}

		err = os.Rename(currentName, nextName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Writers_rollingfilewriter) checkToRollByDate() (bool, error) {

	return time.Now().Format("20060102") != w.currentFileDate, nil
}

func (w *Writers_rollingfilewriter) checkToRollBySize() (bool, error) {
	return w.currentFileSize >= w.MaxSize*1024, nil
}

//检查文件夹与文件情况，如果不存在则创建
func (w *Writers_rollingfilewriter) createFileAndFolderIfNeeded() error {
	if len(w.FileDirPath) != 0 {
		err := os.MkdirAll(w.FileDirPath, 0777)
		if err != nil {
			return err
		}
	}

	w.currentFileDate = time.Now().Format("20060102")
	filePath := w.getLogName(0)
	stat, err := os.Lstat(filePath)
	if err == nil {
		w.currentFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0777)

		stat, err = os.Lstat(filePath)
		if err != nil {
			return err
		}

		w.currentFileSize = stat.Size()
	} else {
		w.currentFile, err = os.Create(filePath)
		w.currentFileSize = 0
	}

	return nil
}

func (w *Writers_rollingfilewriter) Close() error {
	if w.currentFile != nil {
		err := w.currentFile.Close()
		if err != nil {
			return err
		}
		w.currentFile = nil
	}
	return nil
}
