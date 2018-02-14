package scannergo

import (
	"fmt"
	"strings"

	"github.com/rhettcp/keylogger"
)

type Scanner struct {
	keylog       *keylogger.KeyLogger
	Scans        chan string
	keepScanning bool
	stopWork     chan bool
}

func NewScanner(scannerName string) (*Scanner, error) {
	s := &Scanner{Scans: make(chan string), stopWork: make(chan bool)}
	devs, err := keylogger.NewDevices()
	if err != nil {
		return nil, err
	}
	qrDevice := -1
	for _, val := range devs {
		if val.Name == "\""+scannerName+"\"" {
			qrDevice = val.Id
			break
		}
	}
	if qrDevice == -1 {
		return nil, fmt.Errorf("QR SCANNER NOT FOUND")
	}
	//our keyboard..on your system, it will be diffrent
	s.keylog = keylogger.NewKeyLogger(devs[qrDevice])
	return s, nil
}

func (s *Scanner) StartScanning() error {
	in, err := s.keylog.Read()
	if err != nil {
		return err
	}
	s.keepScanning = true
	go func() {
		str := ""
		store := true
		for {
			select {
			case <-s.stopWork:
				return
			case i := <-in:
				//we only need keypress
				if i.Type == keylogger.EV_KEY {
					if store {
						if i.KeyString() == "ENTER" {
							s.Scans <- strings.ToLower(strings.Replace(str, "SPACE", " ", -1))
							str = ""
						} else {
							str = str + i.KeyString()
						}
					}
					store = !store
				}
			}
		}
	}()
	return nil
}

func (s *Scanner) StopScanning() {
	s.stopWork <- true
}
