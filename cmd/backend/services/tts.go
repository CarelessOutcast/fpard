package services

/*
#cgo CFLAGS: -I${SRCDIR}/tts/include
#cgo LDFLAGS: -L${SRCDIR}/tts/lib -lespeak-ng -lstdc++ -lpthread
#include <stdlib.h>
#include "tts/include/speak_lib.h"
extern int synthCallback(short*, int, espeak_EVENT*);
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

type ESpeak struct{}

var pcmBuffer []int16
var pcmMutex sync.Mutex

//export synthCallback
func synthCallback(wav *C.short, numsamples C.int, events *C.espeak_EVENT) C.int {
	if numsamples > 0 {
		pcmMutex.Lock()
		defer pcmMutex.Unlock()

		data := (*[1 << 30]C.short)(unsafe.Pointer(wav))[:numsamples:numsamples]
		for _, sample := range data {
			pcmBuffer = append(pcmBuffer, int16(sample))
		}
	}
	return 0
}

func NewESpeak() (*ESpeak, error) {
	// TODO: Create a pure C implementation, and make this a wrapper
	if C.espeak_Initialize(C.AUDIO_OUTPUT_SYNCHRONOUS, 0, nil, 0) == -1 {
		return nil, fmt.Errorf("failed to initialize espeak NG")
	}
	C.espeak_SetSynthCallback((*C.t_espeak_callback)(C.synthCallback))
	return &ESpeak{}, nil
}

func (e *ESpeak) SetVoice(voiceName string) error {
	cVoice := C.CString(voiceName)
	defer C.free(unsafe.Pointer(cVoice))

	if C.espeak_SetVoiceByName(cVoice) != 0 {
		return fmt.Errorf("failed to set voice: %s", voiceName)
	}
	return nil
}

func (e *ESpeak) Speak(text string) ([]int16, error) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	// syntesize speech
	if C.espeak_Synth(unsafe.Pointer(cText), C.size_t(len(text)), 0, 0, 0, C.espeakCHARS_AUTO, nil, nil) != 0 {
		return nil, fmt.Errorf("failed to synthesize speech")
	}

	// wait for synth to finish
	C.espeak_Synchronize()

	// Convert C buffer into Go buffer
	return pcmBuffer, nil
}
