package input

import (
	"fmt"
	"sync"
	"time"

	"github.com/mitchellh/hashstructure"

	"github.com/elastic/beats/filebeat/channel"
	"github.com/elastic/beats/filebeat/input/file"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

// Input is the interface common to all prospectors
type Input interface {
	Run()
	Stop()
	Wait()
}

// InputRunner contains the prospector
type InputRunner struct {
	config   inputConfig
	input    Input
	done     chan struct{}
	wg       *sync.WaitGroup
	ID       uint64
	Once     bool
	beatDone chan struct{}
}

// NewProspector instantiates a new prospector
func New(
	conf *common.Config,
	outlet channel.Factory,
	beatDone chan struct{},
	states []file.State,
	dynFields *common.MapStrPointer,
) (*InputRunner, error) {
	input := &InputRunner{
		config:   defaultConfig,
		wg:       &sync.WaitGroup{},
		done:     make(chan struct{}),
		Once:     false,
		beatDone: beatDone,
	}

	var err error
	if err = conf.Unpack(&input.config); err != nil {
		return nil, err
	}

	var h map[string]interface{}
	conf.Unpack(&h)
	input.ID, err = hashstructure.Hash(h, nil)
	if err != nil {
		return nil, err
	}

	var f Factory
	f, err = GetFactory(input.config.Type)
	if err != nil {
		return input, err
	}

	context := Context{
		States:        states,
		Done:          input.done,
		BeatDone:      input.beatDone,
		DynamicFields: dynFields,
	}
	var ipt Input
	ipt, err = f(conf, outlet, context)
	if err != nil {
		return input, err
	}
	input.input = ipt

	return input, nil
}

// Start starts the prospector
func (p *InputRunner) Start() {
	p.wg.Add(1)
	logp.Info("Starting prospector of type: %v; ID: %d ", p.config.Type, p.ID)

	onceWg := sync.WaitGroup{}
	if p.Once {
		// Make sure start is only completed when Run did a complete first scan
		defer onceWg.Wait()
	}

	onceWg.Add(1)
	// Add waitgroup to make sure prospectors finished
	go func() {
		defer func() {
			onceWg.Done()
			p.stop()
			p.wg.Done()
		}()

		p.Run()
	}()
}

// Run starts scanning through all the file paths and fetch the related files. Start a harvester for each file
func (p *InputRunner) Run() {
	// Initial prospector run
	p.input.Run()

	// Shuts down after the first complete run of all prospectors
	if p.Once {
		return
	}

	for {
		select {
		case <-p.done:
			logp.Info("Prospector ticker stopped")
			return
		case <-time.After(p.config.ScanFrequency):
			logp.Debug("prospector", "Run prospector")
			p.input.Run()
		}
	}
}

// Stop stops the prospector and with it all harvesters
func (p *InputRunner) Stop() {
	// Stop scanning and wait for completion
	close(p.done)
	p.wg.Wait()
}

func (p *InputRunner) stop() {
	logp.Info("Stopping Prospector: %d", p.ID)

	// In case of once, it will be waited until harvesters close itself
	if p.Once {
		p.input.Wait()
	} else {
		p.input.Stop()
	}
}

func (p *InputRunner) String() string {
	return fmt.Sprintf("prospector [type=%s, ID=%d]", p.config.Type, p.ID)
}
