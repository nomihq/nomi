package cli

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/nullswan/nomi/internal/audio"
	"github.com/nullswan/nomi/internal/config"
	"github.com/nullswan/nomi/internal/logger"
	"github.com/nullswan/nomi/internal/transcription"
)

// InitTranscriptionServer initializes the Transcription Server with predefined buffer settings.
func InitTranscriptionServer(
	oaiKey string,
	audioOpts *audio.StreamParameters,
	log *logger.Logger,
	callback transcription.TranscriptionServerCallbackT,
	language string,
) (*transcription.TranscriptionServer, error) {
	bufferManagerPrimary := transcription.NewSimpleBufferManager(audioOpts)
	bufferManagerPrimary.SetMinBufferDuration(1 * time.Second)

	textReconciler := transcription.NewTextReconciler(log)
	tsHandler := transcription.NewTranscriptionHandler(
		oaiKey,
		audioOpts,
		log,
	)
	if language != "" {
		lang, err := transcription.LoadLangFromValue(language)
		if err != nil {
			return nil, fmt.Errorf("invalid language code: %w", err)
		}
		tsHandler.WithLanguage(lang)
	}

	return transcription.NewTranscriptionServer(
		bufferManagerPrimary,
		nil,
		tsHandler,
		textReconciler,
		log,
		callback,
	), nil
}

func InitVoice(
	cfg *config.Config,
	log *logger.Logger,
	handleTranscription func(text string, isProcessing bool),
	cmdKeyCode uint16,
	language string,
) (*audio.StreamHandler, <-chan struct{}, <-chan struct{}, error) {
	if !cfg.Input.Voice.Enabled {
		return nil, nil, nil, nil
	}

	if err := portaudio.Initialize(); err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to initialize PortAudio: %w",
			err,
		)
	}

	audioOpts, err := audio.ComputeDefaultAudioOptions()
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"error computing audio options: %w",
			err,
		)
	}

	oaiKey := os.Getenv("OPENAI_API_KEY")
	if oaiKey == "" {
		return nil, nil, nil, errors.New("OPENAI_API_KEY is not set")
	}

	ts, err := InitTranscriptionServer(
		oaiKey,
		audioOpts,
		log,
		handleTranscription,
		language,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"error initializing transcription server: %w",
			err,
		)
	}
	err = ts.Start()
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to start transcription server: %w",
			err,
		)
	}

	inputStream, err := audio.NewDefaultInputStream(
		log,
		audioOpts,
		func(buffer []float32) {
			data, err := audio.Float32ToPCM(buffer)
			if err != nil {
				log.With("error", err).
					Error("Failed to convert float32 to PCM")
				return
			}

			ts.AddAudio(data)
		},
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to create input stream: %w",
			err,
		)
	}

	audioStartCh, audioEndCh := SetupKeyHooks(cmdKeyCode, ts)
	return inputStream, audioStartCh, audioEndCh, nil
}
